package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// コマンドライン引数を解析する
func parseArgs() (string, string, error) {
	if len(os.Args) < 2 {
		return "", "", fmt.Errorf("エディタの指定が必要です")
	}

	editor := os.Args[1]
	env := ""
	if len(os.Args) > 2 {
		env = os.Args[2]
	}

	return editor, env, nil
}

// RULES_PATH環境変数を取得する
func getRulesPath() (string, error) {
	rulesPath := os.Getenv("RULES_PATH")
	if rulesPath == "" {
		return "", fmt.Errorf("RULES_PATH環境変数が設定されていません")
	}
	return rulesPath, nil
}

// .mdファイルのリストを取得する
func getMdFiles(rulesPath, env string) ([]string, error) {
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
func combineFiles(files []string) (string, error) {
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
func getOutputPath(editor string) (string, error) {
	outputFileName := fmt.Sprintf(".%srules", editor)
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("現在のディレクトリの取得に失敗しました: %w", err)
	}
	return filepath.Join(currentDir, outputFileName), nil
}

// 使用方法を表示する
func printUsage() {
	fmt.Println("使用方法: rules <エディタ> [環境]")
	fmt.Println()
	fmt.Println("説明:")
	fmt.Println("  RULES_PATH環境変数で指定されたディレクトリから.mdファイルを読み込み、")
	fmt.Println("  指定されたエディタ用のrulesファイルを現在のディレクトリに作成します。")
	fmt.Println()
	fmt.Println("引数:")
	fmt.Println("  <エディタ>  出力ファイル名の一部として使用される（例: cline → .clinerules）")
	fmt.Println("  [環境]     特定の環境用のmdファイルのみを使用する場合に指定（例: frontend）")
	fmt.Println("             省略した場合はすべての.mdファイルが使用されます")
	fmt.Println()
	fmt.Println("環境変数:")
	fmt.Println("  RULES_PATH  .mdファイルが格納されているディレクトリのパス")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  rules cline           # すべての.mdファイルから.clinerules作成")
	fmt.Println("  rules cline frontend  # frontend.mdファイルから.clinerules作成")
	fmt.Println("  rules cursor backend  # backend.mdファイルから.cursorrules作成")
}

// os.Exit関数をモック可能にするための変数
var osExit = os.Exit

func main() {
	// コマンドライン引数の解析
	editor, env, err := parseArgs()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		printUsage()
		osExit(1)
	}

	// RULES_PATH環境変数の取得
	rulesPath, err := getRulesPath()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		printUsage()
		osExit(1)
	}

	// .mdファイルのリスト取得
	mdFiles, err := getMdFiles(rulesPath, env)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		osExit(1)
	}

	// ファイルの内容を結合
	content, err := combineFiles(mdFiles)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		osExit(1)
	}

	// 出力先のパスを取得
	outputPath, err := getOutputPath(editor)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		osExit(1)
	}

	// ファイルに書き込み
	err = os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("エラー: ファイルの書き込みに失敗しました (%s): %v\n", outputPath, err)
		osExit(1)
	}

	fmt.Printf("ファイルが作成されました: %s\n", outputPath)
	fmt.Printf("読み込んだファイル数: %d\n", len(mdFiles))
}
