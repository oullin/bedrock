package data

// StoreFileCounts holds file count metadata for a vector store.
type StoreFileCounts struct {
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Cancelled  int `json:"cancelled"`
	Total      int `json:"total"`
}
