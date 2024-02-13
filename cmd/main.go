package main

import (
	"database/sql"
	"fmt"
	"go-people-crm/database"
	"go-people-crm/model"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const fileName = "sqlite.db"

func main() {
	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	personRepository := database.NewSQLiteRepository(db)

	if err := personRepository.Migrate(); err != nil {
		log.Fatal(err)
	}

	person := model.Person {
		Name: "Ornstein",
		LastTiemMet: time.Now(),
		MeetingFrequecy: time.Duration(time.Hour * 24 * 14),
	}

	person2 := model.Person {
		Name: "Havir",
		Phone: "(47) 99123-4567",
		LastTiemMet: time.Now(),
		MeetingFrequecy: time.Duration(time.Hour * 24 * 14),
	}

	createdPerson, err := personRepository.Create(person)
	if err != nil {
		log.Fatal(err)
	}

	createdPerson2, err := personRepository.Create(person2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(createdPerson)
	fmt.Println(createdPerson2)
}
