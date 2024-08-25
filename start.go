package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
)

// startContainers take a docker compose file and start the services
func Listcontainers() {

	ctx := context.Background()

	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {

		log.Fatalf("Could not create a docker client! %v", err)

	}

	defer dc.Close()

	cli, err := command.NewDockerCli()
	if err != nil {
		log.Fatalf("cannot create a new DockerCli: %v", err)
	}

	clientOptions := flags.NewClientOptions()

	err = cli.Initialize(clientOptions)
	if err != nil {
		log.Fatalf("Cannot initialize docker %v", err)

	}
	cs := compose.NewComposeService(cli)

	projects, err := cs.List(ctx, api.ListOptions{})
	if err != nil {
		log.Fatalf("Could not list containers! %v", err)
	}

	_, err = cli.Client().Ping(ctx)
	if err != nil {
		log.Fatalf("Could not connect to docker daemon!: %v ", err)
	}

	for _, project := range projects {
		fmt.Printf("project: %s\n", project.ConfigFiles)
		// stop project containers

		err = cs.Down(ctx, project.Name, api.DownOptions{
			RemoveOrphans: true,
			Volumes:       false,
		})
		if err != nil {
			log.Fatalf("eror in downing project %v", err)
		}

	}

}
