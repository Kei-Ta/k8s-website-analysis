package main

import (
	"context"
	"os"

	"github.com/Kei-Ta/k8s-website-analysis/internal/app"
)

func main() {

	app := app.NewApp()
	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		os.Exit(1)
	}
}
