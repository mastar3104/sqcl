# sqcl

Go製のターミナルSQLクライアント

## Features

- **MySQL サポート** - MySQL データベースへの接続と操作
- **接続設定の保存** - よく使う接続を名前をつけて保存・管理
- **TAB キーによる自動補完** - SQLキーワード、テーブル名、カラム名を補完
- **マルチライン入力** - `;` で終端するまで複数行にわたってSQLを入力可能
- **コマンド履歴の永続化** - 履歴を保存し、次回起動時に復元
- **テーブル形式の出力** - クエリ結果を見やすいテーブル形式で表示
- **内部コマンド** - データベースのメタ情報を簡単に取得
- **プレースホルダー対応** - `?` を使ったパラメータ化クエリをインタラクティブに実行

## Installation

```bash
go install github.com/mastar3104/sqcl/cmd/sqcl@latest
```

## Usage

```bash
# 直接接続
sqcl -dsn 'user:pass@tcp(host:port)/dbname'

# 保存済み接続を使用
sqcl -c mydb
```

### 接続の保存・管理

```bash
# 接続を保存
sqcl save mydb -dsn 'user:pass@tcp(host:port)/dbname'

# 保存済み接続の一覧
sqcl list

# 保存済み接続を削除
sqcl remove mydb
```

### コマンドラインオプション

| フラグ | 説明 | デフォルト |
|--------|------|-----------|
| `-dsn` | 接続文字列 | - |
| `-c` | 保存済み接続名を指定 | - |
| `-driver` | データベースドライバ | `mysql` |
| `-history` | 履歴ファイルパス | `~/.sqlc_history` |
| `-cache-ttl` | メタデータキャッシュTTL | `60s` |
| `-version` | バージョン表示 | - |

※ `-dsn` または `-c` のいずれかが必須

## 内部コマンド

| コマンド | エイリアス | 説明 |
|----------|-----------|------|
| `:help` | `:h`, `:?` | ヘルプ表示 |
| `:quit` | `:q`, `:exit` | 終了 |
| `:tables` | - | テーブル一覧 |
| `:columns <table>` | `:cols <table>` | カラム一覧 |
| `:databases` | `:dbs` | データベース一覧 |
| `:reload` | `:refresh` | メタデータキャッシュ再読み込み |
| `:status` | - | 接続状態表示 |

## キーバインド

| キー | 説明 |
|------|------|
| `TAB` | 自動補完（次の候補） |
| `Shift+TAB` | 自動補完（前の候補） |
| `Ctrl+C` | 入力キャンセル |
| `Ctrl+D` | 終了 |

## プレースホルダー

SQL文中の `?` はプレースホルダーとして認識され、実行前に値の入力を求められます。

```
sqcl(mydb)> SELECT * FROM users WHERE status = ? AND age > ?;

Query: SELECT * FROM users WHERE status = ? AND age > ?
Enter values for 2 placeholder(s):
  (Press Enter for NULL, Ctrl+C to cancel)

  [1]> active
  [2]> 18

+----+-------+--------+-----+
| id | name  | status | age |
+----+-------+--------+-----+
| 1  | Alice | active | 25  |
+----+-------+--------+-----+
1 row(s) in set
Time: 12.345ms
```

### 型推論

| 入力 | 変換結果 |
|-----|---------|
| (空入力) | `NULL` |
| `NULL` | `NULL` |
| `123` | 整数 |
| `3.14` | 浮動小数点 |
| `hello` | 文字列 |

### 注意事項

- 文字列リテラル内の `?` (`'text?'`) はプレースホルダーとして認識されません
- バッククォート内 (`` `col?` ``) やダブルクォート内 (`"col?"`) も同様です

## プロジェクト構造

```
.
├── cmd/
│   └── sqcl/
│       └── main.go          # エントリポイント
└── internal/
    ├── app/                  # アプリケーション設定・起動
    ├── cache/                # メタデータキャッシュ
    ├── completion/           # 自動補完ロジック
    ├── connections/          # 接続設定の保存・管理
    ├── db/                   # データベース抽象化層
    │   └── mysql/            # MySQL実装
    ├── history/              # 履歴管理
    ├── placeholder/          # プレースホルダー検出・入力処理
    ├── render/               # 出力フォーマッタ
    └── repl/                 # REPL・コマンド処理
```

## 依存関係

- [github.com/chzyer/readline](https://github.com/chzyer/readline) - 行編集・履歴・補完
- [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL ドライバ

## License

MIT License
