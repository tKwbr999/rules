package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// テスト用一時ディレクトリを作成
func createTestEnvironment(t *testing.T) (string, func()) {
	// 一時ディレクトリの作成
	tempDir, err := os.MkdirTemp("", "rules_test")
	if err != nil {
		t.Fatalf("テスト環境の作成に失敗しました: %v", err)
	}

	// テスト用のmarkdownファイル作成
	frontendContent := "# フロントエンドルール\n- React コンポーネントは関数コンポーネントを使用する\n- CSS は Tailwind を使用する"
	backendContent := "# バックエンドルール\n- REST API は RESTful 設計原則に従う\n- エラーハンドリングは統一された形式で行う"
	commonContent := "# 共通ルール\n- コードレビューは必須\n- テストカバレッジは 80% 以上とする"

	if err := os.WriteFile(filepath.Join(tempDir, "frontend.md"), []byte(frontendContent), 0644); err != nil {
		t.Fatalf("フロントエンドルールファイルの作成に失敗しました: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "backend.md"), []byte(backendContent), 0644); err != nil {
		t.Fatalf("バックエンドルールファイルの作成に失敗しました: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "common.md"), []byte(commonContent), 0644); err != nil {
		t.Fatalf("共通ルールファイルの作成に失敗しました: %v", err)
	}

	// クリーンアップ関数
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// 環境変数を管理
func manageEnvVariable(t *testing.T, key string, value string) func() {
	// 現在の環境変数を保存
	originalValue, exists := os.LookupEnv(key)

	// 新しい値を設定
	if value != "" {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("環境変数の設定に失敗しました: %v", err)
		}
	} else {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("環境変数の削除に失敗しました: %v", err)
		}
	}

	// 復元関数
	return func() {
		if exists {
			os.Setenv(key, originalValue)
		} else {
			os.Unsetenv(key)
		}
	}
}

// 出力ファイルを管理
func manageOutputFile(t *testing.T, filename string) func() {
	return func() {
		if _, err := os.Stat(filename); err == nil {
			if err := os.Remove(filename); err != nil {
				t.Logf("出力ファイルの削除に失敗しました: %v", err)
			}
		}
	}
}

// メイン関数のモックと実行
func runMain(t *testing.T, args ...string) {
	// 元のコマンドライン引数を保存
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// テスト用引数を設定
	os.Args = append([]string{"rules"}, args...)

	// メイン関数をキャプチャするためのリダイレクト
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// 終了コードをキャプチャ
	defer func() {
		// if r := recover(); r != nil {
		// 	if exitCode, ok := r.(int); ok && exitCode == 1 {
		// 		// 正常に終了コード1でexitしたとみなす
		// 		t.Errorf("プログラムが終了コード1で終了しました")
		// 	} else {
		// 		// 予期しないパニックの場合はテスト失敗
		// 		t.Fatalf("予期しないパニック: %v", r)
		// 	}
		// }
		// パイプを閉じる
		w.Close()
	}()

	// メイン関数実行
	// os.Exit()を呼び出すので、defer内でリカバリする
	func() {
		defer func() {
			if r := recover(); r != nil {
				panic(r) // パニックを上のdefer関数に伝搬
			}
		}()
		main()
	}()
}

// 全ファイル結合のテスト
func TestCombineAllFiles(t *testing.T) {
	// テスト環境の作成
	tempDir, cleanupDir := createTestEnvironment(t)
	defer cleanupDir()

	// 環境変数の設定
	restoreEnv := manageEnvVariable(t, "RULES_PATH", tempDir)
	defer restoreEnv()

	// 出力ファイルの管理
	cleanupOutput := manageOutputFile(t, ".clinerules")
	defer cleanupOutput()

	// テスト実行
	runMain(t, "cline")

	// 結果確認
	content, err := os.ReadFile(".clinerules")
	if err != nil {
		t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
	}

	// 全ファイルの内容が含まれているか確認
	if !containsString(string(content), "フロントエンドルール") ||
		!containsString(string(content), "バックエンドルール") ||
		!containsString(string(content), "共通ルール") {
		t.Errorf("結合ファイルに全てのコンテンツが含まれていません。内容: %s", string(content))
	}
}

// 特定の環境ファイル結合のテスト
func TestCombineSpecificFile(t *testing.T) {
	// テスト環境の作成
	tempDir, cleanupDir := createTestEnvironment(t)
	defer cleanupDir()

	// 環境変数の設定
	restoreEnv := manageEnvVariable(t, "RULES_PATH", tempDir)
	defer restoreEnv()

	// frontend.mdのテスト
	func() {
		// 出力ファイルの管理
		cleanupOutput := manageOutputFile(t, ".clinerules")
		defer cleanupOutput()

		// テスト実行
		runMain(t, "cline", "frontend")

		// 結果確認
		content, err := os.ReadFile(".clinerules")
		if err != nil {
			t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
		}

		// フロントエンドの内容だけが含まれているか確認
		if !containsString(string(content), "フロントエンドルール") {
			t.Errorf("結合ファイルにフロントエンドのコンテンツが含まれていません")
		}
		if containsString(string(content), "バックエンドルール") {
			t.Errorf("結合ファイルにバックエンドのコンテンツが含まれています")
		}
	}()

	// backend.mdのテスト
	func() {
		// 出力ファイルの管理
		cleanupOutput := manageOutputFile(t, ".cursorrules")
		defer cleanupOutput()

		// テスト実行
		runMain(t, "cursor", "backend")

		// 結果確認
		content, err := os.ReadFile(".cursorrules")
		if err != nil {
			t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
		}

		// バックエンドの内容だけが含まれているか確認
		if !containsString(string(content), "バックエンドルール") {
			t.Errorf("結合ファイルにバックエンドのコンテンツが含まれていません")
		}
		if containsString(string(content), "フロントエンドルール") {
			t.Errorf("結合ファイルにフロントエンドのコンテンツが含まれています")
		}
	}()
}

// 環境変数未設定のテスト
func TestNoRulesPathVariable(t *testing.T) {
	// 現在の環境変数を保存して削除
	restoreEnv := manageEnvVariable(t, "RULES_PATH", "")
	defer restoreEnv()

	// 終了コードのテスト方法を変更
	oldExit := osExit
	defer func() { osExit = oldExit }()

	var exitCode int
	var osExitCalled bool
	osExit = func(code int) {
		exitCode = code
		osExitCalled = true
		panic(code) // パニックを起こして実行を中断
	}

	defer func() {
		if r := recover(); r != nil {
			if code, ok := r.(int); ok && code == 1 {
				// 想定通りの終了
				if exitCode != 1 {
					t.Errorf("期待する終了コード1ではなく %d が返されました", exitCode)
				}
			} else {
				t.Fatalf("予期しないパニック: %v", r)
			}
		} else {
			if !osExitCalled {
				t.Error("環境変数未設定でもエラー終了しませんでした")
			}
		}
	}()

	// テスト実行
	runMain(t, "cline")
}

// 存在しないファイルパターンのテスト
func TestNonExistentFilePattern(t *testing.T) {
	// テスト環境の作成
	tempDir, cleanupDir := createTestEnvironment(t)
	defer cleanupDir()

	// 環境変数の設定
	restoreEnv := manageEnvVariable(t, "RULES_PATH", tempDir)
	defer restoreEnv()

	// 終了コードのテスト方法を変更
	oldExit := osExit
	defer func() { osExit = oldExit }()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
		panic(code) // パニックを起こして実行を中断
	}

	defer func() {
		if r := recover(); r != nil {
			if code, ok := r.(int); ok && code == 1 {
				// 想定通りの終了
				if exitCode != 1 {
					t.Errorf("期待する終了コード1ではなく %d が返されました", exitCode)
				}
			} else {
				t.Fatalf("予期しないパニック: %v", r)
			}
		}
	}()

	// 存在しない環境名でテスト実行
	runMain(t, "cline", "nonexistent")

	// ここに到達した場合は、期待通りにプログラムが終了しなかった
	t.Error("存在しないファイルパターンでもエラー終了しませんでした")
}

// 文字列が別の文字列に含まれるかチェック
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// この変数はmain.goで定義されているので、ここでは不要
// var osExit = os.Exit
