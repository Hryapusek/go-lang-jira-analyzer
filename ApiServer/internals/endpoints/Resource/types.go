package endpoints

type IssueInfo struct {
	IssueID     int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	AuthorID    int    `json:"author_id"`
	AssigneeId  int    `json:"assigned_id"`
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
	ProjectID int    `json:"id"`
	Title     string `json:"title"`
}

type HistoryInfo struct {
	IssueID    int    `json:"issue_id"`
	AuthorID   int    `json:"author_id"`
	ChangeTime uint64 `json:"change_time"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
}

type IssueResponse struct {
	IssueInfo
	ProjectID ProjectInfo `json:"project"`
}

type ProjectResponse struct {
	ProjectInfo
}

type HistoryResponse struct {
	Histories []HistoryInfo `json:"histories"`
	IssueID   IssueInfo     `json:"issue"`
}

type Link struct {
	URL string `json:"href"`
}

type ReferencesLinks struct {
	LinkSelf      Link `json:"self"`
	LinkIssues    Link `json:"issues"`
	LinkProjects  Link `json:"projects"`
	LinkHistories Link `json:"histories"`
}

type RestAPIGetResponseSchema struct {
	Links ReferencesLinks `json:"_links"`
	Info  interface{}     `json:"data"`
}

type RestAPIPostResponseSchema struct {
	Id         int `json:"id"`
	StatusCode int `json:"status_code"`
}
