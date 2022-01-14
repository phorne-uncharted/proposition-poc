package routes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// TreemapItem is a transformed graph to match the expected treemap structure.
type TreemapItem struct {
	Children []*TreemapItem `json:"children"`
	Name     string         `json:"name"`
	ColName  string         `json:"colname,omitempty"`
	Value    int            `json:"value,omitempty"`
}

// Graph is a collection of nodes with a single root.
type Graph struct {
	URL  string `json:"url"`
	Root *Node  `json:"root"`
}

// Node is one entity in a graph.
type Node struct {
	Key        string       `json:"key"`
	Neighbours []*Node      `json:"neighbours"`
	Data       *Proposition `json:"data"`
}

// Proposition is an entity being extracted from a site.
type Proposition struct {
	ID            string   `json:"id"`
	FullName      string   `json:"fullname"`
	Tag           string   `json:"tag"`
	Code          string   `json:"code"`
	URL           string   `json:"url"`
	Key           string   `json:"key"`
	PotentialTags []string `json:"potentialTags"`
	ParentURL     string   `json:"parentUrl"`
}

func buildTreemap(graph *Graph) *TreemapItem {
	return nodeToItem(1, graph.Root)
}

func nodeToItem(depth int, node *Node) *TreemapItem {
	colName := ""
	if depth > 1 {
		colName = fmt.Sprintf("level%d", depth)
	}
	item := &TreemapItem{
		Name:     node.Data.Tag,
		ColName:  colName,
		Children: make([]*TreemapItem, len(node.Neighbours)),
	}

	for i, c := range node.Neighbours {
		item.Children[i] = nodeToItem(depth+1, c)
	}

	if len(item.Children) == 0 {
		item.Value = 20
	}

	return item
}

// ToPropertySlice converts a proposition to a string slice.
func (p *Proposition) ToPropertySlice() []string {
	return []string{
		p.ID,
		p.FullName,
		p.Tag,
		p.Code,
		p.URL,
	}
}

// LinksHandler generates a route handler that returns links.
func LinksHandler(allowedSites []string) func(http.ResponseWriter, *http.Request) {
	allowedSitesMap := map[string]bool{}
	for _, s := range allowedSites {
		allowedSitesMap[s] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params, err := getPostParameters(r)
		if err != nil {
			handleError(w, errors.Wrap(err, "Unable to parse post parameters"))
			return
		}

		urlRaw := params["url"].(string)
		urlParsed, err := url.Parse(urlRaw)
		if err != nil {
			handleError(w, errors.Wrap(err, "Unable to parse url"))
			return
		}

		if !allowedSitesMap[urlParsed.Hostname()] {
			handleError(w, errors.Errorf("host '%s' not allowed", urlParsed.Hostname()))
			return
		}

		graph := buildGraph(urlParsed)
		treemap := buildTreemap(graph)

		// marshal data
		err = handleJSON(w, treemap)
		if err != nil {
			handleError(w, errors.Wrap(err, "unable to marshal model result into JSON"))
			return
		}
	}
}

func buildGraph(urlParsed *url.URL) *Graph {
	nodes := map[string]*Node{}
	c := colly.NewCollector(
		colly.AllowedDomains("onedemo-telco.azurewebsites.net"),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Ctx.Put("parent", e.Request.URL.String())
		e.Request.Visit(e.Attr("href"))
	})

	c.OnResponse(func(r *colly.Response) {
		prop, _ := processLink(r)
		nodes = addNode(nodes, prop)
	})

	c.Visit(urlParsed.String())

	return processGraph(urlParsed.String(), nodes)
}

func processGraph(url string, nodes map[string]*Node) *Graph {
	root := nodes[""]
	root = buildBreadCrumb(nil, root, []*Node{})[0]
	return &Graph{
		URL:  url,
		Root: root,
	}
}

func buildBreadCrumb(parent *Node, current *Node, alreadyProcessed []*Node) []*Node {
	if parent == nil {
		current.Data.FullName = fmt.Sprintf("/")
	} else {
		current.Data.FullName = fmt.Sprintf("%s%s/", parent.Data.FullName, current.Data.Tag)
		alreadyProcessed = append(alreadyProcessed, current)
	}

	for _, c := range current.Neighbours {
		alreadyProcessed = buildBreadCrumb(current, c, alreadyProcessed)
	}

	return alreadyProcessed
}

func addNode(nodes map[string]*Node, proposition *Proposition) map[string]*Node {
	newNode := &Node{
		Key:        proposition.URL,
		Neighbours: []*Node{},
		Data:       proposition,
	}
	nodes[newNode.Key] = newNode
	parentNode := nodes[newNode.Data.ParentURL]
	if parentNode == nil {
		parentNode = &Node{
			Key:        "",
			Neighbours: []*Node{},
			Data:       &Proposition{},
		}
		nodes[parentNode.Key] = parentNode
	}
	parentNode.Neighbours = append(parentNode.Neighbours, newNode)

	return nodes
}

func outputData(outputName string, propositions []*Proposition) error {
	mapped := [][]string{{"Proposition ID", "Proposition Full Name", "Proposition", "Proposition Code", "URL"}}
	for _, p := range propositions {
		mapped = append(mapped, p.ToPropertySlice())
	}
	var outputBuffer bytes.Buffer
	csvWriter := csv.NewWriter(&outputBuffer)
	err := csvWriter.WriteAll(mapped)
	if err != nil {
		return errors.Wrap(err, "unable to write csv data to buffer")
	}
	csvWriter.Flush()

	err = ioutil.WriteFile(outputName, outputBuffer.Bytes(), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "unable to write csv data to file")
	}

	return nil
}

func processLink(r *colly.Response) (*Proposition, error) {
	parent := r.Ctx.Get("parent")
	labels, _ := getLabels(r)
	fmt.Printf("Visiting %v\tParent %v\n", r.Request.URL, parent)
	id, _ := createID()
	return &Proposition{
		ID:            id,
		Tag:           labels[0],
		PotentialTags: labels,
		URL:           r.Request.URL.String(),
		Key:           r.Request.URL.String(),
		ParentURL:     parent,
	}, nil
}

func getLabels(r *colly.Response) ([]string, error) {
	// very quick and dirty title grab
	body := string(r.Body)
	title := strings.TrimSpace(body[strings.Index(body, "<title>")+7 : strings.Index(body, "</title>")])
	linkText := strings.TrimSpace(r.Ctx.Get("link"))

	return []string{title, linkText}, nil
}

func createID() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
