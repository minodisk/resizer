# Contributing

Go 1.4.2 で開発

## パッケージの依存管理

[nitrous-io/goop](https://github.com/nitrous-io/goop)を使います。
下記に説明のない依存管理に関わる操作や設定は`goop help`かgoopのドキュメントを参照して下さい。

### ビルド環境の構築

```bash
go get github.com/nitrous-io/goop
goop install
```

### Run

```bash
goop go run main.go
```

ローカルで立ち上げてブラウザから試す。

```
http://localhost:3000/?width=300&url=http://example.com/foo.jpg
```

### Test

```bash
goop go test -v ./... -race
```

### Build

```bash
goop go build -v
```

## Run/Test/Build に必要な環境変数

### アプリケーション

- `RESIZER_LOG_FILENAME`: ログを出力するファイル名です。空にしておくと標準出力にログを出力します。

### AWS関連(S3)

- `AWS_ACCESS_KEY_ID`: AWSのアクセスキーです。
- `AWS_SECRET_ACCESS_KEY`: AWSのアクセスシークレットです。
- `RESIZER_S3_REGION`: S3のAZ (`ap-northeast-1` など) です。
- `RESIZER_S3_BUCKET`: S3のバケット名 (`your-bucket-name` など) です。

### DB関連(RDS)

- `RESIZER_DB_USERNAME`: DBのユーザーネームです。
- `RESIZER_DB_PASSWORD`: DBのパスワードです。
- `RESIZER_DB_ENDPOINT`: DBのエンドポイントです。(`/resizer` など)

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
