package main

import (
	"fmt"
	"log"
	"order-matching-engine/db"
	"order-matching-engine/router"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := router.SetupRouter()
	fmt.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
