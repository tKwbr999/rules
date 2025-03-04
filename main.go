package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tKwbr999/rules_cli/internal/command"
)


// RULES_PATH環境変数を取得する
func getRulesPath() (string, error) {
	rulesPath := os.Getenv("RULES_PATH")
	if rulesPath == "" {
		return "", fmt.Errorf("RULES_PATH環境変数が設定されていません")
	}
	return rulesPath, nil
}

// .mdファイルのリストを取得する
func getMdFiles(rulesPath string, files []string) ([]string, error) {
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
	fmt.Println("使用方法: rules [オプション] <エディタ> [環境]")
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
	fmt.Println("                ファイル名を指定しない場合は、すべての.mdファイルが使用されます")
	fmt.Println()
	fmt.Println("環境変数:")
	fmt.Println("  RULES_PATH  .mdファイルが格納されているディレクトリのパス")
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
		mdFiles, err := getMdFiles(rulesPath, []string{})
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
	editor, files, err := command.ParseArgs()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		printUsage()
		osExit(1)
	}

	// .mdファイルのリスト取得
	mdFiles, err := getMdFiles(rulesPath, files)
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
