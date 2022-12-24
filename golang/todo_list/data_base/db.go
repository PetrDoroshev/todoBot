package data_base

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DB struct {
	conn *sql.DB
}

type Task struct {
	TaskID      int
	UserID      int64
	Description string
	Priority    string
	NotifyTime  *string
	Status      string
}

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("HOST"),
	os.Getenv("PORT"),
	os.Getenv("USER"),
	os.Getenv("PASSWORD"),
	os.Getenv("DBNAME"))

func NewDB() *DB {
	conn, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../../go/build/todo_list/migrations/1_initial.up.sql", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	return &DB{conn: conn}
}

func (db *DB) CloseDB() {
	err := db.conn.Close()

	if err != nil {
		panic(err)
	}
}

func (db *DB) AddNewTask(newTask *Task) error {

	insertStatement := "INSERT INTO tasks(user_id, description, priority, notify_time, status) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.conn.Exec(insertStatement, newTask.UserID, newTask.Description, newTask.Priority,
		newTask.NotifyTime, newTask.Status)

	if err != nil {
		return fmt.Errorf("unable to insert new task: %w", err)
	}

	return nil
}

func (db *DB) ListUserTasks(userID int64) ([]Task, error) {

	taskList := make([]Task, 0, 0)

	rows, err := db.conn.Query("SELECT * FROM tasks WHERE user_id = $1 AND status = 'in work' ORDER BY id", userID)
	if err != nil {
		return nil, fmt.Errorf("unable to extract data from table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err = rows.Scan(&t.TaskID, &t.UserID, &t.Description, &t.Priority, &t.NotifyTime, &t.Status)

		if err != nil {
			return nil, fmt.Errorf("unable to scan data to struct: %w", err)
		}

		taskList = append(taskList, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return taskList, nil
}

func (db *DB) ListTaskWithPriority(userID int64, priority string) ([]Task, error) {
	taskList := make([]Task, 0, 0)

	rows, err := db.conn.Query("SELECT * FROM tasks WHERE user_id = $1 AND priority = $2 ORDER BY id", userID, priority)
	if err != nil {
		return nil, fmt.Errorf("unable to extract data from table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := Task{}
		err = rows.Scan(&t.TaskID, &t.UserID, &t.Description, &t.Priority, &t.NotifyTime, &t.Status)

		if err != nil {
			return nil, fmt.Errorf("unable to scan data to struct: %w", err)
		}

		taskList = append(taskList, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return taskList, nil
}

func (db *DB) ListFinishedTasks(userID int64) ([]Task, error) {
	taskList := make([]Task, 0, 0)

	rows, err := db.conn.Query("SELECT * FROM tasks WHERE status = 'finished' AND user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, fmt.Errorf("unable to extract data from table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := Task{}
		err = rows.Scan(&t.TaskID, &t.UserID, &t.Description, &t.Priority, &t.NotifyTime, &t.Status)

		if err != nil {
			return nil, fmt.Errorf("unable to scan data: %w", err)
		}

		taskList = append(taskList, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return taskList, nil
}

func (db *DB) FinishTask(userID int64, taskNumber int) error {

	updateStatement := "UPDATE tasks SET status = 'finished' " +
		"WHERE id = (SELECT id " +
		"FROM (SELECT id, row_number() over (order by id) " +
		"FROM tasks " +
		"WHERE user_id = $1 AND status = 'in work' ) as num_row_table " +
		"where row_number = $2)"

	_, err := db.conn.Exec(updateStatement, userID, taskNumber)

	if err != nil {
		return fmt.Errorf("unable to update status : %w", err)
	}

	return nil
}

func (db *DB) CountFinishedTasks(userID int64) (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT count(id) FROM tasks WHERE status = 'finished' AND user_id = $1", userID).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("unable to extract data from table: %w", err)
	}
	return count, nil
}
