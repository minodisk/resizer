package server

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/template"
	"github.com/minodisk/resizer/fetcher"
	"github.com/minodisk/resizer/input"
	"github.com/minodisk/resizer/options"
	"github.com/minodisk/resizer/processor"
	"github.com/minodisk/resizer/storage"
	"github.com/minodisk/resizer/uploader"
	"github.com/pkg/errors"
	"golang.org/x/net/netutil"
)

const (
	addr      = ":3000"
	errorHTML = `<!Doctype html>
<html>
<head>
  <title>{{ .StatusCode }} {{ .StatusText }}</title>
</head>
<body>
  <h1>{{ .StatusText }}</h1>
  <p>{{ .Message }}</p>
  <hr>
  <address>{{ .AppName }}</address>
</body>
</html>
`
)

var (
	contentTypes = map[string]string{
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
	}
	errorHTMLTemplate *template.Template
)

type ErrorHTML struct {
	StatusCode int
	StatusText string
	Message    string
	AppName    string
}

func NewErrorHTML(code int, message string) ErrorHTML {
	return ErrorHTML{
		StatusCode: code,
		StatusText: http.StatusText(code),
		Message:    message,
		AppName:    "Resizer",
	}
}

func init() {
	var err error
	errorHTMLTemplate, err = template.New("error").Parse(errorHTML)
	if err != nil {
		panic(err)
	}
}

func Start() error {
	options, err := options.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	handler, err := NewHandler(options)
	if err != nil {
		return err
	}
	server := http.Server{
		Handler:        &handler,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", options.Port))
	if err != nil {
		return err
	}
	if err := server.Serve(netutil.LimitListener(listener, options.MaxHTTPConnections)); err != nil {
		return errors.Wrap(err, "fail to serve")
	}
	return nil
}

type Handler struct {
	Storage  *storage.Storage
	Uploader *uploader.Uploader
	Hosts    []string
}

func NewHandler(o options.Options) (Handler, error) {
	s, err := storage.New(o)
	if err != nil {
		return Handler{}, err
	}
	u, err := uploader.New(o)
	if err != nil {
		return Handler{}, err
	}
	h := o.AllowedHosts
	return Handler{s, u, h}, nil
}

// ServeHTTP はリクエストに応じて処理を行いレスポンスする。
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprint(resp, "Not Found")
		log.Printf("'%s' not found\n", req.URL.Path)
		return
	}

	if err := h.operate(resp, req); err != nil {
		log.Println(errors.Wrap(err, "fail to operate"))
		resp.WriteHeader(http.StatusBadRequest)

		e := NewErrorHTML(http.StatusBadRequest, errors.Cause(err).Error())
		err := errorHTMLTemplate.Execute(resp, e)
		if err != nil {
			log.Println(errors.Wrap(err, "fail to generate error html from template"))
		}

		return
	}

	log.Println("OK")
}

// operate は手続き的に一連のリサイズ処理を行う。
// エラーを画一的に扱うためにメソッドとして切り分けを行っている
func (h *Handler) operate(resp http.ResponseWriter, req *http.Request) error {
	// 1. URLクエリからリクエストされているオプションを抽出する
	input, err := input.New(req.URL.Query())
	if err != nil {
		return err
	}
	input, err = input.Validate(h.Hosts)
	if err != nil {
		return err
	}
	i, err := storage.NewImage(input)
	if err != nil {
		return err
	}

	// 3. バリデート済みオプションでリサイズをしたキャッシュがあるか調べる
	// 4. キャッシュがあればリサイズ画像のURLにリダイレクトする
	cache := storage.Image{}
	h.Storage.Where(&storage.Image{
		ValidatedHash:    i.ValidatedHash,
		ValidatedWidth:   i.ValidatedWidth,
		ValidatedHeight:  i.ValidatedHeight,
		ValidatedMethod:  i.ValidatedMethod,
		ValidatedFormat:  i.ValidatedFormat,
		ValidatedQuality: i.ValidatedQuality,
	}).First(&cache)
	log.Printf("cache.ID=%d\n", cache.ID)
	if cache.ID != 0 {
		log.Printf("validated cache %+v exists, requested with %+v\n", cache, i)
		url := h.Uploader.CreateURL(cache.Filename)
		http.Redirect(resp, req, url, http.StatusSeeOther)
		return nil
	}
	log.Printf("validated cache doesn't exist, requested with %+v\n", i)

	// 5. 元画像を取得する
	// 6. リサイズの前処理をする
	filename, err := fetcher.Fetch(i.ValidatedURL)
	defer func() {
		if err := fetcher.Clean(filename); err != nil {
			log.Printf("fail to clean fetched file: %s\n", filename)
		}
	}()
	if err != nil {
		return err
	}
	var b []byte
	buf := bytes.NewBuffer(b)
	p := processor.New()
	pixels, err := p.Preprocess(filename)
	if err != nil {
		return err
	}

	// 7. 正規化する
	// 8. 正規化済みのオプションでリサイズをしたことがあるか調べる
	// 9. あればリサイズ画像のURLにリダイレクトする
	i, err = i.Normalize(pixels.Bounds().Size())
	if err != nil {
		return err
	}
	cache = storage.Image{}
	h.Storage.Where(&storage.Image{
		NormalizedHash:   i.NormalizedHash,
		DestWidth:        i.DestWidth,
		DestHeight:       i.DestHeight,
		ValidatedMethod:  i.ValidatedMethod,
		ValidatedFormat:  i.ValidatedFormat,
		ValidatedQuality: i.ValidatedQuality,
	}).First(&cache)
	if cache.ID != 0 {
		log.Printf("normalized cache %+v exists, requested with %+v\n", cache, i)
		url := h.Uploader.CreateURL(cache.Filename)
		http.Redirect(resp, req, url, http.StatusSeeOther)
		return nil
	}
	log.Printf("normalized cache doesn't exist, requested with %+v\n", i)

	// 10. リサイズする
	// 11. ファイルオブジェクトの処理結果フィールドを埋める
	// 12. レスポンスする
	size, err := p.Resize(pixels, buf, i)
	if err != nil {
		return err
	}
	b = buf.Bytes()

	i.ETag = fmt.Sprintf("%x", md5.Sum(b))
	i.Filename = i.CreateFilename()
	i.ContentType = contentTypes[i.ValidatedFormat]
	i.CanvasWidth = size.X
	i.CanvasHeight = size.Y

	resp.Header().Add("Content-Type", i.ContentType)
	io.Copy(resp, bufio.NewReader(buf))

	// レスポンスを完了させるために非同期に処理する
	go h.save(b, i)

	return nil
}

// save はファイルやデータを保存します。
func (h *Handler) save(b []byte, f storage.Image) {
	// 13. アップロードする
	// 14. キャッシュをDBに格納する
	if _, err := h.Uploader.Upload(bytes.NewBuffer(b), f); err != nil {
		log.Println(errors.Wrap(err, "fail to upload"))
		return
	}
	h.Storage.NewRecord(f)
	h.Storage.Create(&f)
	h.Storage.Save(&f)

	log.Println("complete to save")
}
