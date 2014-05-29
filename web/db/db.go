package main

import (
	r "github.com/dancannon/gorethink"
	"log"
	"os"
	"time"
)

func InitDB() *r.Session {

	session, err := r.Connect(r.ConnectOpts{
		Address:     os.Getenv("RETHINKDB_URL"),
		Database:    "test",
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})

	if err != nil {
		log.Println(err)
	}
	err = r.DbCreate("test").Exec(session)
	if err != nil {
		log.Println(err)
	}

	_, err = r.Db("test").TableCreate("todos").RunWrite(session)
	if err != nil {
		log.Println(err)
	}

	return session
}
