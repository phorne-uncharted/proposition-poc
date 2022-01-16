package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/unchartedsoftware/plog"
	"github.com/zenazn/goji/graceful"
	goji "goji.io/v3"
	"goji.io/v3/pat"

	"github.com/phorne-uncharted/proposition-poc/api/env"
	"github.com/phorne-uncharted/proposition-poc/api/middleware"
	"github.com/phorne-uncharted/proposition-poc/api/routes"
)

// Proposition is an entity being extracted from a site.
type Proposition struct {
	ID            string
	FullName      string
	Tag           string
	Code          string
	URL           string
	Key           string
	PotentialTags []string
	ParentURL     string
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

func registerRoute(mux *goji.Mux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Infof("Registering GET route %s", pattern)
	mux.HandleFunc(pat.Get(pattern), handler)
}

func registerRoutePost(mux *goji.Mux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Infof("Registering POST route %s", pattern)
	mux.HandleFunc(pat.Post(pattern), handler)
}

func main() {
	// load config from env
	config, err := env.LoadConfig()
	if err != nil {
		log.Errorf("%+v", err)
		os.Exit(1)
	}
	log.Infof("%+v", spew.Sdump(config))

	allowedSites, err := loadAllowedSites(config.AllowedSitesFile)
	if err != nil {
		log.Errorf("%+v", err)
		os.Exit(1)
	}

	// register routes
	mux := goji.NewMux()
	mux.Use(middleware.Log)
	mux.Use(middleware.Gzip)
	registerRoutePost(mux, "/site/treemap", routes.LinksHandler(allowedSites))
	registerRoutePost(mux, "/site/treegraph", routes.TreeGraphHandler(allowedSites))

	registerRoute(mux, "/*", routes.FileHandler("./dist"))

	// catch kill signals for graceful shutdown
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)

	// kick off the server listen loop
	log.Infof("Listening on port %s", config.AppPort)
	err = graceful.ListenAndServe(":"+config.AppPort, mux)
	if err != nil {
		log.Errorf("%+v", err)
		os.Exit(1)
	}

	// wait until server gracefully exits
	graceful.Wait()
}

func loadAllowedSites(filename string) ([]string, error) {
	log.Infof("loading allowed sites from file '%s'", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open allowed sites file")
	}
	defer file.Close()

	allowedSites := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allowedSites = append(allowedSites, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "unable to read allowed sites file")
	}
	log.Infof("allowed sites loaded from file: %v", allowedSites)

	return allowedSites, nil
}

func generate() {
	nodes := map[string]*node{}
	c := colly.NewCollector(
		colly.AllowedDomains("onedemo-telco.azurewebsites.net"),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Ctx.Put("parent", e.Request.URL.String())
		e.Request.Visit(e.Attr("href"))
	})

	//c.OnRequest(func(r *colly.Request) {
	//	r.Ctx.Put("parent", r.URL)
	//})

	c.OnResponse(func(r *colly.Response) {
		prop, _ := processLink(r)
		nodes = addNode(nodes, prop)
	})

	c.Visit("https://onedemo-telco.azurewebsites.net/")

	propositions := processGraph(nodes)
	outputData("output.csv", propositions)
}

func processGraph(nodes map[string]*node) []*Proposition {
	root := nodes[""]
	return buildBreadCrumb(nil, root, []*Proposition{})
}

func buildBreadCrumb(parent *node, current *node, alreadyProcessed []*Proposition) []*Proposition {
	if parent == nil {
		current.data.FullName = fmt.Sprintf("/")
	} else {
		current.data.FullName = fmt.Sprintf("%s%s/", parent.data.FullName, current.data.Tag)
		alreadyProcessed = append(alreadyProcessed, current.data)
	}

	for _, c := range current.children {
		alreadyProcessed = buildBreadCrumb(current, c, alreadyProcessed)
	}

	return alreadyProcessed
}

func addNode(nodes map[string]*node, proposition *Proposition) map[string]*node {
	newNode := &node{
		key:      proposition.URL,
		children: []*node{},
		data:     proposition,
	}
	nodes[newNode.key] = newNode
	parentNode := nodes[newNode.data.ParentURL]
	if parentNode == nil {
		parentNode = &node{
			key:      "",
			children: []*node{},
			data:     &Proposition{},
		}
		nodes[parentNode.key] = parentNode
	}
	parentNode.children = append(parentNode.children, newNode)

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

type node struct {
	key      string
	children []*node
	data     *Proposition
}
