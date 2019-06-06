package model

// Domain search history response json model
type HistoryResponse struct {
	Items []HostResponse `json:"items"`
}
