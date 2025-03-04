package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tKwbr999/rules_cli/internal/command"
	"github.com/tKwbr999/rules_cli/internal/handler"
)

// RULES_PATH環境変数を取得する
func getRulesPath() (string, error) {
	rulesPath := os.Getenv("RULES_PATH")
	if rulesPath == "" {
		return "", fmt.Errorf("RULES_PATH環境変数が設定されていません")
	}
	return rulesPath, nil
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
	fmt.Println("使用方法: rules [オプション] <エディタ> [ファイル...]")
	fmt.Println()
	fmt.Println("説明:")
	fmt.Println("  RULES_PATH環境変数で指定されたディレクトリから.mdファイルを読み込み、")
	fmt.Println("  指定されたエディタ用のrulesファイルを現在のディレクトリに作成します。")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -l, --list  利用可能な.mdファイルの一覧を表示します")
	fmt.Println()
	fmt.Println("引数:")
	fmt.Println("  <エディタ>  出力ファイル名の一部として使用される（例: cline → .clinerules）")
	fmt.Println("  [ファイル...]  結合する.mdファイルのファイル名を指定します (例: frontend backend)")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  rules cline           # すべての.mdファイルから.clinerulesを作成")
	fmt.Println("  rules cline frontend  # frontend.mdファイルからfrontend.mdを作成")
	fmt.Println("  rules cursor backend  # backend.mdファイルからbackend.mdを作成")
	fmt.Println("  rules cursor frontend backend # frontend.mdとbackend.mdを結合して.cursorrulesを作成")
	fmt.Println("  rules -l              # 利用可能な.mdファイルの一覧を表示")
}

// os.Exit関数をモック可能にするための変数
var osExit = os.Exit

func main() {
	// フラグを定義
	var listFiles bool
	flag.BoolVar(&listFiles, "l", false, "利用可能な.mdファイルの一覧を表示")

	// フラグを解析
	flag.Parse()

	// RULES_PATH環境変数の取得
	rulesPath, err := getRulesPath()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		printUsage()
		osExit(1)
	}

	// リスト表示フラグが設定されている場合
	if listFiles {
		fmt.Println("RULES_PATH:", rulesPath)
		mdFiles, err := handler.GetMdFiles(rulesPath, []string{})
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			osExit(1)
		}
		fmt.Println("利用可能な.mdファイル:")
		for _, file := range mdFiles {
			fmt.Println(filepath.Base(file))
		}
		return
	}

	// コマンドライン引数の解析
	editor, files := command.ParseArgs()

	// .mdファイルのリスト取得
	mdFiles, err := handler.GetMdFiles(rulesPath, files)
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
