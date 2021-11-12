package dynamicmerkletree

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"testing"
)

type TestSHA256Content struct {
	x string
}

//CalculateHash hashes the values of a TestSHA256Content
func (t TestSHA256Content) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (t TestSHA256Content) Equals(other Content) (bool, error) {
	return t.x == other.(TestSHA256Content).x, nil
}

var table = []struct {
	testCaseId          int
	hashStrategy        func() hash.Hash
	hashStrategyName    string
	defaultHashStrategy bool
	contents            []Content
	expectedHash        []byte
	notInContents       Content
	nodeCntOfLeaf       int
}{
	{
		testCaseId:          0,
		hashStrategy:        sha256.New,
		hashStrategyName:    "sha256",
		defaultHashStrategy: true,
		contents: []Content{
			TestSHA256Content{
				x: "Hello",
			},
			TestSHA256Content{
				x: "Hi",
			},
			TestSHA256Content{
				x: "Hey",
			},
			TestSHA256Content{
				x: "Hola",
			},
		},
		notInContents: TestSHA256Content{x: "NotInTestTable"},
		expectedHash:  []byte{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188},
		nodeCntOfLeaf: 2,
	},
	{
		testCaseId:          1,
		hashStrategy:        sha256.New,
		hashStrategyName:    "sha256",
		defaultHashStrategy: true,
		contents: []Content{
			TestSHA256Content{
				x: "Hello",
			},
			TestSHA256Content{
				x: "Hi",
			},
			TestSHA256Content{
				x: "Hey",
			},
			TestSHA256Content{
				x: "Hola",
			},
			TestSHA256Content{
				x: "BABABA",
			},
			TestSHA256Content{
				x: "LALALA",
			},
		},
		notInContents: TestSHA256Content{x: "NotInTestTable"},
		expectedHash:  []byte{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188},
		nodeCntOfLeaf: 2,
	},
	{
		testCaseId:          2,
		hashStrategy:        sha256.New,
		hashStrategyName:    "sha256",
		defaultHashStrategy: true,
		contents: []Content{
			TestSHA256Content{
				x: "Hello",
			},
			TestSHA256Content{
				x: "Hi",
			},
			TestSHA256Content{
				x: "Hey",
			},
			TestSHA256Content{
				x: "Hola Hola",
			},
			TestSHA256Content{
				x: "BABABA",
			},
			TestSHA256Content{
				x: "LALALA",
			},
		},
		notInContents: TestSHA256Content{x: "NotInTestTable"},
		expectedHash:  []byte{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188},
		nodeCntOfLeaf: 2,
	},
	{
		testCaseId:          3,
		hashStrategy:        sha256.New,
		hashStrategyName:    "sha256",
		defaultHashStrategy: true,
		contents: []Content{
			TestSHA256Content{
				x: "Hello",
			},
			TestSHA256Content{
				x: "Hi",
			},
			TestSHA256Content{
				x: "Hey",
			},
			TestSHA256Content{
				x: "Hola Hola",
			},
			TestSHA256Content{
				x: "BABABA",
			},
			TestSHA256Content{
				x: "LALALA",
			},
			TestSHA256Content{
				x: "BACADA",
			},
		},
		notInContents: TestSHA256Content{x: "NotInTestTable"},
		expectedHash:  []byte{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188},
		nodeCntOfLeaf: 3,
	},
}

func TestNewTree(t *testing.T) {

	tree, err := NewTree(table[0].contents, table[0].nodeCntOfLeaf)
	if err != nil {
		t.Errorf("[case:%d] error: unexpected error: %v", table[0].testCaseId, err)
	}
	if !bytes.Equal(tree.MerkleRoot(), table[0].expectedHash) {
		t.Errorf("[case:%d] error: expected hash equal to %v got %v", table[0].testCaseId, table[0].expectedHash, tree.MerkleRoot())
	}

	tree_3, err := NewTree(table[3].contents, table[3].nodeCntOfLeaf)
	if err != nil {
		t.Errorf("[case:%d] error: unexpected error: %v", table[3].testCaseId, err)
	}
	if tree_3.MerkleRoot() == nil {
		t.Errorf("[case:%d] error: root nil error: %v", table[3].testCaseId, err)
	}
}

func TestAppendContent(t *testing.T) {
	tree_0, err := NewTree(table[0].contents, table[0].nodeCntOfLeaf)
	if err != nil {
		t.Errorf("[case:%d] error: unexpected error: %v", table[0].testCaseId, err)
	}

	appendContent := []Content{
		TestSHA256Content{
			x: "BABABA",
		},
		TestSHA256Content{
			x: "LALALA",
		},
	}
	tree_0.AppendContent(appendContent)

	tree_1, err := NewTree(table[1].contents, table[1].nodeCntOfLeaf)
	if err != nil {
		t.Errorf("[case:%d] error: unexpected error: %v", table[1].testCaseId, err)
	}

	if !bytes.Equal(tree_0.MerkleRoot(), tree_1.MerkleRoot()) {
		t.Errorf("[case:%d] error: expected hash equal to %v got %v", table[1].testCaseId, tree_0.MerkleRoot(), tree_1.MerkleRoot())
	}

}

func TestUpdateContent(t *testing.T) {
	tree_1, err := NewTree(table[1].contents, table[1].nodeCntOfLeaf)
	if err != nil {
		t.Errorf("[case:%d] error: unexpected error: %v", table[1].testCaseId, err)
	}

	tree_1.UpdateContent(3, TestSHA256Content{
		x: "Hola Hola",
	})

	tree_2, err := NewTree(table[2].contents, table[2].nodeCntOfLeaf)
	if err != nil {
		t.Errorf("[case:%d] error: unexpected error: %v", table[2].testCaseId, err)
	}

	if !bytes.Equal(tree_1.MerkleRoot(), tree_2.MerkleRoot()) {
		t.Errorf("[case:%d] error: expected hash equal to %v got %v", table[1].testCaseId, tree_1.MerkleRoot(), tree_2.MerkleRoot())
	}
}
