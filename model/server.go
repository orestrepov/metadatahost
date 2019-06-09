package model

import "time"

// Server db model
type Server struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Address  string
	SslGrade string
	Country  string
	Owner    string

	Host   Host
	HostID uint
}

// Server response json model
type ServerResponse struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}
