package data

// StoreFileCounts holds file count metadata for a vector store.
// Mirrors upstream Ai\Responses\Data\StoreFileCounts.
type StoreFileCounts struct {
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Cancelled  int `json:"cancelled"`
	Total      int `json:"total"`
}
