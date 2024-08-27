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
)

type Service struct {
	srv api.Service
	cli *command.DockerCli
}

// failOnError logs the given error with a message and returns a boolean indicating if an error occurred.
func failOnError(message string, err error) bool {
	if err != nil {
		log.Printf("%s : %v", message, err)
		return true

	}
	return false
}

func NewService(ctx context.Context) (*Service, error) {
	dcli, err := command.NewDockerCli()
	if failOnError("Error in creating Docker cli", err) {
		return nil, err
	}

	err = dcli.Initialize(&flags.ClientOptions{})

	if failOnError("Error in initializing docker client", err) {
		return nil, err
	}
	return &Service{
		cli: dcli,
		srv: compose.NewComposeService(dcli),
	}, nil

}

// DockerClientInit initializes a Docker CLI client.

// func (s *Service) DockerClientInit(ctx context.Context) (*command.DockerCli, error) {

// 	// New docker Cli to intract with docker
// 	dcli, err := command.NewDockerCli()
// 	if failOnError("Could not create a new Docker cli", err) {
// 		return nil, err
// 	}

// 	// Initialize docker cli
// 	err = dcli.Initialize(&flags.ClientOptions{})
// 	if failOnError("erro in dcli init", err) {
// 		return nil, err
// 	}

// 	// it pings docker to see if it is alive
// 	_, err = dcli.Client().Ping(ctx)
// 	if failOnError("Could not ping docker!", err) {
// 		return nil, err
// 	}

// 	return dcli, nil
// }

// Up brings up a Docker Compose project.

func (s *Service) Up(ctx context.Context, cf []string) ([]api.Stack, error) {

	return s.composeActions(ctx, cf, "up")
	// ops, err := cli.NewProjectOptions(cf)
	// failOnError("Error creating project options", err)

	// pr, err := ops.LoadProject(ctx)
	// s.ComposeProjectAddLabel(pr)

	// failOnError("Error loading project", err)

	// cs := compose.NewComposeService(dcli)
	// err = cs.Up(ctx, pr, api.UpOptions{})
	// failOnError("Failed to create the project", err)
	// return nil
}

// Down brings down a Docker Compose project.

func (s *Service) Down(ctx context.Context, cf []string) ([]api.Stack, error) {

	return s.composeActions(ctx, cf, "down")

	// ops, err := cli.NewProjectOptions(cf)
	// failOnError("Error creating project options", err)

	// pr, err := ops.LoadProject(ctx)

	// failOnError("Error loading project", err)

	// cs := compose.NewComposeService(dcli)
	// err = cs.Down(ctx, pr.Name, api.DownOptions{})
	// failOnError("Failed to create the project", err)
}

// List lists all Docker Compose projects.

func (s *Service) List(ctx context.Context, cf []string) ([]api.Stack, error) {
	return s.composeActions(ctx, cf, "list")

	// dsrv := compose.NewComposeService(dcli)
	// p, err := dsrv.List(ctx, api.ListOptions{})
	// if failOnError("error in listing compose", err) {
	// 	return nil, err
	// }

	// return p, nil
}

// addLabel handles the common logic for bringing up or down a Docker Compose project.

func (s *Service) addLabel(pr *types.Project) {

	for i, s := range pr.Services {
		if s.CustomLabels == nil {
			s.CustomLabels = map[string]string{}
		}
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
}

func (s *Service) composeActions(ctx context.Context, cf []string, act string) ([]api.Stack, error) {
	ops, err := cli.NewProjectOptions(cf)
	if failOnError("Error Creating Projects", err) {
		return nil, err
	}

	pr, err := ops.LoadProject(ctx)
	if failOnError("Error in Loading Projects", err) {
		return nil, err
	}
	s.addLabel(pr)

	var prl []api.Stack
	switch act {
	case "up":
		err := s.srv.Up(ctx, pr, api.UpOptions{})
		if failOnError("failed to execute up action on the project: ", err) {
			return nil, err

		}
	case "down":
		err := s.srv.Down(ctx, pr.Name, api.DownOptions{})
		if failOnError("failed to execute down action on the project: ", err) {
			return nil, err

		}
	case "list":
		prl, err = s.srv.List(ctx, api.ListOptions{})
		if failOnError("failed to execute list action on the project: ", err) {
			return nil, err

		}
	default:
		return nil, fmt.Errorf("unsupported action: %s", act)
	}

	return prl, nil

}
