package dbPusher

import (
	"JiraConnector/configReader"
	"JiraConnector/jsonmodels"
	"JiraConnector/logging"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DatabasePusher struct {
	configReader *configReader.ConfigReader
	logger       *logging.Logger
	database     *sql.DB
}

func NewDatabasePusher() *DatabasePusher {
	configReaderInstance := configReader.NewConfigReader()
	loggerInstance := logging.NewLogger()

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configReaderInstance.GetDbHost(),
		configReaderInstance.GetDbPort(),
		configReaderInstance.GetDbUsername(),
		configReaderInstance.GetDbPassword(),
		configReaderInstance.GetDbName())

	database, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		loggerInstance.Log(logging.ERROR, "Can not open connection to database"+err.Error())
		log.Fatal("Can not open connection to database", err.Error())
	}

	return &DatabasePusher{
		configReader: configReaderInstance,
		logger:       loggerInstance,
		database:     database,
	}
}

func (databasePusher *DatabasePusher) PushIssues(issues []jsonmodels.TransformedIssue) {
	projectId := databasePusher.extractProjectId(issues[0].Project)
	transaction, err := databasePusher.database.Begin()
	if err != nil {
		databasePusher.logger.Log(logging.ERROR, "Can not open a transaction for project="+issues[0].Project)
		return
	}

	statement, err := transaction.Prepare("INSERT INTO \"issue\" (projectid, authorid, assigneeid, key, summary, description, type, priority, status, createdtime, closedtime, updatedtime, timespent) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id")
	if err != nil {
		databasePusher.logger.Log(logging.ERROR, "Can not create a prepare statement for project="+issues[0].Project)
		return
	}
	defer func(statement *sql.Stmt) {
		_ = statement.Close()
	}(statement)

	for _, issue := range issues {
		authorId := databasePusher.extractAuthorId(issue.Author)
		assigneeId := databasePusher.extractAssigneeId(issue.Assignee)

		issueId := databasePusher.extractIssueId(issue.Key)
		if issueId != 0 {
			databasePusher.deleteIssueById(issueId)
		}

		err = statement.QueryRow(projectId, authorId, assigneeId, issue.Key, issue.Summary,
			issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime,
			issue.UpdatedTime, issue.Timespent).Scan(&issueId)
		if err != nil {
			_ = transaction.Rollback()
			break
		}
	}

	err = transaction.Commit()
	if err != nil {
		databasePusher.logger.Log(logging.ERROR, "Error while committing a transaction for project="+issues[0].Project)
	}

}

func (databasePusher *DatabasePusher) extractProjectId(projectTitle string) int {
	var projectId int
	_ = databasePusher.database.QueryRow("SELECT id FROM \"projects\" WHERE title=$1", projectTitle).Scan(&projectId)
	if projectId == 0 {
		_ = databasePusher.database.QueryRow("INSERT INTO \"projects\" (title) VALUES($1) RETURNING id", projectTitle).Scan(&projectId)
	}

	return projectId
}

func (databasePusher *DatabasePusher) extractAuthorId(authorName string) int {
	var authorId int
	_ = databasePusher.database.QueryRow("SELECT id FROM \"author\" WHERE name=$1", authorName).Scan(&authorId)
	if authorId == 0 {
		_ = databasePusher.database.QueryRow("INSERT INTO \"author\" (name) VALUES($1) RETURNING id", authorName).Scan(&authorId)
	}

	return authorId
}

func (databasePusher *DatabasePusher) extractAssigneeId(assigneeName string) int {
	var assigneeId int
	_ = databasePusher.database.QueryRow("SELECT id FROM \"author\" WHERE name=$1", assigneeName).Scan(&assigneeId)
	if assigneeId == 0 {
		_ = databasePusher.database.QueryRow("INSERT INTO \"author\" (name) VALUES($1) RETURNING id", assigneeName).Scan(&assigneeId)
	}

	return assigneeId
}

func (databasePusher *DatabasePusher) extractIssueId(issueKey string) int {
	var issueId int
	_ = databasePusher.database.QueryRow("SELECT id FROM \"issue\" WHERE key=$1", issueKey).Scan(&issueId)
	return issueId
}

func (databasePusher *DatabasePusher) deleteIssueById(issueId int) {
	_, _ = databasePusher.database.Exec("DELETE FROM \"issue\" WHERE id=$1", issueId)
}
