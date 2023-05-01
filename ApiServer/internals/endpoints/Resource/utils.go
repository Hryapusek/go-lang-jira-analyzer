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

func GetIssueInfoByID(id int) (IssueInfo, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}
	var issue IssueInfo
	err := db.QueryRow(
		"SELECT "+
			"projectId,"+
			"authorId,"+
			"key,"+
			"summary,"+
			"description,"+
			"type,"+
			"priority,"+
			"status,"+
			"EXTRACT(EPOCH FROM createdTime)::bigint,"+
			"EXTRACT(EPOCH FROM closedTime)::bigint,"+
			"EXTRACT(EPOCH FROM updatedTime)::bigint,"+
			"timeSpent "+
			"FROM Issue "+
			"WHERE id = "+
			fmt.Sprintf("%d;", id),
	).Scan(
		&issue.ProjectID, &issue.AuthorID, &issue.Key, &issue.Summary,
		&issue.Description, &issue.Type, &issue.Priority, &issue.Status,
		&issue.CreatedTime, &issue.ClosedTime, &issue.UpdatedTime, &issue.TimeSpent,
	)

	if err != nil {
		log.Printf("Error with querying an issue with id = %d", id)
		return IssueInfo{}, err
	}

	log.Printf("Not implemented GetIssueInfoByID call")
	return issue, nil
}

func GetAllHistoryInfoByIssueID(id int) ([]HistoryInfo, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var history []HistoryInfo
	rows, err := db.Query(
		"SELECT " +
			"authorId," +
			"EXTRACT(EPOCH FROM changeTime)," +
			"fromStatus," +
			"toStatus " +
			"FROM StatusChanges " +
			"WHERE issueId = " +
			fmt.Sprintf("%d;", id),
	)

	if err != nil {
		log.Printf("Error with querying an history of issue with id = %d", id)
		return history, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	for rows.Next() {
		var statusChange HistoryInfo
		err := rows.Scan(&statusChange.AuthorID, &statusChange.ChangeTime, &statusChange.FromStatus, &statusChange.ToStatus)
		if err != nil {
			log.Printf("Error on handling query to the database: %s", err.Error())
			return nil, err
		}
		history = append(history, statusChange)
	}

	log.Printf("GetAllHistoryInfoByIssueID call")
	return history, nil
}

func GetProjectInfoByID(id int) (ProjectInfo, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var project ProjectInfo
	err := db.QueryRow(
		"SELECT " +
			"title " +
			"FROM Projects " +
			"WHERE id = " +
			fmt.Sprintf("%d;", id),
	).Scan(
		&project.Title,
	)

	if err != nil {
		log.Printf("Error with querying an project with id = %d", id)
		return ProjectInfo{}, err
	}

	project.ProjectID = id
	log.Printf("GetProjectByID call")
	return project, nil
}

func PutProjectToDB(data ProjectInfo) (int, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var newID int

	err := db.QueryRow(
		"INSERT INTO Projects (" +
			"title" +
			") VALUES (" +
			fmt.Sprintf("'%s'", data.Title) +
			") RETURNING id",
	).Scan(newID)

	log.Printf("PutProjectToDB call")
	return newID, err
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

	err := db.QueryRow("INSERT INTO StatusChanges (" +
		"issueId," +
		"authorId," +
		"changeTime," +
		"fromStatus," +
		"toStatus" +
		") VALUES " +
		fmt.Sprintf("("+
			"%d,"+
			"%d,"+
			"to_timestamp(%d),"+
			"'%s',"+
			"'%s'"+
			");", data.IssueID, data.AuthorID, data.ChangeTime, data.FromStatus, data.ToStatus,
		),
	).Err()

	log.Printf("PutHistoryToDB call")
	return err
}

func PutIssueToDB(data IssueInfo) (int, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var newID int
	err := db.QueryRow(
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
			") RETURNING id",
	).Scan(&newID)

	log.Printf("PutIssueToDB call")
	return newID, err
}
