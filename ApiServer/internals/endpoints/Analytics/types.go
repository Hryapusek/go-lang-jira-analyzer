package endpoints

type IssueForGraphOne struct {
	Id                int    `json:"id"`
	TimeOpenedSeconds uint64 `json:"time_open_seconds"`
}

type IssueForGraphTwo struct {
	Id              int    `json:"id"`
	TimeOpenSeconds uint64 `json:"time_open_seconds"`
}

type GraphThreeData struct {
	Date         uint64 `json:"timestamp"`
	CreateIssues int    `json:"create_issues"`
	ClosedIssues int    `json:"closed_issues"`
}
