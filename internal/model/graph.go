package model

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	State string `json:"state"`
}

type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Kind string `json:"kind"`
}

func NewGraph() Graph {
	return Graph{Nodes: make([]GraphNode, 0), Edges: make([]GraphEdge, 0)}
}

func (g *Graph) AddNode(node GraphNode) {
	g.Nodes = append(g.Nodes, node)
}

func (g *Graph) AddEdge(edge GraphEdge) {
	g.Edges = append(g.Edges, edge)
}

func (g Graph) HasNode(id string) bool {
	for _, node := range g.Nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}
