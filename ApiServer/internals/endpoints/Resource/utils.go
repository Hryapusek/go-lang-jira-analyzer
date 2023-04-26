package endpoints

import (
	"ApiServer/internals/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB = nil

func init() {
	initDB()
}

func initDB() {
	cfg := config.LoadDBConfig("configs/server.yaml")

	connectionStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.UserDB,
		cfg.PasswordDB,
		cfg.HostDB,
		cfg.PortDB,
		cfg.NameDB,
	)

	var err error
	db, err = sql.Open("postgres", connectionStr)

	if err != nil {
		log.Fatalf("Unable to open Postgresql with %s database", connectionStr)
	}
}

func GetIssueInfoByID(id int) IssueInfo {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	log.Printf("Not implemented GetIssueInfoByID call")
	return IssueInfo{}
}

func GetHistoryInfoByID(id int) HistoryInfo {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	log.Printf("Not implemented GetHistoryInfoByID call")
	return HistoryInfo{}
}

func GetProjectInfoByID(id int) ProjectInfo {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	log.Printf("Not implemented GetProjectByID call")
	return ProjectInfo{}
}

func PutProjectToDB(data ProjectInfo) error {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	log.Printf("Not implemented PutProjectToDB call")
	return nil
}

func PutHistoryToDB(data HistoryInfo) error {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	log.Printf("Not implemented PutHistoryToDB call")
	return nil
}

func PutIssueToDB(data IssueInfo) error {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	_, err := db.Query(
		"INSERT INTO Issue (" +
			"projectId," +
			"authorId," +
			"assigneeId," +
			"key," +
			"summary," +
			"description," +
			"type," +
			"priority," +
			"status," +
			"createdTime," +
			"closedTime," +
			"updatedTime," +
			"timeSpent" +
			") VALUES (" +
			fmt.Sprintf(
				"%d,"+
					"%d,"+
					"%d,"+
					"'%s',"+
					"'%s',"+
					"'%s',"+
					"'%s',"+
					"'%s',"+
					"'%s',"+
					"to_timestamp(%d),"+
					"to_timestamp(%d),"+
					"to_timestamp(%d),"+
					"%d", data.ProjectID, data.AuthorID, data.AssigneeId, data.Key, data.Summary, data.Description,
				data.Type, data.Priority, data.Status, data.CreatedTime, data.ClosedTime, data.UpdatedTime,
				data.TimeSpent,
			) +
			")",
	)

	log.Printf("PutIssueToDB call")
	return err
}
