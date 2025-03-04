#!/bin/bash

# テスト用のディレクトリ作成
mkdir -p test_rules

# テスト用のmarkdownファイル作成
echo "# フロントエンドルール" > test_rules/frontend.md
echo "- React コンポーネントは関数コンポーネントを使用する" >> test_rules/frontend.md
echo "- CSS は Tailwind を使用する" >> test_rules/frontend.md

echo "# バックエンドルール" > test_rules/backend.md
echo "- REST API は RESTful 設計原則に従う" >> test_rules/backend.md
echo "- エラーハンドリングは統一された形式で行う" >> test_rules/backend.md

echo "# 共通ルール" > test_rules/common.md
echo "- コードレビューは必須" >> test_rules/common.md
echo "- テストカバレッジは 80% 以上とする" >> test_rules/common.md

# ビルド
go build -o rules

# RULES_PATH 環境変数を設定
export RULES_PATH="./test_rules"

# テスト実行
echo "===== すべてのファイルから .clinerules 作成 ====="
./rules cline
cat .clinerules
echo ""

echo "===== frontend.md から .clinerules 作成 ====="
./rules cline frontend
cat .clinerules
echo ""

echo "===== backend.md から .cursorrules 作成 ====="
./rules cursor backend
cat .cursorrules
echo ""

# テストファイル削除
echo "===== テスト終了 ====="
rm -f .clinerules .cursorrules
rm -rf test_rules
