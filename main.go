package main

import (
	"fmt"
	"log"
	"order-matching-engine/db"
	"order-matching-engine/router"
)

func main() {
	db.Connect()

	r := router.SetupRouter()
	fmt.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
