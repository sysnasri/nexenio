package helpers

import (
	"context"
	"log"

	"github.com/compose-spec/compose-go/v2/cli"
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
	ctx := context.Background()
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	CheckErr("Could not create docker client", err)
	defer dc.Close()

	// New docker Cli to intract with docker
	dcli, err := command.NewDockerCli()
	CheckErr("Could not create a new Docker cli", err)
	if err != nil {
		return nil, nil, nil, err
	}

	err = dcli.Initialize(flags.NewClientOptions())
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
	CheckErr("Faild to initialize Docker client %v", err)

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

	ctx, _, dcli, err := DockerClientInit()
	CheckErr("Faild to initialize Docker client %v", err)

	//dcs := compose.NewComposeService(dcli)

	ops, err := cli.NewProjectOptions(
		[]string{cf},
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithName(p),
	)

	CheckErr("erro", err)
	pr, err := ops.LoadProject(ctx)
	CheckErr("error", err)
	//pryml, err := pr.MarshalYAML()
	//CheckErr("err", err)

	cs := compose.NewComposeService(dcli)
	cs.Up(ctx, pr, api.UpOptions{})
	//cs.Start(ctx, p, api.StartOptions{})

	//fmt.Println(string(pryml))

}
