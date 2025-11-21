package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"

	"proof-of-reserves/internal/proof"
)

// Result captures the outcome of a verification run.
type Result struct {
	RootHash string
	LeafHash string
	Levels   int
}

// Verify ensures that the provided proof forms a valid path from the self leaf
// to the advertised root.
func Verify(file *proof.File) (*Result, error) {
	if file == nil {
		return nil, errors.New("proof data is nil")
	}

	nodes := file.Nodes()
	grouped := make(map[int][]proof.Node)
	for _, node := range nodes {
		grouped[node.Level] = append(grouped[node.Level], node)
	}

	if len(grouped) == 0 {
		return nil, errors.New("proof contains no levels")
	}

	rootNodes := grouped[0]
	if len(rootNodes) == 0 {
		return nil, errors.New("missing root node in proof path")
	}

	levels := make([]int, 0, len(grouped))
	for level := range grouped {
		levels = append(levels, level)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(levels)))

	var rollingHash string
	for _, level := range levels {
		nodesAtLevel := grouped[level]

		if level == 0 {
			root := nodesAtLevel[0]
			if rollingHash != "" && !hashEqual(rollingHash, root.MerkelLeaf) {
				return nil, fmt.Errorf("computed root %s != provided %s", rollingHash, root.MerkelLeaf)
			}
			rollingHash = root.MerkelLeaf
			continue
		}

		left, right, err := siblingPair(nodesAtLevel)
		if err != nil {
			return nil, fmt.Errorf("level %d: %w", level, err)
		}

		parentHash, err := hashPair(left.MerkelLeaf, right.MerkelLeaf)
		if err != nil {
			return nil, fmt.Errorf("level %d: %w", level, err)
		}

		if left.ParentMerkelLeaf != "" && !hashEqual(left.ParentMerkelLeaf, parentHash) {
			return nil, fmt.Errorf("level %d: left parent hash mismatch", level)
		}
		if right.ParentMerkelLeaf != "" && !hashEqual(right.ParentMerkelLeaf, parentHash) {
			return nil, fmt.Errorf("level %d: right parent hash mismatch", level)
		}

		next, ok := grouped[level-1]
		if !ok || len(next) == 0 {
			return nil, fmt.Errorf("level %d: missing parent nodes for level %d", level, level-1)
		}
		if !containsHash(next, parentHash) {
			return nil, fmt.Errorf("level %d: computed parent hash not present at level %d", level, level-1)
		}

		rollingHash = parentHash
	}

	return &Result{
		RootHash: normalizeHex(rollingHash),
		LeafHash: normalizeHex(file.Self.MerkelLeaf),
		Levels:   len(levels),
	}, nil
}

func siblingPair(nodes []proof.Node) (proof.Node, proof.Node, error) {
	if len(nodes) == 0 {
		return proof.Node{}, proof.Node{}, errors.New("empty level")
	}

	var left, right *proof.Node
	for i := range nodes {
		node := nodes[i]
		switch node.Role {
		case 1:
			if left == nil {
				left = &node
			}
		case 2:
			if right == nil {
				right = &node
			}
		case 3:
			// role 3 indicates root; ignore during intermediate levels.
		default:
			return proof.Node{}, proof.Node{}, fmt.Errorf("unknown role %d", node.Role)
		}
	}

	if left == nil && right == nil {
		return proof.Node{}, proof.Node{}, errors.New("level lacks left/right roles")
	}
	if left == nil {
		left = right
	}
	if right == nil {
		right = left
	}

	return *left, *right, nil
}

func containsHash(nodes []proof.Node, hash string) bool {
	for _, node := range nodes {
		if hashEqual(node.MerkelLeaf, hash) {
			return true
		}
	}
	return false
}

func hashPair(leftHex, rightHex string) (string, error) {
	leftBytes, err := hex.DecodeString(normalizeHex(leftHex))
	if err != nil {
		return "", fmt.Errorf("decode left hash: %w", err)
	}
	rightBytes, err := hex.DecodeString(normalizeHex(rightHex))
	if err != nil {
		return "", fmt.Errorf("decode right hash: %w", err)
	}

	data := append(leftBytes, rightBytes...)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func normalizeHex(input string) string {
	return strings.ToLower(strings.TrimPrefix(input, "0x"))
}

func hashEqual(a, b string) bool {
	return normalizeHex(a) == normalizeHex(b)
}
