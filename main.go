package main

import (
	"context"

	"github.com/sysnasri/nexenio/pkg/helpers"
)

func main() {

	cf := []string{"docker-compose.yml"}
	ctx := context.Background()
	s, _ := helpers.NewService(ctx)
	_, _ = s.Down(ctx, cf)
	_, _ = s.Up(ctx, cf)

}
