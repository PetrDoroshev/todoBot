package main

import (
	"fmt"
	"github.com/PetrDoroshev/HomeWork_db/golang/people_service/service/store"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	s := store.NewStore("postgres://doroshev:doroshev@95.217.232.188:7777/doroshev")
	people, _ := s.ListPeople()
	fmt.Println(people)

	person, err := s.GetPeopleByID("9")

	fmt.Println(person)
	fmt.Print(err)

}
