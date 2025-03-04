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
func GetMdFiles(rulesPath, env string) ([]string, error) {
	filePattern := "*.md"
	if env != "" {
		filePattern = env + ".md"
	}

	mdFiles, err := filepath.Glob(filepath.Join(rulesPath, filePattern))
	if err != nil {
		return nil, fmt.Errorf("ファイルの検索に失敗しました: %w", err)
	}

	if len(mdFiles) == 0 {
		return nil, fmt.Errorf("%sに該当する.mdファイルが見つかりませんでした (パス: %s)", filePattern, rulesPath)
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
