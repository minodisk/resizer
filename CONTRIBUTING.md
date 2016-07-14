# Contributing

Go 1.6.2 で開発

## パッケージの依存管理

[Glide](https://github.com/Masterminds/glide)を使います。
下記に説明のない依存管理に関わる操作や設定は`glide help`かgideのドキュメントを参照して下さい。

### ビルド環境の構築

```bash
brew install glide
glide install
```

### Run

```bash
go run main.go --id=ID --bucket=BUCKET --json=JSON --dbuser=DBUSER --dbprotocol=DBPROTOCOL --dbaddress=DBADDRESS --dbname=DBNAME
```

ローカルで立ち上げてブラウザから試す。

```
http://localhost:3000/?width=300&url=http://example.com/foo.jpg
```

### Test

```bash
go test -v ./... -race
```

### Build

```bash
go build -v
```

## Run/Test/Build に必要なフラグオプションと環境変数

### 環境変数

#### `RESIZER_LOG_FILENAME`
ログを出力するファイル名です。空にしておくと標準出力にログを出力します。

### フラグオプション

#### `--id=ID`
GCPのプロジェクト名です。設定されていなければ環境変数`RESIZER_PROJECT_ID`を使用します。

#### `--bucket=BUCKET`
GCSのバケット名です。設定されていなければ環境変数`RESIZER_BUCKET`を使用します。

#### `--json=JSON`
GCPのサービスアカウントの秘密鍵であるJSONファイルへのパスです。設定されていなければ環境変数`RESIZER_JSON`を使用します。

#### `--dbuser=DBUSER`
DBのユーザーネームです。設定されていなければ環境変数`RESIZER_DB_USER`を使用します。

#### `--dbpassword=""`
DBのパスワードです。設定されていなければ環境変数`RESIZER_DB_PASSWORD`を使用します。デフォルトは空文字です。

#### `--dbprotocol=DBPROTOCOL`
DBへの接続プロトコルです。詳しくは[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#protocol)を参照してください。設定されていなければ環境変数`RESIZER_DB_PROTOCOL`を使用します。

#### `--dbaddress=DBADDRESS`
DBのアドレスです。詳しくは[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#address)を参照してください。設定されていなければ環境変数`RESIZER_DB_ADDRESS`を使用します。

#### `--dbname=DBNAME`
データの保存先となるデータベース名です。設定されていなければ環境変数`RESIZER_DB_NAME`を使用します。

#### `--host=HOST ...`
許可する画像URLのホストです。localhostはデフォルトで許可されています。複数のホストを設定する場合は、`--host=example.com --host=example2.com`のようにします。設定されていなければ環境変数`RESIZER_HOSTS`をカンマで区切ったものを使用します。例: `example.com,example2.com`

## 定数

### サーバーのポート

`3000`番

## デバッグ

### ベンチマーク

```bash
go test processor/processor_bench_test.go -bench . -benchmem
```

### プロファイル

[main.go](main.go)のpprof関連のコメントアウトを復帰しrunする。

[http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/)
