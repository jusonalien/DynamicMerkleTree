package dynamicmerkletree

import (
	"crypto/sha256"
	"errors"
	"hash"
	"sync"
)

type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
}

type DMerkleTree struct {
	sync.RWMutex
	Root           *Node
	merkleRootHash []byte
	Leafs          []*Node
	NodeCntOfLeaf  int
	hashStrategy   func() hash.Hash
}

type Node struct {
	Tree      *DMerkleTree
	Parent    *Node
	LeafNodes []*Node
	leaf      bool
	dup       bool
	Hash      []byte
	C         Content
}

func buildWithContent(cs []Content, t *DMerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var leafs []*Node
	for _, c := range cs {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Hash: hash,
			C:    c,
			leaf: true,
			Tree: t,
		})
	}
	if len(leafs)%t.NodeCntOfLeaf != 0 {
		dupNum := t.NodeCntOfLeaf - (len(leafs) % t.NodeCntOfLeaf)
		for i := 0; i < dupNum; i++ {
			duplicate := &Node{
				Hash: leafs[len(leafs)-1].Hash,
				C:    leafs[len(leafs)-1].C,
				leaf: true,
				dup:  true,
				Tree: t,
			}
			leafs = append(leafs, duplicate)
		}
	}
	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

func buildIntermediate(nl []*Node, t *DMerkleTree) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += t.NodeCntOfLeaf {
		var subLeafs []*Node
		h := t.hashStrategy()
		var chash []byte
		for j := i; j < len(nl) && j < i+t.NodeCntOfLeaf; j += 1 {
			chash = append(chash, nl[j].Hash...)
			subLeafs = append(subLeafs, nl[j])
		}

		if _, err := h.Write(chash); err != nil {
			return nil, err
		}

		n := &Node{
			LeafNodes: subLeafs,
			Hash:      h.Sum(nil),
			Tree:      t,
		}
		for idx := 0; idx < len(subLeafs); idx++ {
			subLeafs[idx].Parent = n
		}

		nodes = append(nodes, n)
		if len(nl) <= t.NodeCntOfLeaf {
			return n, nil
		}
	}

	return buildIntermediate(nodes, t)
}

//NewTree creates a new Merkle Tree using the content cs.
func NewTree(cs []Content, nodeCntOfLeaf int) (*DMerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &DMerkleTree{
		hashStrategy:  defaultHashStrategy,
		NodeCntOfLeaf: nodeCntOfLeaf,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRootHash = root.Hash
	return t, nil
}

func NewTreeWithHashStrategy(cs []Content, nodeCntOfLeaf int, hashStrategy func() hash.Hash) (*DMerkleTree, error) {
	t := &DMerkleTree{
		hashStrategy:  hashStrategy,
		NodeCntOfLeaf: nodeCntOfLeaf,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRootHash = root.Hash
	return t, nil
}

//MerkleRoot returns the unverified Merkle Root (hash of the root node) of the tree.
func (m *DMerkleTree) MerkleRoot() []byte {

	m.RLock()
	defer m.RUnlock()

	return m.merkleRootHash
}

func (m *DMerkleTree) rebuildTreeWithContent(cs []Content) error {

	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRootHash = root.Hash
	return nil
}

func (m *DMerkleTree) AppendContent(appendContent []Content) error {

	m.Lock()
	defer m.Unlock()

	var content []Content
	for _, c := range m.Leafs {
		content = append(content, c.C)
	}

	content = append(content, appendContent...)

	err := m.rebuildTreeWithContent(content)
	if err != nil {
		return err
	}
	return nil
}

func (m *DMerkleTree) UpdateContent(targetIdx int, targetContent Content) error {

	m.Lock()
	defer m.Unlock()

	m.Leafs[targetIdx].C = targetContent
	var content []Content
	for _, c := range m.Leafs {
		content = append(content, c.C)
	}

	err := m.rebuildTreeWithContent(content)
	if err != nil {
		return err
	}
	return nil
}
