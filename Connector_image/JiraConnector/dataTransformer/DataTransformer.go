package dataTransformer

import (
	"JiraConnector/jsonmodels"
	"time"
)

type DataTransformer struct {
}

func (dataTransformer *DataTransformer) TransformIssues(issues map[jsonmodels.Issue]struct{}) []jsonmodels.TransformedIssue {
	var transformedIssues []jsonmodels.TransformedIssue

	for key := range issues {
		createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", key.Fields.CreatedTime)
		updatedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", key.Fields.UpdatedTime)
		closedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", key.Fields.ClosedTime)
		if len(key.Fields.ClosedTime) == 0 {
			closedTime = createdTime.Add(time.Hour * 24 * 365 * 10)
		}
		timespent := closedTime.Sub(createdTime)

		transformedIssues = append(transformedIssues, jsonmodels.TransformedIssue{
			Project:     key.Fields.Project.Name,
			Author:      key.Fields.Creator.Name,
			Assignee:    key.Fields.AssigneeName.Name,
			Key:         key.Key,
			Summary:     key.Fields.Summary,
			Description: key.Fields.Description,
			Type:        key.Fields.Type.Name,
			Priority:    key.Fields.Priority.Name,
			Status:      key.Fields.Status.Name,
			CreatedTime: createdTime,
			UpdatedTime: updatedTime,
			ClosedTime:  closedTime,
			Timespent:   int64(timespent.Seconds()),
		})
	}

	return transformedIssues
}

func NewDataTransformer() *DataTransformer {
	return &DataTransformer{}
}
