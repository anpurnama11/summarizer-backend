package api

import (
	"anpurnama/summarizer-backend/internal/repository"
	"anpurnama/summarizer-backend/internal/service"
	"anpurnama/summarizer-backend/internal/service/extractor"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	historyRepo repository.HistoryRepository
	styleRepo   repository.StyleRepository
	extractor   extractor.ContentExtractor
	summarizer  service.Summarizer
}

func NewHandler(
	historyRepo repository.HistoryRepository,
	styleRepo repository.StyleRepository,
	extractor extractor.ContentExtractor,
	summarizer service.Summarizer,
) *Handler {
	return &Handler{
		historyRepo: historyRepo,
		styleRepo:   styleRepo,
		extractor:   extractor,
		summarizer:  summarizer,
	}
}

func (h *Handler) HandleSummarize(c *gin.Context) {
	reqInterface, exists := c.Get("summarizeRequest")
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	req, ok := reqInterface.(SummarizeRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request type"})
		return
	}

	// Extract content from URL
	extracted, err := h.extractor.Extract(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to extract content: " + err.Error()})
		return
	}

	// Generate summary using Gemini
	styleName := req.Style
	if styleName == "" {
		styleName = "concise" // Change default to an existing style
	}

	// Get style ID before generating summary
	style, err := h.styleRepo.GetByName(c.Request.Context(), styleName)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid style: " + err.Error()})
		return
	}

	summary, err := h.summarizer.Summarize(c.Request.Context(), extracted.Content, styleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate summary: " + err.Error()})
		return
	}

	// Save to history
	history := &repository.History{
		URL:     req.URL,
		Title:   &extracted.Title,
		Content: extracted.Content,
		Summary: summary,
		StyleID: &style.ID,
	}

	// Ensure language code meets ISO 639-1 format
	if extracted.Language != "" {
		if len(extracted.Language) == 2 {
			lowerLang := strings.ToLower(extracted.Language)
			history.Language = &lowerLang
		}
	}

	if err := h.historyRepo.Create(c.Request.Context(), history); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, SummarizeResponse{
		Summary: summary,
		Title:   *history.Title,
		URL:     req.URL,
	})
}

func (h *Handler) HandleGetHistory(c *gin.Context) {
	limit := 10 // Default limit
	offset := 0 // Default offset

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	repoHistories, err := h.historyRepo.ListWithStyles(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch history: " + err.Error()})
		return
	}

	totalSize, err := h.historyRepo.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch total size: " + err.Error()})
		return
	}

	histories := make([]History, len(repoHistories))
	for i, h := range repoHistories {
		histories[i] = toAPIHistory(h)
	}

	c.JSON(http.StatusOK, HistoryResponse{
		Histories: histories,
		TotalSize: totalSize,
	})
}

func (h *Handler) HandleGetHistoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
		return
	}

	history, err := h.historyRepo.GetWithStyle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch history: " + err.Error()})
		return
	}

	if history == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "History not found"})
		return
	}

	c.JSON(http.StatusOK, toAPIHistory(*history))
}

func (h *Handler) HandleSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Search query is required"})
		return
	}

	limit := 10 // Default limit
	offset := 0 // Default offset

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	results, err := h.historyRepo.Search(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to search: " + err.Error()})
		return
	}

	histories := make([]History, len(results))
	for i, h := range results {
		histories[i] = toAPIHistory(h)
	}

	c.JSON(http.StatusOK, histories)
}

func toAPIHistory(h repository.History) History {
	title := ""
	if h.Title != nil {
		title = *h.Title
	}

	return History{
		ID:        strconv.Itoa(h.ID),
		URL:       h.URL,
		Summary:   h.Summary,
		Title:     title,
		CreatedAt: h.CreatedAt.Format(time.RFC3339),
	}
}
