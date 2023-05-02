package apiServer

import (
	"JiraConnector/configReader"
	"JiraConnector/connector"
	"JiraConnector/dataTransformer"
	"JiraConnector/dbPusher"
	"JiraConnector/logging"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"fmt"
)

type Server struct {
	configReader    *configReader.ConfigReader
	config          *ServerConfig
	logger          *logging.Logger
	connector       *connector.JiraConnector
	dataTransformer *dataTransformer.DataTransformer
	databasePusher  *dbPusher.DatabasePusher
}

func NewServer() *Server {
	reader := configReader.NewConfigReader()
	connectorInstance := connector.NewJiraConnector()

	return &Server{
		configReader:    reader,
		config:          NewServerConfig(reader.GetLocalServerHost(), reader.GetLocalServerPort()),
		logger:          logging.NewLogger(),
		connector:       connectorInstance,
		dataTransformer: dataTransformer.NewDataTransformer(),
		databasePusher:  connectorInstance.DatabasePusher,
	}
}

func (server *Server) updateProject(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		server.logger.Log(logging.ERROR, "Incorrect http method for /updateProject")
		writer.WriteHeader(400)
		return
	}

	projectName := request.URL.Query().Get("project")
	if len(projectName) == 0 {
		server.logger.Log(logging.ERROR, "Project name was not passed to /updateProject")
		writer.WriteHeader(400)
		return
	}

	issues, err := server.connector.GetProjectIssues(projectName, server.configReader.GetMinTimeSleep())
	if err != nil {
		server.logger.Log(logging.ERROR, "Error while downloading issues for project \""+projectName+"\"")
		writer.WriteHeader(400)
		return
	}

	transformedIssues := server.dataTransformer.TransformIssues(issues)
	server.databasePusher.PushIssues(transformedIssues)
}

func (server *Server) projects(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		server.logger.Log(logging.ERROR, "Incorrect http method for /projects")
		writer.WriteHeader(400)
		return
	}

	limit, page, search := extractProjectParameters(request)

	writer.Header().Set("Content-Type", "application/json")
	projects, err := server.connector.GetProjects(limit, page, search)
	if err != nil {
		server.logger.Log(logging.ERROR, "Error while downloading list of projects")
		writer.WriteHeader(400)
		return
	}
	response, _ := json.Marshal(projects)
	_, _ = writer.Write(response)
}

func extractProjectParameters(request *http.Request) (int, int, string) {
	limit := 20
	page := 1
	search := ""

	limitParam := request.URL.Query().Get("limit")
	if len(limitParam) != 0 {
		limit, _ = strconv.Atoi(limitParam)
	}

	pageParam := request.URL.Query().Get("page")
	if len(pageParam) != 0 {
		page, _ = strconv.Atoi(pageParam)
	}

	searchParam := request.URL.Query().Get("search")
	if len(searchParam) != 0 {
		search = searchParam
	}

	return limit, page, search
}

func (server *Server) routes() {
	http.HandleFunc("/api/v1/connector/updateProject", server.updateProject)
	http.HandleFunc("/api/v1/connector/projects", server.projects)
}

func (server *Server) start() {
	err := http.ListenAndServe(server.config.host+":"+strconv.Itoa(int(server.config.port)), nil)
	fmt.Println(server.config.host+":"+strconv.Itoa(int(server.config.port)))
	if err != nil {
		server.logger.Log(logging.ERROR, "Error while starting a server...")
		log.Fatal()
	}
}

func (server *Server) Start() {
	server.logger.Log(logging.INFO, "Starting server...")
	server.routes()
	server.start()
}
