package action

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Kei-Ta/k8s-website-analysis/internal/utils"
	"github.com/Kei-Ta/k8s-website-analysis/pkg/git"
)

const repoURL = "https://github.com/kubernetes/website.git"

func Init() {
	if utils.FolderExists() {
		fmt.Printf("websiteフォルダは既に存在します。\n")
		fmt.Println("Pullしますか？ (yes/no)")
		err := git.GitPull()
		if err != nil {
			fmt.Printf("Pullエラー: %v\n", err)
		} else {
			fmt.Println("Pullが完了しました。")
		}
	} else {
		fmt.Printf("websiteフォルダは存在しません。\n")
		fmt.Println("Cloneしますか？ (yes/no)")
		cmd := exec.Command("git", "clone", repoURL)
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Git clone failed: %s", err)
		}
	}
}

func Update() {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("入力エラー:", err)
	}
	input = strings.TrimSpace(strings.ToLower(input))
	switch input {
	case "yes", "y":
		fmt.Println("Pullを実行します。")
		err := git.GitPull()
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
}

// RunK8sWebsiteDiff は指定されたフォルダ内の .md ファイルを比較し、タグ付けを行う関数です。
func Diff(language string, directory string, tag string) {

	contentPath := "website/content"
	// ja フォルダと en フォルダのパスを設定します
	enPath := filepath.Join(contentPath, "en")
	selectLanPath := filepath.Join(contentPath, language)

	// ja フォルダと en フォルダの .md ファイル一覧と文字数を取得します
	selectLanPathFiles, err := utils.ListMDFilesAndSizes(selectLanPath)
	if err != nil {
		log.Fatalf("Failed to list .md files in ja folder: %v", err)
	}

	enFiles, err := utils.ListMDFilesAndSizes(enPath)
	if err != nil {
		log.Fatalf("Failed to list .md files in en folder: %v", err)
	}
	enFileCount := len(enFiles)

	// en フォルダにあって ja フォルダにない .md ファイルのタグ付けを行います

	c := 0

	for i := 0; i < len(enFiles); i++ {
		if !utils.Contains(selectLanPathFiles, enFiles[i]) {

			count, size, err := utils.AnalyzeFile(enFiles[i])
			if err != nil {
				log.Fatalf("Failed to count words in file: %v", err)
			}
			if tag == "" || tag == size {
				c++
				fmt.Printf("%s,Count: %d,Size: %s\n", enFiles[i], count, size)
			}
		}
	}
	fmt.Printf("no translate file count: %d\n", c)
	fmt.Printf("enFile count: %d\n", enFileCount)
}
