package config

import (
	"slices"
)

// mergeTree 构建一颗类似节点数，负责按照 diffDegree 指定的级别合并节点
type mergeTree struct {
	root       *trieMergeNode
	diffDegree int
}

func newMergeTree(diffDegree int) *mergeTree {
	return &mergeTree{
		root:       nil,
		diffDegree: diffDegree,
	}
}

type trieMergeNode struct {
	name     string
	parent   *trieMergeNode
	children []*trieMergeNode
}

func (node *trieMergeNode) addChild(child *trieMergeNode) {
	node.children = append(node.children, child)
}

func newTrieNode(name string, parent *trieMergeNode) *trieMergeNode {
	return &trieMergeNode{
		name:     name,
		parent:   parent,
		children: make([]*trieMergeNode, 0),
	}
}

// AddNodes 将 nodes 添加到树中
func (t *mergeTree) addNodes(nodes []string) {
	if len(nodes) == 0 {
		return
	}

	if t.root == nil {
		t.root = newTrieNode("/", nil)
	}

	cur := t.root
	for _, n := range nodes {

		exists := false
		for _, c := range cur.children {
			exists = c.name == n
			if exists {
				cur = c
				break
			}
		}

		if exists {
			// move to next step
			continue
		}

		newNode := newTrieNode(n, cur)
		cur.addChild(newNode)
		cur = newNode
	}
}

// LookCommonPath 在 diffDegree 所允许的差异层级中查找最大公共前缀节点
func (t *mergeTree) lookCommonPath(_ int) [][]string {

	watchPaths := make([][]string, 0)
	if len(t.root.children) == 0 {
		return watchPaths
	}

	cur := t.root.children[0]
	walkUp := true

	for cur != nil && cur != t.root {

		walkedNode, paths := t.walkOnNode(cur, walkUp)
		watchPaths = append(watchPaths, paths)

		// find sibling node
		cur = t.findSibling(walkedNode)
		if cur == nil {
			break
		}
		walkUp = cur.parent != walkedNode.parent
	}
	return watchPaths
}

// findSibling 查找 node 的下一个兄弟节点，如果兄弟节点不存在 (如只有一个节点或已经是最后一个节点)
// 则查找父节点的兄弟节点，如果没有匹配的节点，则一直递归到根节点
func (t *mergeTree) findSibling(node *trieMergeNode) *trieMergeNode {
	if node == nil || node == t.root {
		return nil
	}

	for i, n := range node.parent.children {
		if n == node {
			if len(node.parent.children) == i+1 { // last one
				return t.findSibling(node.parent)
			}

			return node.parent.children[i+1]
		}
	}
	return nil
}

// WalkOnNode 从 node 开始，遍历到最深叶子节点，上升到 diffDegree 父节点，
// 如果子节点只有一个，继续往下找到叶子或有多个子节点的节点，建立监听
// 如果子节点超过一个，则直接监听当前节点
func (t *mergeTree) walkOnNode(node *trieMergeNode, walkUp bool) (*trieMergeNode, []string) {
	leaf := node
	for len(leaf.children) > 0 {
		leaf = leaf.children[0]
	}

	parent := leaf
	for diff := t.diffDegree; walkUp && diff > 0 && parent.parent != nil && parent.parent != t.root; diff-- {
		parent = parent.parent
	}

	for walkUp && len(parent.children) < 2 && len(parent.children) > 0 {
		parent = parent.children[0]
	}

	path := make([]string, 0)
	resultNode := parent
	for parent.parent != nil {
		path = append(path, parent.name)
		parent = parent.parent
	}

	slices.Reverse(path)
	return resultNode, path
}
