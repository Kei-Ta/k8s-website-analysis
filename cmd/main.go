package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Kei-Ta/k8s-website-untranslated-finder/internal/app"
)

func main() {

	app := app.NewApp()
	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		fmt.Println("入力エラー:", err)
		log.Fatalf("app faild")
		os.Exit(1)
	}
}
