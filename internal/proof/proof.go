package proof

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// File models the structure of the downloaded JSON payload that contains the
// self node plus every sibling required to rebuild the Merkle branch.
type File struct {
	Path []Node `json:"path"`
	Self Node   `json:"self"`
}

// Node represents a single position in the proof path.
type Node struct {
	AuditID          string            `json:"auditId"`
	Balances         map[string]string `json:"balances"`
	Level            int               `json:"level"`
	MerkelLeaf       string            `json:"merkelLeaf"`
	Nonce            string            `json:"nonce"`
	Role             int               `json:"role"`
	ParentMerkelLeaf string            `json:"parentMerkelLeaf"`
	EncryptUID       string            `json:"encryptUid"`
}

// Load reads the provided file and unmarshals it into a File structure.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read proof file: %w", err)
	}

	var proof File
	if err := json.Unmarshal(data, &proof); err != nil {
		return nil, fmt.Errorf("parse proof json: %w", err)
	}

	if len(proof.Path) == 0 {
		return nil, fmt.Errorf("proof contains no path entries")
	}

	return &proof, nil
}

// Nodes returns all nodes including the self entry, de-duplicated by hash.
func (f *File) Nodes() []Node {
	nodes := make([]Node, 0, len(f.Path)+1)
	nodes = append(nodes, f.Path...)

	hasSelf := false
	for _, node := range nodes {
		if equalHash(node.MerkelLeaf, f.Self.MerkelLeaf) {
			hasSelf = true
			break
		}
	}
	if !hasSelf && f.Self.MerkelLeaf != "" {
		nodes = append(nodes, f.Self)
	}
	return nodes
}

func equalHash(a, b string) bool {
	return strings.EqualFold(strings.TrimPrefix(a, "0x"), strings.TrimPrefix(b, "0x"))
}
