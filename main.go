package main

import "github.com/sysnasri/nexenio/pkg/helpers"

func main() {

	helpers.ComposeProjectCreation("test", "./docker-compose.yml")
	//Listcontainers()

}
