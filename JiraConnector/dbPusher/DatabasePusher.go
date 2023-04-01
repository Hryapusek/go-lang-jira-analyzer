package dbPusher

import (
	"JiraConnector/dataTransformer"
)

type DatabasePusher struct {
	transformer *dataTransformer.DatabaseTransformer
}

func NewDatabasePusher() *DatabasePusher {
	return &DatabasePusher{
		transformer: dataTransformer.NewDatabaseTransformer(),
	}
}

func (databasePusher *DatabasePusher) PushIssues() {

}
