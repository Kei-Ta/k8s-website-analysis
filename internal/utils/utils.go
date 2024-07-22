package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 指定されたフォルダ内の .md ファイル一覧と各ファイルの文字数を取得する関数
func ListMDFilesAndSizes(folderPath string) ([]string, error) {
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
			size, err := CountCharacters(path)
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
func CountCharacters(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	return len(content), nil
}

// 文字列スライス内に特定の文字列が含まれているかを確認する関数
func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// 指定されたファイルの単語数を数える関数
func CountWordsInFile(filePath string) (int, error) {
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

func FolderExists() bool {
	_, err := os.Stat("website")
	return err == nil
}

func AnalyzeFile(filename string) (int, string, error) {
	// ファイルを開く
	file, err := os.Open("website/content/en/" + filename)
	if err != nil {
		return 0, "", err // ファイルオープンに失敗した場合はエラーを返す
	}
	defer file.Close() // 関数終了時にファイルを確実に閉じる

	// ファイルを行単位で読み込むためのスキャナを作成
	scanner := bufio.NewScanner(file)
	lineCount := 0

	// スキャナを使ってファイルを一行ずつ読み込み、行数をカウントする
	for scanner.Scan() {
		lineCount++
	}

	// スキャン中にエラーが発生した場合はエラーを返す
	if err := scanner.Err(); err != nil {
		return 0, "", err
	}

	return lineCount, AssignTag(lineCount), nil
}

func AssignTag(lineCount int) string {
	switch {
	case lineCount < 10:
		return "XS"
	case lineCount < 30:
		return "S"
	case lineCount < 100:
		return "M"
	case lineCount < 500:
		return "L"
	default:
		return "XL"
	}
}
