package connector

import (
	"JiraConnector/configReader"
	"JiraConnector/dbPusher"
	"math"
	"math/rand"
	"sync"
	"time"
)

type JiraConnector struct {
	configReader   *configReader.ConfigReader
	repositoryUrl  string
	databasePusher *dbPusher.DatabasePusher
}

func NewJiraConnector() *JiraConnector {
	reader := configReader.NewConfigReader()
	return &JiraConnector{
		configReader:   reader,
		repositoryUrl:  reader.GetJiraRepositoryUrl(),
		databasePusher: dbPusher.NewDatabasePusher(),
	}
}

func (connector *JiraConnector) GetProjectIssues(projectName string, timeToWaitMs int) {
	waitGroup := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wasError := false
	for i := 0; i < int(connector.configReader.GetThreadCount()); i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			mutex.Lock()
			// projectIssues = append(projectIssues, nil)
			mutex.Unlock()
		}()
	}
	waitGroup.Wait()

	if !wasError {
		connector.databasePusher.PushIssues()
	} else {
		time.Sleep(time.Duration(rand.Intn(timeToWaitMs)) * time.Millisecond)
		connector.GetProjectIssues(projectName, int(float64(timeToWaitMs)*math.Phi))
	}
}
