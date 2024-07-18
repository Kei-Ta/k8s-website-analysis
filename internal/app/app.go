package app

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Kei-Ta/k8s-website-untranslated-finder/internal/action"
	"github.com/urfave/cli/v2"
)

type App struct {
	Cli     *cli.App //ポインタ型
	Profile string
}

func NewApp() *App {
	repoURL := "https://github.com/kubernetes/website.git"
	app := App{} //構造体作成
	app.Cli = &cli.App{
		Name:  "k8swebsite-diff",
		Usage: "Check and manage folders",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "clone kubernetes/website",
				Action: func(c *cli.Context) error {
					if folderExists() {
						fmt.Printf("websiteフォルダは既に存在します。\n")
						fmt.Println("Pullしますか？ (yes/no)")
						cmd := exec.Command("cd", "website")
						err := cmd.Run()
						if err != nil {
							log.Fatalf("Git clone failed: %s", err)
						}
						cmd = exec.Command("git", "pull", repoURL)
						err = cmd.Run()
						if err != nil {
							log.Fatalf("Git clone failed: %s", err)
						}
					} else {
						fmt.Printf("websiteフォルダは存在しません。\n")
						fmt.Println("Cloneしますか？ (yes/no)")

						cmd := exec.Command("git", "pull", repoURL)
						err := cmd.Run()
						if err != nil {
							log.Fatalf("Git clone failed: %s", err)
						}
					}
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "update kubernetes/website",
				Action: func(c *cli.Context) error {
					reader := bufio.NewReader(os.Stdin)
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("入力エラー:", err)
						return nil
					}
					fmt.Println("test")
					input = strings.TrimSpace(strings.ToLower(input))
					switch input {
					case "yes", "y":
						fmt.Println("Pullを実行します。")
						err := gitPull("website")
						if err != nil {
							fmt.Printf("Pullエラー: %v\n", err)
						} else {
							fmt.Println("Pullが完了しました。")
						}
					case "no", "n":
						fmt.Println("操作をキャンセルしました。")
					default:
						fmt.Println("無効な入力です。操作をキャンセルしました。")
					}
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
						Destination: &app.Profile,
					},
				},
				Before: validateLanguage,
				Action: func(c *cli.Context) error {

					action.Diff(app.Profile)
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
						Destination: &app.Profile,
					},
				},
				Before: validateLanguage,
				Action: func(c *cli.Context) error {
					fmt.Printf("%s", app.Profile)
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
	supportedLanguages := []string{"ja", "ko", "Fc"}
	for _, l := range supportedLanguages {
		if language == l {
			return nil
		}
	}
	return fmt.Errorf("invalid language: %s. Supported languages are: English, Spanish, French", language)
}
func folderExists() bool {
	_, err := os.Stat("website")
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func gitPull(folderPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = folderPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, output)
	}

	return nil
}
