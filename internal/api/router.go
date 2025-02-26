package api

import (
	"anpurnama/summarizer-backend/internal/repository"
	"anpurnama/summarizer-backend/internal/service/extractor"
	"anpurnama/summarizer-backend/internal/service/gemini"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	historyRepo repository.HistoryRepository,
	styleRepo repository.StyleRepository,
	extractor extractor.ContentExtractor,
	geminiClient *gemini.Client,
) *gin.Engine {
	router := gin.Default()
	handler := NewHandler(historyRepo, styleRepo, extractor, geminiClient)

	// Enable CORS
	router.Use(CORSMiddleware())

	// Setup error handling middleware
	router.Use(ErrorHandler())

	// API routes group
	api := router.Group("/api")
	{
		api.POST("/summarize", validateSummarizeRequest(), handler.HandleSummarize)
		api.GET("/history", handler.HandleGetHistory)
		api.GET("/history/:id", handler.HandleGetHistoryById)
		api.GET("/search", handler.HandleSearch)
	}

	return router
}
