package apiServer

import (
	"JiraConnector/configReader"
	"JiraConnector/connector"
	"JiraConnector/logging"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	configReader *configReader.ConfigReader
	config       *ServerConfig
	logger       *logging.Logger
	connector    *connector.JiraConnector
}

func NewServer() *Server {
	reader := configReader.NewConfigReader()
	return &Server{
		configReader: reader,
		config:       NewServerConfig(reader.GetLocalServerPort()),
		logger:       logging.NewLogger(),
		connector:    connector.NewJiraConnector(),
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

	server.connector.GetProjectIssues(projectName, 0)
}

func (server *Server) projects(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		server.logger.Log(logging.ERROR, "Incorrect http method for /projects")
		writer.WriteHeader(400)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
}

func (server *Server) routes() {
	http.HandleFunc("/updateProject", server.updateProject)
	http.HandleFunc("/projects", server.projects)
}

func (server *Server) start() {
	err := http.ListenAndServe(":"+strconv.Itoa(int(server.config.port)), nil)
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
