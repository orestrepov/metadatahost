package model

import "time"

// Host db model
type Host struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name             string
	ServersChanged   bool
	SslGrade         string
	PreviousSslGrade string
	Logo             string
	Title            string
	IsDown           bool
}

// Host response json model
type HostResponse struct {
	Name             string           `json:"name"`
	ServersChanged   bool             `json:"servers_changed"`
	SslGrade         string           `json:"ssl_grade"`
	PreviousSslGrade string           `json:"previous_ssl_grade"`
	Logo             string           `json:"logo"`
	Title            string           `json:"title"`
	IsDown           bool             `json:"is_down"`
	Servers          []ServerResponse `json:"servers"`
	TotalServers     int			  `json:"totals_servers"`
	LastSearch       string			  `json:"last_search"`
}
