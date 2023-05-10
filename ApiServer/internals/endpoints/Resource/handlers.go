package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

func HandlerGetIssue(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid Issue ID in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	issue, err := GetIssueInfoByID(id)
	if err != nil {
		log.Printf("Request ended up with mistake of database: %s", err.Error())
		rw.WriteHeader(400)
		return
	}

	project, err := GetProjectInfoByID(issue.ProjectID)
	if err != nil {
		log.Printf("Request ended up with mistake of database: %s", err.Error())
		rw.WriteHeader(400)
		return
	}

	var issueResponse = RestAPIGetResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/issues/%d", id)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Info: IssueResponse{
			IssueInfo: issue,
			ProjectID: project,
		},
	}

	data, err := json.MarshalIndent(issueResponse, "", "\t")
	if err != nil {
		log.Printf("Error with extracting info about issue project with id=%d", id)
		rw.WriteHeader(400)
		return
	}

	rw.WriteHeader(200)
	_, err = rw.Write(data)
	if err != nil {
		return
	}
}

func HandlerGetHistory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid Issue ID in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	history, err := GetAllHistoryInfoByIssueID(id)
	if err != nil {
		log.Printf("Request for histories ended up with mistake of database: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	issue := IssueInfo{}
	issue, err = GetIssueInfoByID(id)
	if err != nil {
		log.Printf("Request ended up with mistake of database: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var historyResponse = RestAPIGetResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/histories/%d", id)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Info: HistoryResponse{
			Histories: history,
			IssueID:   issue,
		},
	}

	data, err := json.MarshalIndent(historyResponse, "", "\t")
	if err != nil {
		log.Printf("Error with extracting info about history with id=%d", id)
		rw.WriteHeader(400)
		return
	}

	rw.WriteHeader(200)
	_, err = rw.Write(data)
	if err != nil {
		return
	}
}

func HandlerGetProject(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid Issue ID in path \"%s\"", r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := GetProjectInfoByID(id)

	if err != nil {
		log.Printf("Request ended up with mistake of database: %s", err.Error())
		rw.WriteHeader(400)
		return
	}

	var projectResponse = RestAPIGetResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/projects/%d", id)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Info: ProjectResponse{
			ProjectInfo: project,
		},
	}

	data, err := json.MarshalIndent(projectResponse, "", "\t")
	if err != nil {
		log.Printf("Error with extracting info about project with id=%d", id)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)
	if err != nil {
		return
	}

}

func HandlerPostIssue(rw http.ResponseWriter, r *http.Request) {
	var data IssueInfo
	body, err := io.ReadAll(r.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &data); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var statusCode int
	id, err := PutIssueToDB(data)
	if err != nil {
		log.Printf("Error %s occured while puting issue to DB", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		id = -1
		statusCode = http.StatusInternalServerError
	} else {
		rw.WriteHeader(http.StatusCreated)
		statusCode = http.StatusCreated
	}

	resp, err := json.MarshalIndent(RestAPIPostResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/issues/%d", id)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Id:         id,
		StatusCode: statusCode,
	},
		"", "\t")

	if err != nil {
		log.Println(err.Error())
	}
	_, err = rw.Write(resp)
	if err != nil {
		return
	}
}

func HandlerPostHistory(rw http.ResponseWriter, r *http.Request) {
	var data HistoryInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &data); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var statusCode int
	var id int
	err = PutHistoryToDB(data)
	if err != nil {
		log.Printf("Error %s occured while puting history to DB", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		id = -1
		statusCode = http.StatusInternalServerError
	} else {
		rw.WriteHeader(http.StatusCreated)
		statusCode = http.StatusCreated
		id = data.IssueID
	}

	resp, err := json.MarshalIndent(RestAPIPostResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/histories/%d", data.IssueID)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Id:         id,
		StatusCode: statusCode,
	},
		"", "\t")

	if err != nil {
		log.Println(err.Error())
	}
	_, err = rw.Write(resp)
	if err != nil {
		return
	}
}

func HandlerPostProject(rw http.ResponseWriter, r *http.Request) {
	var data ProjectInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &data); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var statusCode int
	id, err := PutProjectToDB(data)
	if err != nil {
		log.Printf("Error %s occured while puting issue to DB", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		id = -1
		statusCode = http.StatusInternalServerError
	} else {
		rw.WriteHeader(http.StatusCreated)
		statusCode = http.StatusCreated
	}

	resp, err := json.MarshalIndent(RestAPIPostResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/projects/%d", id)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Id:         id,
		StatusCode: statusCode,
	},
		"", "\t")

	if err != nil {
		log.Println(err.Error())
	}
	_, err = rw.Write(resp)
	if err != nil {
		return
	}
}

func HandlerGetProjectByTitle(rw http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := GetProjectInfoByTitle(title)

	if err != nil {
		log.Printf("Request ended up with mistake of database: %s", err.Error())
		rw.WriteHeader(400)
		return
	}

	var projectResponse = RestAPIGetResponseSchema{
		Links: ReferencesLinks{
			LinkSelf:      Link{fmt.Sprintf("/api/v1/projects/%d", project.ProjectID)},
			LinkIssues:    Link{"/api/v1/issues"},
			LinkProjects:  Link{"/api/v1/projects"},
			LinkHistories: Link{"/api/v1/histories"},
		},
		Info: ProjectResponse{
			ProjectInfo: project,
		},
	}

	data, err := json.MarshalIndent(projectResponse, "", "\t")
	if err != nil {
		log.Printf("Error with extracting info about project with id=%d", project.ProjectID)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)
	if err != nil {
		return
	}
}

func HandlerGetIssuesByProjectId(rw http.ResponseWriter, r *http.Request) {
	projectId, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		log.Printf("Invalid id=%s for issues request with project id.", r.URL.Query().Get("projectId"))
		rw.WriteHeader(http.StatusBadRequest)
	}

	var offset int
	offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 20
	}

	var limit int
	limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	issues, err := GetIssuesWithProjectId(projectId, offset, limit)
	if err != nil {
		log.Printf("Error with database: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	project, err := GetProjectInfoByID(projectId)
	if err != nil {
		log.Printf("Error on getting project info.")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var completeResponse []RestAPIGetResponseSchema
	for _, issue := range issues {
		var issueResponse = RestAPIGetResponseSchema{
			Links: ReferencesLinks{
				LinkSelf:      Link{fmt.Sprintf("/api/v1/issues/%d", issue.IssueID)},
				LinkIssues:    Link{"/api/v1/issues"},
				LinkProjects:  Link{"/api/v1/projects"},
				LinkHistories: Link{"/api/v1/histories"},
			},
			Info: IssueResponse{
				IssueInfo: issue,
				ProjectID: project,
			},
		}

		completeResponse = append(completeResponse, issueResponse)
	}

	data, err := json.MarshalIndent(completeResponse, "", "\t")
	if err != nil {
		log.Printf("Error with extracting info about issues with projectId=%d", projectId)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)
	if err != nil {
		return
	}

	log.Printf("Getting issues request by project id=%d", projectId)
}

func HandlerGetAllProject(rw http.ResponseWriter, r *http.Request) {
	var offset int
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 20
	}

	var limit int
	limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	projects, err := GetAllProjects(offset, limit)
	if err != nil {
		log.Printf("Error with database: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var completeResponse []RestAPIGetResponseSchema
	for _, project := range projects {
		var projectResponse = RestAPIGetResponseSchema{
			Links: ReferencesLinks{
				LinkSelf:      Link{fmt.Sprintf("/api/v1/projects/%d")},
				LinkIssues:    Link{"/api/v1/issues"},
				LinkProjects:  Link{"/api/v1/projects"},
				LinkHistories: Link{"/api/v1/histories"},
			},
			Info: project,
		}

		completeResponse = append(completeResponse, projectResponse)
	}

	data, err := json.MarshalIndent(completeResponse, "", "\t")
	if err != nil {
		log.Printf("Error with extracting info about all projects.")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)
	if err != nil {
		return
	}

	log.Printf("Getting all projects request.")
	rw.WriteHeader(http.StatusNotImplemented)
}
