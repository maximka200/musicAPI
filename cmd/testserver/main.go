package main

import (
	"log"
	"musicAPI/tests/testserver"
)

func main() {
	server := testserver.NewTestServer()
	if err := server.Run("localhost:1818"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
