package dbPusher

import (
	"JiraConnector/configReader"
	"JiraConnector/jsonmodels"
	"JiraConnector/logging"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type DatabasePusher struct {
	configReader *configReader.ConfigReader
	logger       *logging.Logger
	database     *sql.DB
}

func NewDatabasePusher() *DatabasePusher {
	configReaderInstance := configReader.NewConfigReader()
	loggerInstance := logging.NewLogger()

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		configReaderInstance.GetDbHost(),
		configReaderInstance.GetLocalServerPort(),
		configReaderInstance.GetDbUsername(),
		configReaderInstance.GetDbPassword(),
		configReaderInstance.GetDbName())

	database, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		loggerInstance.Log(logging.ERROR, "Can not open connection to database"+err.Error())
		log.Fatal("Can not open connection to database", err.Error())
	}

	fmt.Println("Successfully connected to database")

	return &DatabasePusher{
		configReader: configReaderInstance,
		logger:       loggerInstance,
		database:     database,
	}
}

func (databasePusher *DatabasePusher) PushIssues(issues []jsonmodels.TransformedIssue) {
}
