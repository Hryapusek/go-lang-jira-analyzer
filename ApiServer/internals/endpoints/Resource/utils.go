package endpoints

import "log"

func GetIssueInfoByID(id int) IssueInfo {
	log.Printf("Not implemented GetIssueInfoByID call")
	return IssueInfo{}
}

func GetHistoryInfoByID(id int) HistoryInfo {
	log.Printf("Not implemented GetHistoryInfoByID call")
	return HistoryInfo{}
}

func GetProjectInfoByID(id int) ProjectInfo {
	log.Printf("Not implemented GetProjectByID call")
	return ProjectInfo{}
}

func PutProjectToDB(data ProjectInfo) error {
	log.Printf("Not implemented PutProjectToDB call")
	return nil
}

func PutHistoryToDB(data HistoryInfo) error {
	log.Printf("Not implemented PutHistoryToDB call")
	return nil
}

func PutIssueToDB(data IssueInfo) error {
	log.Printf("Not implemented PutIssueToDB call")
	return nil
}
