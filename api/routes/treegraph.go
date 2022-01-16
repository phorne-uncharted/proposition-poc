package routes

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	log "github.com/unchartedsoftware/plog"
)

// TreeGraphItem is an item in the treegraph structure.
type TreeGraphItem struct {
	ID    string `json:"id"`
	Value int    `json:"value,omitempty"`
}

// TreeGraph is a transformed graph to match the expected treegraph structure.
type TreeGraph struct {
	Items []*TreeGraphItem `json:"items"`
}

func buildTreeGraph(graph *Graph, maxDepth int) *TreeGraph {
	return &TreeGraph{nodeToGraphItem(map[string]bool{}, maxDepth, 1, graph.Root)}
}

func nodeToGraphItem(ids map[string]bool, maxDepth int, depth int, node *Node) []*TreeGraphItem {
	item := &TreeGraphItem{
		ID: node.Data.FullName[:len(node.Data.FullName)-1],
	}

	items := []*TreeGraphItem{}
	if !ids[node.Data.FullName] {
		items = append(items, item)
		ids[node.Data.FullName] = true
	}

	if len(node.Neighbours) == 0 {
		item.Value = 1
	}

	children := []*TreeGraphItem{}
	if depth < maxDepth {
		for _, c := range node.Neighbours {
			children = append(children, nodeToGraphItem(ids, maxDepth, depth+1, c)...)
		}
	}

	return append(items, children...)
}

// TreeGraphHandler generates a route handler that returns a treegraph structure.
func TreeGraphHandler(allowedSites []string) func(http.ResponseWriter, *http.Request) {
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
		maxDepth := int(params["maxDepth"].(float64))
		log.Infof("starting processing of site '%s' to a max depth of %d", urlRaw, maxDepth)

		urlParsed, err := url.Parse(urlRaw)
		if err != nil {
			handleError(w, errors.Wrap(err, "Unable to parse url"))
			return
		}

		if !allowedSitesMap[urlParsed.Hostname()] {
			handleError(w, errors.Errorf("host '%s' not allowed", urlParsed.Hostname()))
			return
		}

		graph := processTreegraph(urlParsed.String(), buildNodes(urlParsed))
		treemap := buildTreeGraph(graph, maxDepth)

		// marshal data
		err = handleJSON(w, treemap)
		if err != nil {
			handleError(w, errors.Wrap(err, "unable to marshal model result into JSON"))
			return
		}
	}
}

func processTreegraph(url string, nodes map[string]*Node) *Graph {
	root := nodes[""]
	root.Key = url
	root.Data = &Proposition{
		URL:      url,
		ID:       url,
		FullName: "HOME.",
	}

	for i, c := range root.Neighbours {
		root.Neighbours[i] = buildBreadCrumb(root, c, ".", []*Node{})[0]
	}

	return &Graph{
		URL:  url,
		Root: root,
	}
}
