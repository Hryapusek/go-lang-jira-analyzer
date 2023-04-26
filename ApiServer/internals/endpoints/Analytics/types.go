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

type GraphFourData struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type GraphFiveAndSixData struct {
	PriorityType string `json:"priority_type"`
	Count        int    `json:"count"`
}
