package action

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// RunK8sWebsiteDiff は指定されたフォルダ内の .md ファイルを比較し、タグ付けを行う関数です。
func Diff(language string) {

	contentPath := "website/content"
	// ja フォルダと en フォルダのパスを設定します
	enPath := filepath.Join(contentPath, "en")
	selectLanPath := filepath.Join(contentPath, language)

	// ja フォルダと en フォルダの .md ファイル一覧と文字数を取得します
	selectLanPathFiles, err := listMDFilesAndSizes(selectLanPath)
	if err != nil {
		log.Fatalf("Failed to list .md files in ja folder: %v", err)
	}

	enFiles, err := listMDFilesAndSizes(enPath)
	if err != nil {
		log.Fatalf("Failed to list .md files in en folder: %v", err)
	}

	// en フォルダにあって ja フォルダにない .md ファイルのタグ付けを行います

	for i := 0; i < len(enFiles); i++ {
		if !contains(selectLanPathFiles, enFiles[i]) {
			count, err := countWordsInFile(enFiles[i])
			if err != nil {
				log.Fatalf("Failed to count words in file: %v", err)
			}
			fmt.Printf("File: %s,Count: %d\n", enFiles[i], count)
		}
	}
}

// 指定されたフォルダ内の .md ファイル一覧と各ファイルの文字数を取得する関数
func listMDFilesAndSizes(folderPath string) ([]string, error) {
	var mdFiles []string
	var sizes []int

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// ファイルのみ処理します
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			relPath, err := filepath.Rel(folderPath, path)
			if err != nil {
				return err
			}
			mdFiles = append(mdFiles, relPath)

			// ファイルの文字数を取得します
			size, err := countCharacters(path)
			if err != nil {
				return fmt.Errorf("failed to count characters for file %s: %v", relPath, err)
			}
			sizes = append(sizes, size)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mdFiles, nil
}

// 指定されたファイルの文字数を数える関数
func countCharacters(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	return len(content), nil
}

// 文字列スライス内に特定の文字列が含まれているかを確認する関数
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// 文字数に応じてタグを付ける関数
func tagBySize(size int) string {
	if size <= 3000 {
		return "S"
	} else if size <= 5000 {
		return "M"
	} else if size <= 10000 {
		return "L"
	} else {
		return "XL"
	}
}

// 指定されたファイルの単語数を数える関数
func countWordsInFile(filePath string) (int, error) {
	file, err := os.Open("website/content/en/" + filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordCount := 0
	for scanner.Scan() {
		wordCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan error: %v", err)
	}

	return wordCount, nil
}
