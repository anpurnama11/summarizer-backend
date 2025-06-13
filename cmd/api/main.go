package main

import (
	"anpurnama/summarizer-backend/internal/api"
	"anpurnama/summarizer-backend/internal/database"
	"anpurnama/summarizer-backend/internal/repository"
	"anpurnama/summarizer-backend/internal/service/extractor"
	"anpurnama/summarizer-backend/internal/service/openrouter"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database connection
	db, err := database.NewDB("./db/database.sqlite")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	historyRepo := repository.NewHistoryRepository(db)
	styleRepo := repository.NewStyleRepository(db)

	// Initialize services
	extractor, err := extractor.NewContentExtractor()
	if err != nil {
		log.Fatalf("Failed to create content extractor: %v", err)
	}

	openrouterClient, err := openrouter.NewClient(styleRepo)
	if err != nil {
		log.Fatalf("Failed to create OpenRouter client: %v", err)
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// Setup router
	router := api.SetupRouter(historyRepo, styleRepo, extractor, openrouterClient)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
