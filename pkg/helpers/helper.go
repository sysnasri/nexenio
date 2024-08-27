package helpers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
)

func CheckErr(s string, err error) {

	if err != nil {
		log.Fatalf("%s : %v", s, err)
	}

}

// DockerClientInit initialize docker client
func DockerClientInit() (context.Context, *client.Client, *command.DockerCli, error) {

	// Create a new context in background to connect to docker api
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	CheckErr("Could not create docker client", err)
	//defer dc.Close()

	// New docker Cli to intract with docker
	dcli, err := command.NewDockerCli()
	CheckErr("Could not create a new Docker cli", err)
	if err != nil {
		return nil, nil, nil, err
	}

	err = dcli.Initialize(&flags.ClientOptions{})
	CheckErr("erro in dcli init", err)

	_, err = dcli.Client().Ping(ctx)
	CheckErr("Could not ping docker!", err)
	if err != nil {
		return nil, nil, nil, err
	}

	return ctx, dc, dcli, nil

}

// ListComposeProject lists docker stack projects
func ListComposeProject() []api.Stack {

	ctx, _, dcli, err := DockerClientInit()
	CheckErr("Faild to initialize Docker client", err)

	// Creates docker Compose service to make a compose file up/down!
	dcs := compose.NewComposeService(dcli)

	lo := api.ListOptions{
		All: true,
	}

	// List of Compose projects
	pl, err := dcs.List(ctx, lo)
	CheckErr("Could not list projects", err)

	return pl

}

func ComposeProjectCreation(p string, cf string) {

	ctx, cancel := context.WithCancel(context.Background())

	dcli, err := command.NewDockerCli()
	CheckErr("Could not create a new Docker cli", err)

	err = dcli.Initialize(&flags.ClientOptions{})
	CheckErr("Error initializing Docker CLI", err)

	ops, err := cli.NewProjectOptions([]string{cf}, cli.WithWorkingDirectory("/Users/nasri/nexenio/"))
	CheckErr("Error creating project options", err)

	pr, err := ops.LoadProject(ctx)
	for i, s := range pr.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     pr.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(pr.ComposeFiles, ","),
			api.OneoffLabel:      "False", // default, will be overridden by `run` command
		}
		pr.Services[i] = s
	}

	CheckErr("Error loading project", err)

	cs := compose.NewComposeService(dcli)
	err = cs.Create(ctx, pr, api.CreateOptions{Services: pr.ServiceNames()})
	CheckErr("Failed to create the project", err)

	err = cs.Start(ctx, pr.Name, api.StartOptions{Project: pr, Services: pr.ServiceNames()})
	CheckErr("Failed to start the project", err)
	defer cancel()
}

func ComposeProjectCreation1(p string, cf string) {

	// ctx, _, dcli, err := DockerClientInit()
	// CheckErr("Faild to initialize Docker client %v", err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	CheckErr("Could not create docker client", err)

	//dcs := compose.NewComposeService(dcli)

	ops, err := cli.NewProjectOptions(
		[]string{cf},
		cli.WithWorkingDirectory("/Users/nasri/nexenio/"),
	)
	CheckErr("erro", err)

	dcli, err := command.NewDockerCli()
	CheckErr("erro", err)

	err = dcli.Initialize(&flags.ClientOptions{})
	CheckErr("error init docker", err)
	pr, err := ops.LoadProject(ctx)
	c, b := pr.ComposeFiles, pr.Configs

	fmt.Println(c, b)
	CheckErr("error", err)
	pryml, err := pr.MarshalYAML()
	CheckErr("err", err)
	fmt.Println(string(pryml))

	cs := compose.NewComposeService(dcli)
	err = cs.Create(ctx, &types.Project{
		Name:     "nexenio",
		Services: pr.Services,
	}, api.CreateOptions{
		Services: pr.ServiceNames(),
	})
	CheckErr("Failed to bring up the project", err)
	err = cs.Start(ctx, "nexenio", api.StartOptions{
		Project:  pr,
		Services: pr.ServiceNames(),
	})
	CheckErr("Failed to bring up the project", err)
}

//cs.Start(ctx, p, api.StartOptions{})

//fmt.Println(string(pryml))
