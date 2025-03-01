package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {
	// Generate key pair
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate keys: %v", err)
	}

	// Convert to hex
	publicKeyHex := hex.EncodeToString(publicKey)
	privateKeyHex := hex.EncodeToString(privateKey)

	// Print keys
	fmt.Printf("Public Key (hex): %s\n", publicKeyHex)
	fmt.Printf("Private Key (hex): %s\n", privateKeyHex)

	// Write to .env file
	envContent := fmt.Sprintf(`# Generated PASETO keys
AUTH_PUBLIC_KEY=%s
AUTH_PRIVATE_KEY=%s
`, publicKeyHex, privateKeyHex)

	if err := os.WriteFile(".env.keys", []byte(envContent), 0600); err != nil {
		log.Fatalf("Failed to write keys to file: %v", err)
	}

	fmt.Println("\nKeys have been written to .env.keys")
	fmt.Println("Please copy these values to your .env file")
}
