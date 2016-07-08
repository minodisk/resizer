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
go run main.go
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

## Run/Test/Build に必要な環境変数

### アプリケーション

- `RESIZER_LOG_FILENAME`: ログを出力するファイル名です。空にしておくと標準出力にログを出力します。

### GCP関連(GCS)

- `RESIZER_GCS_SERVICE_ACCOUNT`: GCPのサービスアカウントの秘密鍵であるJSONファイルへのパスです。
- `RESIZER_GCS_PROJECT_ID`: GCPのプロジェクト名です。
- `RESIZER_S3_BUCKET`: GCSのバケット名です。

### DB関連(Cloud SQL)

- `RESIZER_DB_USERNAME`: DBのユーザーネームです。
- `RESIZER_DB_PASSWORD`: DBのパスワードです。
- `RESIZER_DB_PROTOCOL`: DBへの接続プロトコルです。詳しくは[こちら](https://github.com/go-sql-driver/mysql#protocol)を参照してください。
- `RESIZER_DB_ADDRESS`: DBのアドレスです。詳しくは[こちら](https://github.com/go-sql-driver/mysql#address)を参照してください。
- `RESIZER_DB_NAME`: データベース名です。

## 定数

### サーバーのポート

`3000`番

## デバッグ

### ベンチマーク

```bash
goop go test processor/processor_bench_test.go -bench . -benchmem
```

### プロファイル

[main.go](main.go)のpprof関連のコメントアウトを復帰しrunする。

[http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/)
