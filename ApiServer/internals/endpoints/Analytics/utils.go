package endpoints

// TODO: надо оптимизировать запросы 5-6 (тесты в админке работают быстрее, а там оптимизатор)

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
	// Гистограмма, отражающая время, которое задачи провели в открытом состоянии (время в секундах) и только для закрытых
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
			" FLOOR((EXTRACT(EPOCH FROM (i.closedTime)) - EXTRACT(EPOCH FROM (i.createdTime))))::bigint AS time_open_seconds" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			" i.status IN ('Closed', 'Resolved')" +
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
		err := rows.Scan(&issue.Id, &issue.TimeOpenedSeconds)
		if err != nil {
			log.Fatal(err)
		}
		issues = append(issues, issue)
	}
	log.Printf("We have result on /api/v1/graph/1 route!")
	return issues
}

func GraphTwo(projectName string) []IssueForGraphTwo {
	// Диаграмма, демонстрирующая распределение времени по состоянию "Open" (я так понимаю отсортировать issues по открытым и времени)
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
		" SELECT" +
			" i.id," +
			" i.timespent AS time_open_seconds" + // timespent === FLOOR((EXTRACT(EPOCH FROM (now()::timestamp)) - EXTRACT(EPOCH FROM (i.createdTime))))::bigint
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			" i.status = 'Open'" +
			fmt.Sprintf(" AND p.title = '%s'", projectName) +
			" ORDER BY" +
			" time_open_seconds",
	)
	if err != nil {
		log.Printf("Unable to query a database with /api/v1/graph/2 route: %s", err.Error())
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	var issues []IssueForGraphTwo

	for rows.Next() {
		var issue IssueForGraphTwo
		err := rows.Scan(&issue.Id, &issue.TimeOpenSeconds)
		if err != nil {
			log.Fatal(err)
		}
		issues = append(issues, issue)
	}
	log.Printf("We have result on /api/v1/graph/2 route!")
	return issues
}

func GraphThree(projectName string) []GraphThreeData {
	// Здесь строки с датой, где был создан и/или закрыт хотя бы один issue.
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
		"WITH created_issues AS (" +
			" SELECT" +
			" DATE(i.createdTime) AS date," +
			" COUNT(*) AS created_count" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			fmt.Sprintf(" p.title = '%s'", projectName) +
			" GROUP BY" +
			" DATE(i.createdTime)" +
			" )," +
			" closed_issues AS (" +
			" SELECT" +
			" DATE(i.closedTime) AS date," +
			" COUNT(*) AS closed_count" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			fmt.Sprintf(" p.title = '%s'", projectName) +
			" AND i.status = 'Closed'" +
			" GROUP BY" +
			" DATE(i.closedTime)" +
			" )" +
			" SELECT" +
			" FLOOR(EXTRACT(EPOCH FROM (COALESCE(ci.date, cl.date))::timestamp))::bigint AS unix_date," +
			" COALESCE(ci.created_count, 0) AS created_issues," +
			" COALESCE(cl.closed_count, 0) AS closed_issues" +
			" FROM" +
			" created_issues ci" +
			" FULL OUTER JOIN" +
			" closed_issues cl ON ci.date = cl.date" +
			" ORDER BY" +
			" unix_date;",
	)
	if err != nil {
		log.Printf("Unable to query a database with /api/v1/graph/3 route: %s", err.Error())
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	var data []GraphThreeData

	for rows.Next() {
		var entry GraphThreeData
		err := rows.Scan(&entry.Date, &entry.CreateIssues, &entry.ClosedIssues)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, entry)
	}
	log.Printf("We have result on /api/v1/graph/3 route!")
	return data
}

func GraphFour(projectName string) []GraphFourData {
	// График по типу задачи (сложность видимо).
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
		"SELECT" +
			" i.type AS issue_type," +
			" COUNT(*) AS issue_count" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			fmt.Sprintf(" p.title = '%s'", projectName) +
			" GROUP BY" +
			" i.type" +
			" ORDER BY" +
			" i.type;",
	)
	if err != nil {
		log.Printf("Unable to query a database with /api/v1/graph/4 route: %s", err.Error())
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	var data []GraphFourData

	for rows.Next() {
		var entry GraphFourData
		err := rows.Scan(&entry.Type, &entry.Count)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, entry)
	}
	log.Printf("We have result on /api/v1/graph/4 route!")
	return data
}

func GraphFive(projectName string) []GraphFiveAndSixData {
	// Подсчет кол-ва задач по статусу (приоритет).
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
		"SELECT" +
			" i.priority," +
			" COUNT(*) AS issue_count" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			fmt.Sprintf(" p.title = '%s' AND", projectName) +
			" i.priority IN ('Minor', 'Major', 'Blocker', 'Critical')" +
			" GROUP BY" +
			" i.priority" +
			" ORDER BY" +
			" i.priority;",
	)
	if err != nil {
		log.Printf("Unable to query a database with /api/v1/graph/5 route: %s", err.Error())
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	var data []GraphFiveAndSixData

	for rows.Next() {
		var entry GraphFiveAndSixData
		err := rows.Scan(&entry.PriorityType, &entry.Count)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, entry)
	}
	log.Printf("We have result on /api/v1/graph/5 route!")
	return data
}

func GraphSix(projectName string) []GraphFiveAndSixData {
	// Подсчет кол-ва задач по статусу (приоретет закрытых).
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
		"SELECT" +
			" i.priority," +
			" COUNT(*) AS issue_count" +
			" FROM" +
			" Issue i" +
			" JOIN" +
			" Projects p ON p.id = i.projectId" +
			" WHERE" +
			fmt.Sprintf(" p.title = '%s' AND", projectName) +
			" i.status = 'Closed' AND" +
			" i.priority IN ('Minor', 'Major', 'Blocker', 'Critical')" +
			" GROUP BY" +
			" i.priority" +
			" ORDER BY" +
			" i.priority;",
	)
	if err != nil {
		log.Printf("Unable to query a database with /api/v1/graph/6 route: %s", err.Error())
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to Close() on rows.")
		}
	}(rows)

	var data []GraphFiveAndSixData

	for rows.Next() {
		var entry GraphFiveAndSixData
		err := rows.Scan(&entry.PriorityType, &entry.Count)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, entry)
	}
	log.Printf("We have result on /api/v1/graph/6 route!")
	return data
}
