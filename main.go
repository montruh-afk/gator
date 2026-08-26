package main

import (
	"fmt"

	"github.com/montruh-afk/gator/internal/Config"
)

func main() {
	data, err := Config.Read()
	if err != nil {
		fmt.Errorf("Something went wrong: %v", err)
	}
	data.SetUser("montruh")
}
