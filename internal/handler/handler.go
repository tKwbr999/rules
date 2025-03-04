package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RULES_PATH環境変数を取得する
func GetRulesPath() (string, error) {
	rulesPath := os.Getenv("RULES_PATH")
	if rulesPath == "" {
		return "", fmt.Errorf("RULES_PATH環境変数が設定されていません")
	}
	return rulesPath, nil
}

// .mdファイルのリストを取得する
func GetMdFiles(rulesPath string, files []string) ([]string, error) {
	var mdFiles []string
	var missingFiles []string

	if len(files) == 0 {
		// 環境が指定されていない場合は、すべての.mdファイルを取得
		pattern := filepath.Join(rulesPath, "*.md")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("ファイルの検索に失敗しました: %w", err)
		}
		mdFiles = append(mdFiles, matches...)
	} else {
		// 環境が指定されている場合は、指定されたファイルのみを検索
		for _, file := range files {
			pattern := filepath.Join(rulesPath, file+".md")
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return nil, fmt.Errorf("ファイルの検索に失敗しました: %w", err)
			}
			if len(matches) > 0 {
				mdFiles = append(mdFiles, matches...)
			} else {
				missingFiles = append(missingFiles, file+".md")
			}
		}
	}

	if len(mdFiles) == 0 && len(missingFiles) == 0 {
		return nil, fmt.Errorf("指定されたファイルが見つかりませんでした (パス: %s)", rulesPath)
	}

	if len(missingFiles) > 0 {
		fmt.Printf("警告: ファイルが見つかりませんでした: %s\n", strings.Join(missingFiles, ", "))
	}

	return mdFiles, nil
}

// ファイルの内容を結合する
func CombineFiles(files []string) (string, error) {
	var content strings.Builder
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("ファイルの読み込みに失敗しました (%s): %w", file, err)
		}
		content.Write(data)
		content.WriteString("\n\n")
	}
	return content.String(), nil
}

// 出力ファイルのパスを取得する
func GetOutputPath(editor string) (string, error) {
	outputFileName := fmt.Sprintf(".%srules", editor)
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("現在のディレクトリの取得に失敗しました: %w", err)
	}
	return filepath.Join(currentDir, outputFileName), nil
}
