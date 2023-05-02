package jsonmodels

import "time"

type TransformedIssue struct {
	Project     string
	Author      string
	Assignee    string
	Key         string
	Summary     string
	Description string
	Type        string
	Priority    string
	Status      string
	CreatedTime time.Time
	ClosedTime  time.Time
	UpdatedTime time.Time
	Timespent   int64
}
