package action

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Kei-Ta/k8s-website-analysis/internal/utils"
	"github.com/Kei-Ta/k8s-website-analysis/pkg/git"
)

func Init() {
	if utils.FolderExists() {
		fmt.Printf("kubernetes/website project exist.\n")
		fmt.Printf("If you update project, Please run the update command.\n")
	} else {
		fmt.Printf("kubernetes/website project doesn't exist.\n")
		fmt.Println("run git clone.")
		err := git.GitClone()
		if err != nil {
			log.Fatalf("Git clone failed: %s", err)
		} else {
			fmt.Printf("Success init command.\n")
		}
	}
}

func Update() {
	if !utils.FolderExists() {
		fmt.Printf("kubernetes/website project doesn't exist.\n")
		fmt.Printf("Please run the init command.\n")
		return
	}
	err := git.GitPull()
	if err != nil {
		log.Fatalf("Git pull failed: %s", err)
	} else {
		fmt.Printf("Success update command.\n")
	}

}

// RunK8sWebsiteDiff は指定されたフォルダ内の .md ファイルを比較し、タグ付けを行う関数です。
func Diff(language string, directory string, tag string) {
	if !utils.FolderExists() {
		fmt.Printf("kubernetes/website project doesn't exist.\n")
		fmt.Printf("Please run the init command.\n")
		return
	}
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

	// 出力ディレクトリのパス
	outputDir := "output"

	// ディレクトリが存在しない場合は作成
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	// 出力ファイルを作成します
	outputFile, err := os.Create(outputDir + "/diff_" + language + ".txt")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// en フォルダにあって ja フォルダにない .md ファイルのタグ付けを行います

	c := 0
	// ファイルサイズごとにファイルを集計するためのマップ
	sizeMap := map[string]int{
		"XS": 0,
		"S":  0,
		"M":  0,
		"L":  0,
		"XL": 0,
	}
	// writer := tabwriter.NewWriter(outputFile, 0, 0, 1, ' ', tabwriter.AlignRight)

	for i := 0; i < len(enFiles); i++ {
		if !utils.Contains(selectLanPathFiles, enFiles[i]) {

			count, size, err := utils.AnalyzeFile(enFiles[i])
			if err != nil {
				log.Fatalf("Failed to count words in file: %v", err)
			}
			if tag == "" {
				c++
				sizeMap[size]++
				output := fmt.Sprintf("%s,Count: %d,Size: %s\n", enFiles[i], count, size)
				fmt.Print(output)                       // コンソールに出力
				_, err = outputFile.WriteString(output) // ファイルに出力
				if err != nil {
					log.Fatalf("Failed to write to output file: %v", err)
				}
			} else {
				c++
				output := fmt.Sprintf("%s,Count: %d,Size: %s\n", enFiles[i], count, size)
				fmt.Print(output)                       // コンソールに出力
				_, err = outputFile.WriteString(output) // ファイルに出力
				if err != nil {
					log.Fatalf("Failed to write to output file: %v", err)
				}
			}
		}
	}

	// 出力ファイルを作成します
	summaryFile, err := os.Create(outputDir + "/summary_" + language + ".md")
	if err != nil {
		log.Fatalf("Failed to create summary file: %v", err)
	}
	defer summaryFile.Close()

	// Markdown形式の表のヘッダーをログとファイルに出力
	header := "| XS | S | M | L | XL | Total |"
	separator := "|:------:|:------:|:------:|:------:|:------:|:------:|"
	fmt.Println(header)
	// fmt.Println(separator)
	fmt.Fprintln(summaryFile, header)
	fmt.Fprintln(summaryFile, separator)

	// サイズごとのカウントをログとファイルに出力
	countLine := fmt.Sprintf("| %d | %d | %d | %d | %d | %d |", sizeMap["XS"], sizeMap["S"], sizeMap["M"], sizeMap["L"], sizeMap["XL"], c)
	fmt.Println(countLine)
	fmt.Fprintln(summaryFile, countLine)

	fmt.Printf("no translate file count: %d\n", c)
	fmt.Printf("enFile count: %d\n", enFileCount)
}
