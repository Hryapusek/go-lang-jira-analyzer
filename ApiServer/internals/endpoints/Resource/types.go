package endpoints

type IssueInfo struct {
	Id          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	AuthorID    int    `json:"author_id"`
	AssigneeId  int    `json:"assignee_id"`
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	CreatedTime uint64 `json:"created_time"`
	ClosedTime  uint64 `json:"closed_time"`
	UpdatedTime uint64 `json:"updated_time"`
	TimeSpent   uint64 `json:"timespent"`
}

type ProjectInfo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type HistoryInfo struct {
	IssueID    int    `json:"issue_id"`
	AuthorID   int    `json:"author_id"`
	ChangeTime uint64 `json:"change_time"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
}
