package main

import (
	"github.com/matjam/grupal/api"
	"github.com/matjam/grupal/database"
	"github.com/matjam/grupal/routes"
)

// 	dsn := fmt.Sprintf"host=localhost user=postgres password=mysecretpassword dbname=grupal " +
//		"port=5432 sslmode=disable TimeZone=America/Los_Angeles"

func main() {
	db := database.NewDB("localhost", "postgres", "mysecretpassword", "grupal", 5432, []any{&api.User{}})

	routes.Start(db)
}
