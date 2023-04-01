package jsonmodels

type IssuesList struct {
	IssuesCount int     `json:"total"`
	Issues      []Issue `json:"issues"`
}

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary string `json:"summary"`
	Type    struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"issuetype"`
	Status struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"status"`
	Priority struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"priority"`
	Creator struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"creator"`
	Project struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"project"`
}
