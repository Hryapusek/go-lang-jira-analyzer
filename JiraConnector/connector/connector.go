package connector

import (
	"JiraConnector/configReader"
	"JiraConnector/dbPusher"
	"JiraConnector/jsonmodels"
	"JiraConnector/logging"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type JiraConnector struct {
	configReader   *configReader.ConfigReader
	repositoryUrl  string
	databasePusher *dbPusher.DatabasePusher
	logger         *logging.Logger
}

func NewJiraConnector() *JiraConnector {
	reader := configReader.NewConfigReader()
	return &JiraConnector{
		configReader:   reader,
		repositoryUrl:  reader.GetJiraRepositoryUrl(),
		databasePusher: dbPusher.NewDatabasePusher(),
		logger:         logging.NewLogger(),
	}
}

func (connector *JiraConnector) GetProjectIssues(projectName string, timeToWaitMs int) (map[jsonmodels.Issue]struct{}, error) {
	httpClient := &http.Client{}
	response, err := httpClient.Get(connector.configReader.GetJiraRepositoryUrl() +
		"/rest/api/2/search?jql=project=" + projectName + "&expand=changelog&startAt=0&maxResults=1")

	if err != nil || response.StatusCode != http.StatusOK {
		connector.logger.Log(logging.ERROR, "Unable to get issues for project "+projectName)
		return map[jsonmodels.Issue]struct{}{}, nil
	}

	body, _ := io.ReadAll(response.Body)
	var issueResponse jsonmodels.IssuesList
	_ = json.Unmarshal(body, &issueResponse)

	totalIssuesCount := issueResponse.IssuesCount

	if totalIssuesCount == 0 {
		return map[jsonmodels.Issue]struct{}{}, nil
	}

	issues := map[jsonmodels.Issue]struct{}{}
	issues[issueResponse.Issues[0]] = struct{}{}

	waitGroup := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wasError := false

	threadCount := connector.configReader.GetThreadCount()
	issuesPerRequest := connector.configReader.GetIssuesPerRequest()

	stop := make(chan struct{})

	for i := 0; i < threadCount; i++ {
		waitGroup.Add(1)
		go func(threadNumber int) {
			defer waitGroup.Done()
			select {
			case <-stop:
				connector.logger.Log(logging.ERROR, "Error while reading issues in thread... Stopping all other threads...")
				return
			default:
				threadStartIndex := (totalIssuesCount/threadCount)*threadNumber + 1
				requestCount := int(math.Ceil(float64(totalIssuesCount) / float64(threadCount*issuesPerRequest)))
				for j := 0; j < requestCount; j++ {
					startAt := threadStartIndex + j*issuesPerRequest
					if startAt < totalIssuesCount {
						requestString := connector.configReader.GetJiraRepositoryUrl() + "/rest/api/2/search?jql=project=" +
							projectName + "&expand=changelog&startAt=" + strconv.Itoa(startAt) +
							"&maxResults=" + strconv.Itoa(issuesPerRequest)

						response, requestErr := httpClient.Get(requestString)
						body, responseReadErr := io.ReadAll(response.Body)

						if requestErr != nil || responseReadErr != nil {
							wasError = true
							close(stop)
							return
						}

						var issueResponse jsonmodels.IssuesList
						_ = json.Unmarshal(body, &issueResponse)

						mutex.Lock()
						for _, elem := range issueResponse.Issues {
							issues[elem] = struct{}{}
						}
						mutex.Unlock()
					}
				}
			}
		}(i)
	}
	waitGroup.Wait()

	if wasError {
		time.Sleep(time.Duration(timeToWaitMs) * time.Millisecond)
		newTimeToSleep := int(math.Ceil(float64(timeToWaitMs) * math.Phi))
		connector.logger.Log(logging.INFO, "Error while downloading issues for project \""+
			projectName+"\", waiting now"+strconv.Itoa(timeToWaitMs)+"ms")

		if newTimeToSleep > connector.configReader.GetMaxTimeSleep() {
			return map[jsonmodels.Issue]struct{}{}, errors.New("A lot of time to sleep")
		}

		return connector.GetProjectIssues(projectName, newTimeToSleep)
	}

	return issues, nil
}

func (connector *JiraConnector) GetProjects(limit int, page int, search string) (jsonmodels.ProjectsResponse, error) {
	response, err := http.Get(connector.configReader.GetJiraRepositoryUrl() + "/rest/api/2/project")
	if err != nil {
		connector.logger.Log(logging.ERROR, "Unable to get projects list ")
		return jsonmodels.ProjectsResponse{}, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return jsonmodels.ProjectsResponse{}, err
	}

	var jiraProjects []jsonmodels.JiraProject
	err = json.Unmarshal(body, &jiraProjects)

	if err != nil {
		return jsonmodels.ProjectsResponse{}, err
	}

	var projects []jsonmodels.Project

	projectsCount := 0

	for _, elem := range jiraProjects {
		if isProjectNameSatisfy(elem.Name, search) {
			projectsCount++
			projects = append(projects, jsonmodels.Project{
				Name: elem.Name,
				Link: elem.Link,
			})
		}
	}

	startIndex := limit * (page - 1)
	endIndex := startIndex + limit
	if endIndex >= len(projects) {
		endIndex = len(projects)
	}

	return jsonmodels.ProjectsResponse{
		Projects: projects[startIndex:endIndex],
		PageInfo: jsonmodels.PageInfo{
			PageCount:     int(math.Ceil(float64(projectsCount) / float64(limit))),
			CurrentPage:   page,
			ProjectsCount: projectsCount,
		},
	}, nil
}

func isProjectNameSatisfy(projectName string, search string) bool {
	return strings.Contains(strings.ToLower(projectName), strings.ToLower(search))
}
