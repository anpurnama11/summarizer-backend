package api

import (
	"anpurnama/summarizer-backend/internal/repository"
	"anpurnama/summarizer-backend/internal/service"
	"anpurnama/summarizer-backend/internal/service/extractor"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	historyRepo repository.HistoryRepository,
	styleRepo repository.StyleRepository,
	extractor extractor.ContentExtractor,
	summarizer service.Summarizer,
) *gin.Engine {
	router := gin.Default()
	handler := NewHandler(historyRepo, styleRepo, extractor, summarizer)

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
