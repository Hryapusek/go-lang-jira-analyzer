package endpoints

type IssueInfo struct {
	IssueID     int    `json:"id,omitempty" require:"true"`
	ProjectID   int    `json:"project_id,omitempty" require:"true"`
	AuthorID    int    `json:"author_id,omitempty" require:"true"`
	AssigneeId  int    `json:"assigned_id,omitempty" require:"true"`
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
	ProjectID int    `json:"id,omitempty" require:"true"`
	Title     string `json:"title"`
}

type HistoryInfo struct {
	IssueID    int    `json:"issue_id,omitempty" require:"true"`
	AuthorID   int    `json:"author_id,omitempty" require:"true"`
	ChangeTime uint64 `json:"change_time"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
}

type IssueResponse struct {
	IssueInfo
	ProjectID ProjectInfo `json:"project,omitempty" require:"true"`
}

type ProjectResponse struct {
	ProjectInfo
}

type HistoryResponse struct {
	Histories []HistoryInfo `json:"histories,omitempty" require:"true"`
	IssueID   IssueInfo     `json:"issue,omitempty" require:"true"`
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
	Info  interface{}     `json:"data,omitempty" require:"true"`
}

type RestAPIPostResponseSchema struct {
	Links      ReferencesLinks `json:"_links"`
	Id         int             `json:"id"`
	StatusCode int             `json:"status_code"`
}
