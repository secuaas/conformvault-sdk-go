// Example: quickstart usage of the ConformVault Go SDK
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cv "github.com/secuaas/conformvault-sdk-go"
)

func main() {
	apiKey := os.Getenv("CONFORMVAULT_API_KEY")
	if apiKey == "" {
		log.Fatal("CONFORMVAULT_API_KEY environment variable is required")
	}

	// Create client (defaults to production URL)
	client := cv.NewClient(apiKey)
	ctx := context.Background()

	// List files
	files, err := client.Files.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	fmt.Printf("Found %d files\n", len(files))
	for _, f := range files {
		fmt.Printf("  - %s (%s, %d bytes)\n", f.OriginalName, f.ID, f.Size)
	}

	// List folders
	folders, err := client.Folders.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list folders: %v", err)
	}
	fmt.Printf("Found %d folders\n", len(folders))

	// Create a signature envelope (requires a file ID)
	if len(files) > 0 {
		env, err := client.Signatures.Create(ctx, cv.CreateSignatureRequest{
			FileID:  files[0].ID,
			Subject: "Please sign this document",
			Signers: []cv.CreateSignatureSigner{
				{
					Email: "signer@example.com",
					Name:  "Jane Doe",
					Role:  "signer",
				},
			},
			ExpiryDays: 30,
		})
		if err != nil {
			if cv.IsRateLimited(err) {
				fmt.Println("Rate limited! Try again later.")
			} else {
				log.Fatalf("Failed to create signature: %v", err)
			}
		} else {
			fmt.Printf("Signature envelope created: %s (status: %s)\n", env.ID, env.Status)
		}
	}
}
