package internal

import (
	"github.com/montruh-afk/gator/internal/database"
)


type State struct {
	Configuration *Config
	Db *database.Queries
}