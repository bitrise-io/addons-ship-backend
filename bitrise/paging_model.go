package bitrise

type pagingResponseModel struct {
	TotalItemCount int64  `json:"total_item_count"`
	PageItemLimit  uint   `json:"page_item_limit"`
	Next           string `json:"next,omitempty"`
}
