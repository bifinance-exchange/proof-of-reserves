package merkle

import (
	"path/filepath"
	"testing"

	"proof-of-reserves/internal/proof"
)

func TestVerifySampleProof(t *testing.T) {
	path := filepath.Join("..", "..", "test.json")
	data, err := proof.Load(path)
	if err != nil {
		t.Fatalf("load proof: %v", err)
	}

	result, err := Verify(data)
	if err != nil {
		t.Fatalf("verify proof: %v", err)
	}

	const expectedRoot = "e5f2729cdbc8c1e2989a4dfcce63ee14ef6c8891348cb14f218ae2659432c0ad"
	if result.RootHash != expectedRoot {
		t.Fatalf("unexpected root hash: got %s want %s", result.RootHash, expectedRoot)
	}
}

func TestVerifyDetectsRootMismatch(t *testing.T) {
	left := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	right := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	parent, err := hashPair(left, right)
	if err != nil {
		t.Fatalf("build parent: %v", err)
	}

	proofFile := &proof.File{
		Self: proof.Node{
			AuditID:          "test",
			Level:            1,
			MerkelLeaf:       left,
			ParentMerkelLeaf: parent,
			Role:             1,
		},
		Path: []proof.Node{
			{
				AuditID:          "test",
				Level:            1,
				MerkelLeaf:       left,
				ParentMerkelLeaf: parent,
				Role:             1,
			},
			{
				AuditID:          "test",
				Level:            1,
				MerkelLeaf:       right,
				ParentMerkelLeaf: parent,
				Role:             2,
			},
			{
				AuditID:    "test",
				Level:      0,
				MerkelLeaf: "deadbeef",
				Role:       3,
			},
		},
	}

	if _, err := Verify(proofFile); err == nil {
		t.Fatalf("expected verification to fail")
	}
}
