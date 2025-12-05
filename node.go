package tdms

import "github.com/egregors/sortedmap"

type Node struct {
	name       string
	path       string
	properties *sortedmap.SortedMap[map[string]any, string, any]
	childMap   *sortedmap.SortedMap[map[string]*Node, string, *Node]
}

func NewNode(name string, path string) *Node {
	return &Node{
		name: name,
		path: path,
		properties: sortedmap.New[map[string]any, string, any](func(i, j sortedmap.KV[string, any]) bool {
			return i.Key < j.Key
		}),
		childMap: sortedmap.New[map[string]*Node, string, *Node](func(i, j sortedmap.KV[string, *Node]) bool {
			return i.Key < j.Key
		}),
	}
}

func (node *Node) Name() string {
	return node.name
}

func (node *Node) Path() string {
	return node.path
}

func (node *Node) Properties() *sortedmap.SortedMap[map[string]any, string, any] {
	return node.properties
}

func (node *Node) Children() []*Node {
	return node.childMap.CollectValues()
}

func (node *Node) GetChildByName(name string) *Node {
	child, _ := node.childMap.Get(name)
	return child
}

func (node *Node) AddChild(child *Node) {
	node.childMap.Insert(child.Name(), child)
}
