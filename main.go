package main

import (
	"jwt-authentication/datatabase"
	"jwt-authentication/routes"
)

func main() {

	datatabase.DbInit()


	server :=routes.SetupRouter()

	server.Run(":5000")
}