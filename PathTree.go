package rest

import (
	"strings"
)

const (
	PathTreeNodeNorm      = iota //eg /user
	PathTreeNodeNonGreedy        //eg. /user/:id
	PathTreeNodeGreedy           //eg. /static/::path
)
const (
	PathTreeNodeNonGreedyPrefix = ":"  //eg. /user/:id
	PathTreeNodeGreedyPrefix    = "::" //eg. /static/::path
)

type OrderHandler struct {
	Method   string
	Base     string //eg. /user
	Path     string //eg. /id
	IsFilter bool
	Handle   func(req *Request, res *Response, next func(e error))
	Order    int
}

type PathTreeNode struct {
	Name     string
	Type     int
	Parent   *PathTreeNode
	Children []*PathTreeNode
	Handlers []*OrderHandler
}

func ParsePathNode(raw string) PathTreeNode {
	if strings.HasPrefix(raw, PathTreeNodeGreedyPrefix) {
		return &PathTreeNode{Name: strings.TrimPrefix(raw, PathTreeNodeGreedyPrefix), Type: PathTreeNodeGreedy}
	}
	if strings.HasPrefix(raw, PathTreeNodeNonGreedyPrefix) {
		return &PathTreeNode{Name: strings.TrimPrefix(raw, PathTreeNodeNonGreedyPrefix), Type: PathTreeNodeNonGreedy}
	}
	return &PathTreeNode{Name: raw, Type: PathTreeNodeNorm}
}

func ParsePath(path string) *PathTreeNode {
	ps := strings.Split(path, "/")
	var root = &PathTreeNode(Name)
	var node *PathTreeNode

}
func (this *PathTreeNode) Match(path string) string {

}
func (this *PathTreeNode) Get(path string) PathTreeNode {

}

func (this *PathTreeNode) Find(path string) []PathTreeNode {

}
func (this *PathTreeNode) FindHandlers() []OrderHandler {

}
func (this *PathTreeNode) Append(node *PathTreeNode) {
	this.Children = append(this.Children, node)
}
func (this *PathTreeNode) Mount(path string, method string, Handle func(req *Request, res *Response, next func(e error))) {

}
func (this *PathTreeNode) Root() PathTreeNode {
	for node := this; nil != node.Parent; node = node.Parent {
	}
	return node
}
func (this *PathTreeNode) MaxOrder() PathTreeNode {
	order := 0
	this.Root().Walk(func(node *PathTreeNode) {
		order += len(node.Handlers)
	})
	return order
}
func (this *PathTreeNode) Walk(visitor func(node *PathTreeNode)) {
	visitor(this)
	if 0 >= len(this.Children) {
		return
	}
	for _, n := range this.Children {
		n.Walk(visitor)
	}
}
