package main

import (
	"fmt"

	"github.com/montruh-afk/gator/internal"
)


func main() {
	data, err := internal.Read()
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}

	handler := make(map[string]func(*internal.State, internal.Command) error)

	state := &internal.State {
		Configuration: &data,
	}

	cmd := &internal.Commands{
		Handler: handler,
	}
	repl(state, cmd)
}
