package main

import (
	"fmt"

	"github.com/montruh-afk/gator/internal"
)

func main() {
	data, err := config.Read()
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}
	data.SetUser("montruh")
	fmt.Println(data)
}
