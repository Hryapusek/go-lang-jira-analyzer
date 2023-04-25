package endpoints

import "time"

type IssueForGraphOne struct {
	Id              int       `json:"id"`
	Summary         string    `json:"summary"`
	CreatedTime     time.Time `json:"created_time"`
	ClosedTime      time.Time `json:"closed_time"`
	TimeOpenSeconds float64   `json:"time_open_seconds"`
}
