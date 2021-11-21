package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	//"strings"

	_ "github.com/mattn/go-sqlite3"
)

//implements the storage interface.

type Storage struct {
	db *sql.DB
}

type DataEntry struct {
	name string
	url  string
	desc string
}

func (s *Storage) init() {
	log.Println("Opening DB")
	alreadyExists := true
	if _, err := os.Stat("./data/data.db"); errors.Is(err, os.ErrNotExist) {
		alreadyExists = false
	}
	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatal(err)
	}
	//*only* initialize the table if already exists is false
	if !alreadyExists {
		log.Println("Initializing agdata table")
		sql := `
			CREATE TABLE agdata (
				name text,
				url text,
				desc text
			)
		`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}
	s.db = db
}

func (s *Storage) addEntry(name, url, desc string) {
	_, err := s.db.Exec(fmt.Sprintf("INSERT INTO agdata(name,url,desc) values (\"%v\",\"%v\",\"%v\")", name, url, desc))
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Storage) getEntries() []DataEntry {
	var entries []DataEntry
	rows, err := s.db.Query("select name, url, desc from agdata")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var url string
		var desc string
		err := rows.Scan(&name, &url, &desc)
		if err != nil {
			log.Fatal(err)
		}
		entries = append(entries, DataEntry{
			name, url, desc,
		})

	}
	return entries
}
