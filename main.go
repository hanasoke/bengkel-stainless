package main

import (
	"log"
	"os"

	"bengkel-stainless/config"
	"bengkel-stainless/routes"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	log.Println("Database terhubung")

	r := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server jalan di :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
