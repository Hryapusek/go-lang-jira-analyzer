package jsonmodels

type JiraProject struct {
	Name string `json:"name"`
	Link string `json:"self"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
	PageInfo PageInfo  `json:"pageInfo"`
}

type Project struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type PageInfo struct {
	PageCount     int `json:"pageCount"`
	CurrentPage   int `json:"currentPage"`
	ProjectsCount int `json:"projectsCount"`
}
