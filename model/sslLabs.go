package model

// Host struct from SSLLabs
type SSLHost struct {
	Name      string        `json:"host"`
	Status    string        `json:"status"`
	Endpoints []SSLEndpoint `json:"endpoints"`
}

// Host Endpoint from SSLLabs
type SSLEndpoint struct {
	IPAddress     string `json:"ipAddress"`
	Grade         string `json:"grade"`
	StatusMessage string `json:"statusMessage"`
}
