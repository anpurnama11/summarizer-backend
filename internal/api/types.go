package api

type SummarizeRequest struct {
	URL   string `json:"url" binding:"required"`
	Style string `json:"style,omitempty"`
}

type SummarizeResponse struct {
	Summary string `json:"summary"`
	Title   string `json:"title"`
	URL     string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HistoryResponse struct {
	Histories []History `json:"histories"`
	TotalSize int       `json:"total_size"`
}

type History struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Summary   string `json:"summary"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}
