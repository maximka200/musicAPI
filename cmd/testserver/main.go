package main

import (
	"log"
	"musicAPI/testserver"
)

func main() {
	server := testserver.NewTestServer()
	if err := server.Run("localhost:1808"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
