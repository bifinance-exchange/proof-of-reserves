package main

import (
	"flag"
	"fmt"
	"os"

	"proof-of-reserves/internal/merkle"
	"proof-of-reserves/internal/proof"
)

func main() {
	filePath := flag.String("file", "test.json", "path to a proof JSON file")
	flag.Parse()

	proofFile, err := proof.Load(*filePath)
	if err != nil {
		exitWithError(err)
	}

	result, err := merkle.Verify(proofFile)
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("Verification successful!\n")
	fmt.Printf("  Audit ID : %s\n", proofFile.Self.AuditID)
	fmt.Printf("  Leaf Hash: %s\n", result.LeafHash)
	fmt.Printf("  Root Hash: %s\n", result.RootHash)
	fmt.Printf("  Levels   : %d\n", result.Levels)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
	os.Exit(1)
}
