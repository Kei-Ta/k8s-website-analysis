package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Kei-Ta/k8s-website-analysis/internal/action"
	"github.com/urfave/cli/v2"
)

type App struct {
	Cli       *cli.App //ポインタ型
	Language  string
	Directory string
	Tag       string
}

func NewApp() *App {
	app := App{} //構造体作成
	app.Cli = &cli.App{
		Name:  "k8swebsite-diff",
		Usage: "Check and manage folders",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "clone kubernetes/website",
				Action: func(c *cli.Context) error {
					action.Init()
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "update kubernetes/website",
				Action: func(c *cli.Context) error {
					action.Update()
					return nil
				},
			},
			{
				Name:  "diff",
				Usage: "diff kubernetes/website",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "language",
						Aliases:     []string{"l"},
						Usage:       "Select language",
						Required:    true,
						Destination: &app.Language,
					},
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Select directory(Optional)",
						Destination: &app.Directory,
					},
					&cli.StringFlag{
						Name:        "tag",
						Aliases:     []string{"t"},
						Usage:       "Select tag(Optional)",
						Destination: &app.Tag,
					},
				},
				Before: validateLanguage,
				Action: func(c *cli.Context) error {
					action.Diff(app.Language, app.Directory, app.Tag)
					return nil
				},
			},
			{
				Name:  "coverage",
				Usage: "coverage kubernetes/website",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "language",
						Aliases:     []string{"l"},
						Usage:       "Select language",
						Required:    true,
						Destination: &app.Language,
					},
				},
				Before: validateLanguage,
				Action: func(c *cli.Context) error {
					fmt.Printf("%s", app.Language)
					return nil
				},
			},
		},
	}
	return &app
}
func (a *App) Run(ctx context.Context) error {
	return a.Cli.RunContext(ctx, os.Args)
}

func validateLanguage(c *cli.Context) error {
	language := c.String("language")
	supportedLanguages := []string{"bn", "de", "es", "fr", "hi", "id", "it", "ja", "ko", "no", "pl", "pt-br", "ru", "uk", "vi", "zh-cn"}
	for _, l := range supportedLanguages {
		if language == l {
			return nil
		}
	}
	return fmt.Errorf("invalid language: %s. Supported languages are: English, Spanish, French", language)
}
