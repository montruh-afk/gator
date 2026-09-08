# Gator command-line interface
Gator is a simple blog aggregator that lets multiple users add, browse and keep up with feeds. It can also be used by a single user to group feed by different tastes. 

## Getting started
Gator requires both postgres and Go to be installed to run correctly.

### Installing postgres
###### Using macOS with brew:
type `brew install postgresql@15` in the terminal

###### Using linux / WSL (Debian):
`sudo apt update` 
`sudo apt install postgresql poatgresql-contrib`

###### Ensure the installation worked:
`psql --version`

### Installing Go
Docs available at [Installing Go](https://go.dev/doc/install).

## Installing Gator
Gator cli can be installed by running  `go install github.com/montruh-afk/gator` in the terminal.

## Config

Create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your database connection string.

## Running Gator
Gator can be executed using the command `gator`.
Most commands return the proper usage format when not used correctly.

## Gator Commands
#### register `username`
Creates a new user environment and returns user details on success.  

#### login `username`
Switches between multiple users.  

#### reset
Takes no arguments and clears all stored data.  

#### users
Takes no arguments and returns all currently registered users. It also indicated the currently logged in user.  

#### addfeed `feed name` `feed url`
Adds a new feed url to be tracked. Feed url must be a valid url.  

Any user adding a feed automatically follows the feed at that url.  

#### feeds
Takes no arguments and lists all feeds currently followed by a user.  

#### agg `cooldown period`
Continually updates on any new changes with a cooldown period e.g. 1m, 10m, 1h.  

#### follow `feed url`
Follows the provided feed url.  

#### unfollow `feed url`
Unfollows the provided feed url.  

#### following
Takes no arguments and lists all feeds currently being followed by a user.  

#### browse `[limit]`
Returns all posts under a followed feed with a set limit, defaults to 2 if no limit is provided.  

