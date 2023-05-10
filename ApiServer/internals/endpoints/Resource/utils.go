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
	var issue = IssueInfo{}
	err := db.QueryRow(
		"SELECT "+
			"id,"+
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
			"WHERE id = $1", id,
	).Scan(
		&issue.IssueID,
		&issue.ProjectID, &issue.AuthorID, &issue.Key, &issue.Summary,
		&issue.Description, &issue.Type, &issue.Priority, &issue.Status,
		&issue.CreatedTime, &issue.ClosedTime, &issue.UpdatedTime, &issue.TimeSpent,
	)

	if err != nil {
		log.Printf("Error with querying an issue with id = %d", id)
		return IssueInfo{}, err
	}

	issue.IssueID = id
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
		"SELECT "+
			"authorId,"+
			"EXTRACT(EPOCH FROM changeTime)::bigint,"+
			"fromStatus,"+
			"toStatus "+
			"FROM StatusChanges "+
			"WHERE issueId = $1", id,
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
		var statusChange = HistoryInfo{}
		err := rows.Scan(&statusChange.AuthorID, &statusChange.ChangeTime, &statusChange.FromStatus, &statusChange.ToStatus)
		if err != nil {
			log.Printf("Error on handling query to the database: %s", err.Error())
			return history, err
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

	var project = ProjectInfo{}

	err := db.QueryRow(
		"SELECT "+
			"Projects.id,"+
			"Projects.title, "+
			"("+
			"SELECT COUNT(*) from Issue WHERE Issue.projectId=$1"+
			") as issues_count "+
			"FROM Projects "+
			"WHERE Projects.id = $1", id,
	).Scan(
		&project.ProjectID,
		&project.Title,
		&project.IssuesCount,
	)

	if err != nil {
		log.Printf("Error with querying an project with id = %d", id)
		return ProjectInfo{}, err
	}

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
		"INSERT INTO Projects (title) VALUES (COALESCE($1, '')) RETURNING id", data.Title,
	).Scan(&newID)

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

	err := db.QueryRow("INSERT INTO StatusChanges ("+
		"issueId,authorId,changeTime,fromStatus,toStatus) VALUES "+
		"($1, $2, now(), COALESCE($3, ''), COALESCE($4, ''));",
		data.IssueID, data.AuthorID, data.FromStatus, data.ToStatus,
	).Err()

	if err != nil {
		log.Printf("Error with creating history entry: %s", err.Error())
		return err
	}

	err = db.QueryRow("UPDATE Issue SET status = COALESCE($1, status), updatedTime = now(), timespent = EXTRACT(EPOCH FROM now()-createdTime)::integer WHERE id = $2", data.ToStatus, data.IssueID).Err()

	if err != nil {
		log.Printf("Error with updating issue entry: %s", err.Error())
		return err
	}

	log.Printf("PutHistoryToDB call")
	return nil
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
		"INSERT INTO Issue (projectId,authorId,assigneeId,key,summary,description,type,"+
			"priority,status,createdTime,closedTime,updatedTime,timeSpent) VALUES ("+
			"$1, $2, $3, $4, $5, $6, $7, $8, $9, to_timestamp($10), to_timestamp($11), to_timestamp($12), $13"+
			") RETURNING id",
		data.ProjectID, data.AuthorID, data.AssigneeId, data.Key, data.Summary, data.Description,
		data.Type, data.Priority, data.Status, data.CreatedTime, data.ClosedTime, data.UpdatedTime,
		data.TimeSpent,
	).Scan(&newID)

	log.Printf("PutIssueToDB call")
	return newID, err
}

func GetIssuesWithProjectId(projectId int, offset int, limit int) ([]IssueInfo, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var issues []IssueInfo
	rows, err := db.Query(
		"SELECT "+
			"Issue.id,"+
			"Issue.projectId,"+
			"Issue.authorId,"+
			"Issue.key,"+
			"Issue.summary,"+
			"Issue.description,"+
			"Issue.type,"+
			"Issue.priority,"+
			"Issue.status,"+
			"EXTRACT(EPOCH FROM Issue.createdTime)::bigint,"+
			"EXTRACT(EPOCH FROM Issue.closedTime)::bigint,"+
			"EXTRACT(EPOCH FROM Issue.updatedTime)::bigint,"+
			"Issue.timeSpent, "+
			"(SELECT COUNT(*) FROM statuschanges WHERE statuschanges.issueId=Issue.id) AS status_changes_count "+
			"FROM Issue "+
			"WHERE projectId=$1 "+
			"ORDER BY id "+
			"LIMIT $2 OFFSET $3",
		projectId, limit, offset,
	)

	if err != nil {
		log.Printf("Error with querying issues with projectId = %d", projectId)
		return issues, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	for rows.Next() {
		var issue = IssueInfo{}
		err := rows.Scan(
			&issue.IssueID,
			&issue.ProjectID, &issue.AuthorID, &issue.Key, &issue.Summary,
			&issue.Description, &issue.Type, &issue.Priority, &issue.Status,
			&issue.CreatedTime, &issue.ClosedTime, &issue.UpdatedTime, &issue.TimeSpent, &issue.ChangeStatusCount,
		)
		if err != nil {
			log.Printf("Error on handling query to the database: %s", err.Error())
			return issues, err
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

func GetProjectInfoByTitle(title string) (ProjectInfo, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var projectId int
	err := db.QueryRow(
		"SELECT "+
			"id "+
			"FROM Projects "+
			"WHERE Projects.title = $1", title,
	).Scan(
		&projectId,
	)

	if err != nil {
		log.Printf("Error with querying an project with title = %s", title)
		return ProjectInfo{}, err
	}

	log.Printf("GetProjectByID call")
	return GetProjectInfoByID(projectId)
}

func GetAllProjects(offset int, limit int) ([]ProjectInfo, error) {
	if db == nil {
		initDB()
	} else {
		log.Println("Try to re-establish database connection.")

		err := db.Ping()
		if err != nil {
			log.Fatalf("Can't connect to database.")
		}
	}

	var projects []ProjectInfo
	rows, err := db.Query(
		`SELECT 
    	projects.id, 
    	projects.title, 
    	COUNT(issue.id) AS issues_count
		FROM 
    	projects
		LEFT JOIN 
    	issue ON projects.id = issue.projectId
		GROUP BY 
    	projects.id, 
    	projects.title
		ORDER BY 
    	projects.id
		LIMIT 
    	$1
		OFFSET 
    	$2`,
		limit, offset,
	)

	if err != nil {
		log.Printf("Error with querying all projects.")
		return projects, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	for rows.Next() {
		var project = ProjectInfo{}
		err := rows.Scan(
			&project.ProjectID, &project.Title, &project.IssuesCount,
		)
		if err != nil {
			log.Printf("Error on handling query to the database: %s", err.Error())
			return projects, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}
