# rules

rulesファイルを作成するコマンドラインツール

## 概要

このツールは、RULES_PATH環境変数で指定されたディレクトリにある.mdファイルを読み込み、指定されたエディタ用のrulesファイルを生成します。

## インストール

```bash
go install ./...
```

## 使用方法

```
rules <エディタ> [環境]
```

### 引数

- `<エディタ>`: 出力ファイル名の一部として使用されます（例: cline → .clinerules）
- `[環境]`（任意）: 特定の環境用のmdファイルのみを使用する場合に指定します（例: frontend）
  - 省略した場合はすべての.mdファイルが使用されます

### 環境変数

- `RULES_PATH`: .mdファイルが格納されているディレクトリのパス

### 使用例

```bash
# すべての.mdファイルから.clinerules作成
rules cline

# frontend.mdファイルから.clinerules作成
rules cline frontend

# backend.mdファイルから.cursorrules作成
rules cursor backend

# frontend.mdとbackend.mdを結合しから.windsurfrules作成
rules windsurf backend
```

## 出力ファイル

コマンドが実行されたディレクトリに、以下の命名規則でファイルが作成されます:

```
.<エディタ>rules
```

例: `.clinerules`, `.cursorrules`
