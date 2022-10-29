package store

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Store struct {
	conn *sqlx.DB
}

type People struct {
	ID   int32
	Name string
}

func NewStore(connString string) *Store {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := pgx.WithInstance(conn.DB, &pgx.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../Golang/HomeWork_db/golang/people_service/migrations/",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	return &Store{
		conn: conn,
	}
}

func (s *Store) ListPeople() ([]People, error) {

	people := make([]People, 0, 0)
	var id int32
	var name string

	rows, err := s.conn.Query("select * from people")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("error while reding: %d", err)
		}

		people = append(people, People{
			ID:   id,
			Name: name,
		})
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in rows: %v\n", err)
	}

	return people, nil
}

func (s *Store) GetPeopleByID(id string) (People, error) {

	var name string
	var _id int32

	row := s.conn.QueryRow("select * from people where id = " + id)

	err := row.Scan(&_id, &name)
	if err != nil {
		return People{}, err
	}
	return People{ID: _id, Name: name}, nil
}
