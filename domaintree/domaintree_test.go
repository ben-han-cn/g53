package domaintree

import (
	ut "cement/unittest"
	//"fmt"
	"testing"

	"g53"
)

/* The initial structure of rbtree
 *
 *             b
 *           /   \
 *          a    d.e.f
 *              /  |   \
 *             c   |    g.h
 *                 |     |
 *                w.y    i
 *              /  |  \
 *             x   |   z
 *                 |   |
 *                 p   j
 *               /   \
 *              o     q
 */
func nameFromString(n string) *g53.Name {
	name, _ := g53.NameFromString(n)
	return name
}

func treeInsertString(tree *DomainTree, n string) (*Node, error) {
	return tree.Insert(nameFromString(n))
}

func createDomainTree(returnEmptyNode bool) *DomainTree {
	domains := []string{
		"c", "b", "a", "x.d.e.f", "z.d.e.f", "g.h", "i.g.h", "o.w.y.d.e.f",
		"j.z.d.e.f", "p.w.y.d.e.f", "q.w.y.d.e.f"}

	tree := NewDomainTree(returnEmptyNode)
	for i, d := range domains {
		node, _ := treeInsertString(tree, d)
		node.data = i + 1
	}
	return tree
}

func TestTreeNodeCount(t *testing.T) {
	ut.Equal(t, createDomainTree(false).nodeCount, 13)
}

func TestTreeInsert(t *testing.T) {
	tree := createDomainTree(false)
	node, err := treeInsertString(tree, "c")
	ut.Equal(t, err, ErrAlreadyExist)

	node, err = treeInsertString(tree, "d.e.f")
	ut.Equal(t, err, nil)
	ut.Equal(t, node.name.String(true), "d.e.f")
	ut.Equal(t, tree.nodeCount, 13)

	node, err = tree.Insert(g53.Root)
	ut.Assert(t, err == nil, "inert root domain should ok but get %v", err)
	ut.Assert(t, g53.Root.Equals(node.name), "insert return node name should equals to insert name")
	ut.Equal(t, tree.nodeCount, 14)

	node, err = treeInsertString(tree, "example.com")
	ut.Assert(t, err == nil, "inert new domain should ok but get %v", err)
	ut.Equal(t, tree.nodeCount, 15)
	node.data = 12

	node, err = treeInsertString(tree, "example.com")
	ut.Equal(t, err, ErrAlreadyExist)
	ut.Equal(t, node.name.String(true), "example.com")
	ut.Equal(t, tree.nodeCount, 15)

	// split the node "d.e.f"
	node, err = treeInsertString(tree, "k.e.f")
	ut.Equal(t, node.name.String(true), "k")
	ut.Equal(t, tree.nodeCount, 17)

	// split the node "g.h"
	node, err = treeInsertString(tree, "h")
	ut.Equal(t, err, nil)
	ut.Equal(t, node.name.String(true), "h")
	ut.Equal(t, tree.nodeCount, 18)

	// add child domain
	node, err = treeInsertString(tree, "m.p.w.y.d.e.f")
	ut.Equal(t, node.name.String(true), "m")
	ut.Equal(t, tree.nodeCount, 19)

	node, err = treeInsertString(tree, "n.p.w.y.d.e.f")
	ut.Assert(t, err == nil, "insert new child name should ok but get %v", err)
	ut.Equal(t, node.name.String(true), "n")
	ut.Equal(t, tree.nodeCount, 20)

	node, err = treeInsertString(tree, "l.a")
	ut.Equal(t, node.name.String(true), "l")
	ut.Equal(t, tree.nodeCount, 21)

	_, err = treeInsertString(tree, "r.d.e.f")
	ut.Assert(t, err == nil, "insert new child name should ok but get %v", err)
	_, err = treeInsertString(tree, "s.d.e.f")
	ut.Assert(t, err == nil, "insert new child name should ok but get %v", err)
	ut.Equal(t, tree.nodeCount, 23)
	_, err = treeInsertString(tree, "h.w.y.d.e.f")
	ut.Assert(t, err == nil, "insert new child name should ok but get %v", err)

	node, err = treeInsertString(tree, "f")
	ut.Assert(t, err == nil, "f node has no data")
	node.SetData(1000)
	_, err = treeInsertString(tree, "f")
	ut.Assert(t, err == ErrAlreadyExist, "insert already exists domain should get error")

	newNames := []string{"m", "nm", "om", "k", "l", "fe", "ge", "i", "ae", "n"}
	for _, newName := range newNames {
		_, err = treeInsertString(tree, newName)
		ut.Assert(t, err == nil, "insert new child name should ok but get %v", err)
	}
}

func TestTreeSearch(t *testing.T) {
	tree := createDomainTree(false)
	node, ret := tree.Search(nameFromString("a"))
	ut.Equal(t, ret, ExactMatch)
	ut.Equal(t, node.name.String(true), "a")

	notExistsNames := []string{
		"d.e.f", "y.d.e.f", "x", "m.n",
	}
	for _, n := range notExistsNames {
		_, ret := tree.Search(nameFromString(n))
		ut.Equal(t, ret, NotFound)
	}

	tree = createDomainTree(true)
	exactMatchNames := []string{
		"d.e.f", "w.y.d.e.f",
	}
	for _, n := range exactMatchNames {
		_, ret := tree.Search(nameFromString(n))
		ut.Equal(t, ret, ExactMatch)
	}

	// partial match
	node, ret = tree.Search(nameFromString("m.b"))
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, node.name.String(true), "b")
	node, ret = tree.Search(nameFromString("m.d.e.f"))
	ut.Equal(t, ret, PartialMatch)

	// find rbtnode
	node, ret = tree.Search(nameFromString("q.w.y.d.e.f"))
	ut.Equal(t, ret, ExactMatch)
	ut.Equal(t, node.name.String(true), "q")
}

func TestTreeFlag(t *testing.T) {
	tree := createDomainTree(false)
	node, ret := treeInsertString(tree, "flags.example")
	ut.Equal(t, ret, nil)
	ut.Assert(t, node.GetFlag(NF_CALLBACK) == false, "by default, node has no flag")
	node.SetFlag(NF_CALLBACK, true)
	ut.Assert(t, node.GetFlag(NF_CALLBACK) == true, "node should has flag after set")
	node.SetFlag(NF_CALLBACK, false)
	ut.Assert(t, node.GetFlag(NF_CALLBACK) == false, "node should has no flag after reset")
}

func testCallback(node *Node, callbackChecker interface{}) bool {
	*(callbackChecker.(*bool)) = true
	return false
}

func TestTreeNodeCallback(t *testing.T) {
	tree := createDomainTree(false)
	node, err := treeInsertString(tree, "callback.example")
	ut.Equal(t, err, nil)
	node.data = 1
	ut.Assert(t, node.GetFlag(NF_CALLBACK) == false, "by default, node has no flag")
	node.SetFlag(NF_CALLBACK, true)
	// add more levels below and above the callback node for partial match.

	subNode, err := treeInsertString(tree, "sub.callback.example")
	ut.Equal(t, err, nil)
	node.data = 2
	parentNode, _ := treeInsertString(tree, "example")
	node, ret := tree.Search(nameFromString("callback.example"))
	ut.Assert(t, node.GetFlag(NF_CALLBACK) == true, "node has set flag")
	ut.Assert(t, subNode.GetFlag(NF_CALLBACK) == false, "node hasn't set flag")
	ut.Assert(t, parentNode.GetFlag(NF_CALLBACK) == false, "node hasn't  set flag")

	// check if the callback is called from find()
	nodePath := NewNodeChain()
	callbackCalled := false
	node, ret = tree.SearchExt(nameFromString("sub.callback.example"), nodePath, testCallback, &callbackCalled)
	ut.Equal(t, callbackCalled, true)

	// enable callback at the parent node, but it doesn't have data so
	// the callback shouldn't be called.
	nodePath2 := NewNodeChain()
	parentNode.SetFlag(NF_CALLBACK, true)
	callbackCalled = false
	node, ret = tree.SearchExt(nameFromString("callback.example"), nodePath2, testCallback, &callbackCalled)
	ut.Equal(t, ret, ExactMatch)
	ut.Equal(t, callbackCalled, false)
}

func TestTreeNodeChain(t *testing.T) {
	chain := NewNodeChain()
	ut.Equal(t, chain.GetLevelCount(), 0)

	tree := NewDomainTree(true)
	treeInsertString(tree, ".")
	_, ret := tree.SearchExt(nameFromString("."), chain, nil, nil)
	ut.Equal(t, ret, ExactMatch)
	ut.Equal(t, chain.GetLevelCount(), 1)

	/*
	 * Now creating a possibly deepest tree with MAX_LABELS levels.
	 * it should look like:
	 *           (.)
	 *            |
	 *            a
	 *            |
	 *            a
	 *            : (MAX_LABELS - 1) "a"'s
	 *
	 * then confirm that find() for the deepest name succeeds without any
	 * disruption, and the resulting chain has the expected level.
	 * Note that the root name (".") solely belongs to a single level,
	 * so the levels begin with 2.
	 */
	nodeName := g53.Root
	for i := 2; i <= g53.MAX_LABELS; i++ {
		nodeName, _ = nameFromString("a").Concat(nodeName)
		_, err := tree.Insert(nodeName)
		ut.Equal(t, err, nil)

		chain := NewNodeChain()
		_, ret := tree.SearchExt(nodeName, chain, nil, nil)
		ut.Equal(t, ret, ExactMatch)
		ut.Equal(t, chain.GetLevelCount(), i)
	}
}

//
//the domain order should be:
// a, b, c, d.e.f, x.d.e.f, w.y.d.e.f, o.w.y.d.e.f, p.w.y.d.e.f, q.w.y.d.e.f,
// z.d.e.f, j.z.d.e.f, g.h, i.g.h
//             b
//           /   \
//          a    d.e.f
//              /  |   \
//             c   |    g.h
//                 |     |
//                w.y    i
//              /  |  \
//             x   |   z
//                 |   |
//                 p   j
//               /   \
//              o     q
///
func TestTreeNextNode(t *testing.T) {
	names := []string{
		"a", "b", "c", "d.e.f", "x.d.e.f", "w.y.d.e.f", "o.w.y.d.e.f",
		"p.w.y.d.e.f", "q.w.y.d.e.f", "z.d.e.f", "j.z.d.e.f", "g.h", "i.g.h"}
	tree := createDomainTree(false)
	nodePath := NewNodeChain()
	node, ret := tree.SearchExt(nameFromString(names[0]), nodePath, nil, nil)
	ut.Equal(t, ret, ExactMatch)
	for i := 0; i < len(names); i++ {
		ut.Assert(t, node != nil, "node shouldn't be nil")
		ut.Equal(t, names[i], nodePath.GetAbsoluteName().String(true))
		node = tree.nextNode(nodePath)
	}

	// We should have reached the end of the tree.
	ut.Assert(t, node == nil, "node will reach the end")
}

func comparisonChecks(t *testing.T, chain *NodeChain, expectedOrder int, expectedCommonLabels int, expectedRelation g53.NameRelation) {
	if expectedOrder > 0 {
		ut.Assert(t, chain.lastComparison.Order > 0, "")
	} else if expectedOrder < 0 {
		ut.Assert(t, chain.lastComparison.Order < 0, "")
	} else {
		ut.Equal(t, chain.lastComparison.Order, 0)
	}

	ut.Equal(t, expectedCommonLabels, chain.lastComparison.CommonLabelCount)
	ut.Equal(t, expectedRelation, chain.lastComparison.Relation)
}

func TestTreeNodeChainLastComparison(t *testing.T) {
	chain := NewNodeChain()
	ut.Equal(t, chain.lastCompared, (*Node)(nil))

	emptyTree := NewDomainTree(false)
	node, ret := emptyTree.SearchExt(nameFromString("a"), chain, nil, nil)
	ut.Equal(t, ret, NotFound)
	ut.Equal(t, chain.lastCompared, (*Node)(nil))
	chain.clear()

	tree := createDomainTree(true)
	node, _ = tree.SearchExt(nameFromString("x.d.e.f"), chain, nil, nil)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, 0, 2, g53.EQUAL)
	chain.clear()

	_, ret = tree.Search(nameFromString("i.g.h"))
	ut.Equal(t, ret, ExactMatch)
	node, ret = tree.SearchExt(nameFromString("x.i.g.h"), chain, nil, nil)
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, 1, 2, g53.SUBDOMAIN)
	chain.clear()

	// Partial match, search stopped in the subtree below the matching node
	// after following a left branch.
	node, ret = tree.Search(nameFromString("x.d.e.f"))
	ut.Equal(t, ret, ExactMatch)
	_, ret = tree.SearchExt(nameFromString("a.d.e.f"), chain, nil, nil)
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, -1, 1, g53.COMMONANCESTOR)
	chain.clear()

	// Partial match, search stopped in the subtree below the matching node
	// after following a right branch.
	node, ret = tree.Search(nameFromString("z.d.e.f"))
	ut.Equal(t, ret, ExactMatch)
	_, ret = tree.SearchExt(nameFromString("zz.d.e.f"), chain, nil, nil)
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, 1, 1, g53.COMMONANCESTOR)
	chain.clear()

	// Partial match, search stopped at a node for a super domain of the
	// search name in the subtree below the matching node.
	node, ret = tree.Search(nameFromString("w.y.d.e.f"))
	ut.Equal(t, ret, ExactMatch)
	_, ret = tree.SearchExt(nameFromString("y.d.e.f"), chain, nil, nil)
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, -1, 2, g53.SUPERDOMAIN)
	chain.clear()

	// Partial match, search stopped at a node that share a common ancestor
	// with the search name in the subtree below the matching node.
	// (the expected node is the same as the previous case)
	_, ret = tree.SearchExt(nameFromString("z.y.d.e.f"), chain, nil, nil)
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, 1, 2, g53.COMMONANCESTOR)
	chain.clear()

	// Search stops in the highest level after following a left branch.
	node, ret = tree.Search(nameFromString("c"))
	ut.Equal(t, ret, ExactMatch)
	_, ret = tree.SearchExt(nameFromString("bb"), chain, nil, nil)
	ut.Equal(t, ret, NotFound)
	//ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, -1, 1, g53.COMMONANCESTOR)
	chain.clear()

	// Search stops in the highest level after following a right branch.
	// (the expected node is the same as the previous case)
	_, ret = tree.SearchExt(nameFromString("d"), chain, nil, nil)
	ut.Equal(t, ret, NotFound)
	ut.Equal(t, chain.lastCompared, node)
	comparisonChecks(t, chain, 1, 1, g53.COMMONANCESTOR)
	chain.clear()
}

func TestRootZone(t *testing.T) {
	tree := NewDomainTree(false)
	node, _ := treeInsertString(tree, ".")
	node.SetData(1)

	node, ret := tree.Search(nameFromString("."))
	ut.Equal(t, ret, ExactMatch)

	node, ret = tree.Search(nameFromString("example.com"))
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, node.name.String(true), ".")
	ut.Equal(t, node.Data().(int), 1)

	node, _ = treeInsertString(tree, "com")
	node.SetData(2)
	node, ret = tree.Search(nameFromString("example.com"))
	ut.Equal(t, ret, PartialMatch)
	ut.Equal(t, node.name.String(true), "com")
	ut.Equal(t, node.Data().(int), 2)
}
