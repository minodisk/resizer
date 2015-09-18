# resizer

## Keywords

- リサイズ: 回転やスケールなど、このアプリケーションで行う一連の画像処理全体のこと。
- スケール: 拡大縮小処理。
- アスペクト比: `width/height`。

## 仕様

- 目的のサイズと同じか小さなサイズにリサイズする。
  - リサイズ後の画像が目的のサイズと同サイズである保証はない。
- 元画像を縦横同じ比率で1倍以下にスケールする。
  - リサイズ後の画像が一方向に潰れることはない。
  - 拡大することはない。
- スケール後に再エンコードを行うので元画像のメタ情報は引き継がれない。
  - 元画像がjpegで且つEXIFにOrientationタグが存在する場合、画像情報に回転を反映してからスケール処理を行う。

## API

### エンドポイント

```http:Endpoint
GET http://your.host.name/
```

### リクエスト

- 値に`&`等の記号が入っている場合はURLエンコードする必要がある。
  - 予め値が判明していて`&`等の記号が含まれていない場合はURLエンコードする必要がない。
  - 値がユーザー入力で決定する等、予測不可能な場合は必ずURLエンコードするべき。
- オプショナルなパラメータは設定されていなければデフォルト値を自動的に設定する。
- パラメータに値が設定されていて、且つ適切でない値が設定された場合はエラーとなる。

#### url

画像のURL。必須。

- 特定のホストしているサービスのホスト名でなければならない。

#### width, height

幅と高さ。単位px。`0`〜。オプショナル(最低でも width, height のどちらか1つは指定する必要がある)。デフォルト`0`。

- 「`0`以上の整数値」以外の値は許可しない。
- 共に`0`は許可しない。
- *width*だけが`0`なら*height*と元画像のアスペクト比から*width*を算出して設定する。
- *height*だけが`0`なら*width*と元画像のアスペクト比から*height*を算出して設定する。
- 元画像のサイズより大きければ無視し、元画像のサイズを採用する。

#### method

リサイズ方法。`normal`|`thumbnail`。オプショナル。デフォルト`normal`。

- *width*または*height*のどちらかが`0`なら無視する。
- 選択肢以外の値は許可しない。
- `normal`の場合、必ず目的のサイズに収まるように画像をスケールする。画像の一部を切り取ることはない。
- `thumbnail`の場合、なるべく目的のサイズの全ピクセルを塗るように画像をスケールする。目的サイズの外側を切り取り、内側を結果とする。

#### format

フォーマット。`jpeg`|`png`|`gif`。オプショナル。デフォルト`jpeg`。

- 選択肢以外の値は許可しない。

#### quality

jpegのクオリティ。`1`〜`100`。オプショナル。デフォルト`100`。

- *format*が`jpeg`以外なら無視する。
- 「`1`〜`100`の整数値」以外の値は許可しない。

### レスポンス

#### エラー時

- `4xx`系レスポンスの場合はレスポンスボディに理由の詳細を記載する。
  - `message`: エラーメッセージです。

#### 正常時

- 同条件のパラメーターでリサイズを行ったことがない場合、リサイズ済の画像データをレスポンスする。
- 同条件のパラメーターでリサイズを行ったことがある場合、S3の画像URLにリダイレクトする。

### 例

リクエスト

```http:HTTPRequest
GET http://your.host.name/?url=http%3A%2F%2Fexample.com%2Fimage.jpeg&width=800 HTTP/1.1
Host: your.host.name
```

---

同条件のパラメーターでリサイズを行ったことがない場合のレスポンス

```http:HTTPResponse
HTTP/1.1 200 OK
Content-Type: image/jpeg
```

---

同条件のパラメーターでリサイズを行ったことがある場合のレスポンス

```http:HTTPResponse
HTTP/1.1 303 See Other
Content-Type: text/html; charset=utf-8
Location: https://s3-ap-northeast-1.amazonaws.com/your-bucket-name/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-xxx.jpeg
```

---

エラーした場合

```http:HTTPResponse
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"message":"detail of the error"}
```

## 開発方法

[CONTRIBUTING.md](CONTRIBUTING.md)を参照。
