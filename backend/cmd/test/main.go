package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	// Set environment variables
	viper.AutomaticEnv()

	// Print environment variables
	fmt.Println("DATABASE_URL:", viper.GetString("DATABASE_URL"))
	fmt.Println("DATABASE_AUTH_TOKEN:", viper.GetString("DATABASE_AUTH_TOKEN"))

	// Print environment variables directly
	fmt.Println("DATABASE_URL (os):", os.Getenv("DATABASE_URL"))
	fmt.Println("DATABASE_AUTH_TOKEN (os):", os.Getenv("DATABASE_AUTH_TOKEN"))

	// Exit with success
	os.Exit(0)
}
