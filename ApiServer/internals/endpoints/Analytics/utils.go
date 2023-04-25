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

func GraphOne(projectName string) []IssueForGraphOne {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	rows, err := db.Query(
		"SELECT " +
			" i.id," +
			" i.summary," +
			" i.createdTime," +
			" i.closedTime," +
			" (EXTRACT(EPOCH FROM (i.closedTime)) - EXTRACT(EPOCH FROM (i.createdTime))) AS time_open_seconds" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			" i.status = 'Closed'" +
			fmt.Sprintf(" AND p.title = '%s'", projectName) +
			" ORDER BY" +
			" time_open_seconds;",
	)
	if err != nil {
		log.Printf("Unable to query a database with /api/v1/graph/1 route: %s", err.Error())
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	var issues []IssueForGraphOne

	for rows.Next() {
		var issue IssueForGraphOne
		err := rows.Scan(&issue.Id, &issue.Summary, &issue.CreatedTime, &issue.ClosedTime, &issue.TimeOpenSeconds)
		if err != nil {
			log.Fatal(err)
		}
		issues = append(issues, issue)
	}
	log.Printf("We have result on /api/v1/graph/1 route!")
	return issues
}
