package rest

import (
	"fmt"
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
	Method string
	//Base     string //eg. /user
	//Path     string //eg. /id
	IsFilter bool
	Handler  func(req *Request, res *Response, next func(e error))
	Order    int
}

type PathTreeNode struct {
	Name     string
	Type     int
	Parent   *PathTreeNode
	Children []*PathTreeNode
	Handlers []*OrderHandler
}

func (this *PathTreeNode) Equal(node *PathTreeNode) bool {
	return this.Name == node.Name && this.Type == node.Type
}

func CreatePathNode(name string) *PathTreeNode {
	//if strings.HasPrefix(name, PathTreeNodeGreedyPrefix) {
	//	return &PathTreeNode{Name: strings.TrimPrefix(name, PathTreeNodeGreedyPrefix), Type: PathTreeNodeGreedy}
	//}
	if strings.HasPrefix(name, PathTreeNodeNonGreedyPrefix) {
		return &PathTreeNode{Name: strings.TrimPrefix(name, PathTreeNodeNonGreedyPrefix), Type: PathTreeNodeNonGreedy}
	}
	return &PathTreeNode{Name: name, Type: PathTreeNodeNorm}
}

func ParsePath(path string) []*PathTreeNode {
	ps := strings.Split(path, "/")
	nodes := make([]*PathTreeNode, len(ps))
	for i, p := range ps {
		nodes[i] = CreatePathNode(p)
	}
	for i, _ := range nodes {
		if i > 0 {
			nodes[i].Append(nodes[i-1])
		}
	}
	return nodes
}

func (this *PathTreeNode) Matches(path string) bool {
	return this.Name == path || this.Type == PathTreeNodeNonGreedy
}
func (this *PathTreeNode) String() string {
	if PathTreeNodeNorm == this.Type {
		return this.Name
	}
	if PathTreeNodeNonGreedy == this.Type {
		return PathTreeNodeNonGreedyPrefix + this.Name
	}
	return ""
}

//func (this *PathTreeNode) Get(path string) PathTreeNode {

//}

//func (this *PathTreeNode) Find(path string) []PathTreeNode {

//}
//func (this *PathTreeNode) FindHandlers() []OrderHandler {

//}
func (this *PathTreeNode) Append(node *PathTreeNode) {
	this.Children = append(this.Children, node)
}

func (this *PathTreeNode) AppendPath(name string) *PathTreeNode {
	fmt.Printf("AppendPath  %s->%s\n", this.Name, name)
	//var node = CreatePathNode(name);
	for _, child := range this.Children {
		if this.String() == name {
			return child
		}
	}

	node := CreatePathNode(name)
	this.Children = append(this.Children, node)
	fmt.Printf("create node %s , child count %d ,node.Name:%s\n", name, len(this.Children), node.Name)
	return node
}

func (this *PathTreeNode) Mount(path string, method string, isFilter bool, handler func(req *Request, res *Response, next func(e error))) *PathTreeNode {
	names := strings.Split(path, "/")
	fmt.Println(names)
	node := this
	for _, name := range names {
		if 0 == len(name) {
			continue
		}
		node = node.AppendPath(name)
	}
	node.Handlers = append(this.Handlers, &OrderHandler{Method: strings.ToUpper(method), Handler: handler, IsFilter: isFilter})
	return this
}
func (this *PathTreeNode) Root() *PathTreeNode {
	var node *PathTreeNode
	for node = this; nil != node.Parent; node = node.Parent {
	}
	return node
}
func (this *PathTreeNode) MaxOrder() int {
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
