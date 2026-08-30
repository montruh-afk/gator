package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/montruh-afk/gator/internal"
	"github.com/montruh-afk/gator/internal/database"
)


func main() {
	data, err := internal.Read()
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}

	handlers := make(map[string]func(*internal.State, internal.Command) error)
	db, err := sql.Open("postgres", data.Url)
	dbQueries := database.New(db)

	state := &internal.State {
		Configuration: &data,
		Db: dbQueries,
	}
	

	cmd := &internal.Commands{
		Handlers: handlers,
	}
	start(state, cmd)
}
