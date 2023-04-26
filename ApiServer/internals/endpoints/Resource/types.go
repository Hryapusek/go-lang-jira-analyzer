package endpoints

import "time"

type IssueInfo struct {
	Id          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	AuthorID    int       `json:"author_id"`
	AssigneeId  int       `json:"assignee_id"`
	Key         string    `json:"key"`
	Summary     string    `json:"summary"`
	Type        string    `json:"type"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedTime time.Time `json:"created_time"`
	ClosedTime  time.Time `json:"closed_time"`
	UpdatedTime time.Time `json:"updated_time"`
	TimeSpent   time.Time `json:"timespent"`
}

type ProjectInfo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type HistoryInfo struct {
	IssueID    int       `json:"issue_id"`
	AuthorID   int       `json:"author_id"`
	ChangeTime time.Time `json:"change_time"`
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
}
