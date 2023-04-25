package endpoints

type IssueForGraphOne struct {
	Id                int     `json:"id"`
	TimeOpenedSeconds float64 `json:"time_open_seconds"`
}

type IssueForGraphTwo struct {
	Id              int     `json:"id"`
	TimeOpenSeconds float64 `json:"time_open_seconds"`
}
