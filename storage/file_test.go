package storage_test

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"testing"

	"github.com/go-microservices/resizer/option"
	"github.com/go-microservices/resizer/storage"
)

func TestNewImage(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	if _, err := storage.NewImage(map[string][]string{}, o.Hosts); err == nil {
		t.Fatalf("fail to NewFile: error=%v", err)
	}

	q := map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"method":  []string{storage.MethodThumbnail},
		"width":   []string{"400"},
		"height":  []string{"300"},
		"format":  []string{storage.FormatPng},
		"quality": []string{"80"},
	}
	f, err := storage.NewImage(q, o.Hosts)
	if err != nil {
		t.Fatalf("fail to New: error=%v", err)
	}
	if f.ValidatedURL != "http://0.0.0.0" {
		t.Errorf("URL is expected '%s', but actual '%s'", "http://0.0.0.0", f.ValidatedURL)
	}
	if f.ValidatedMethod != storage.MethodThumbnail {
		t.Errorf("Method is expected '%s', but actual '%s'", storage.MethodThumbnail, f.ValidatedMethod)
	}
	if f.ValidatedWidth != 400 {
		t.Errorf("RequestedWidth is expected %d, but actual %d", 400, f.ValidatedWidth)
	}
	if f.ValidatedHeight != 300 {
		t.Errorf("RequestedHeight is expected %d, but actual %d", 300, f.ValidatedHeight)
	}
	if f.ValidatedFormat != storage.FormatPng {
		t.Errorf("Format is expected '%s', but actual '%s'", storage.FormatPng, f.ValidatedFormat)
	}
	if f.ValidatedQuality != 0 {
		t.Errorf("Quality is expected %d, but actual %d", 0, f.ValidatedQuality)
	}
}

func TestValidateURL(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	// ローカルホストと特定のホストとそのサブドメインを許可する
	for _, host := range []string{
		"0.0.0.0",
		"255.255.255.255",
		"localhost",
	} {
		// http, https プロトコルを許可する
		for _, protocol := range []string{
			"http",
			"https",
		} {
			url := fmt.Sprintf("%s://%s", protocol, host)
			if _, err := storage.NewImage(map[string][]string{
				"url":   []string{url},
				"width": []string{"400"},
			}, o.Hosts); err != nil {
				t.Errorf("should allow %s error=%v", url, err)
			}
			url = fmt.Sprintf("%s:%d", url, 99999)
			if _, err := storage.NewImage(map[string][]string{
				"url":   []string{url},
				"width": []string{"400"},
			}, o.Hosts); err != nil {
				t.Errorf("should allow %s error=%v", url, err)
			}
		}
		// その他のプロトコルを許可しない
		for _, protocol := range []string{
			"ftp",
			"ssh",
			"git",
		} {
			url := fmt.Sprintf("%s://%s", protocol, host)
			if _, err := storage.NewImage(map[string][]string{
				"url":   []string{url},
				"width": []string{"400"},
			}, o.Hosts); err == nil {
				t.Errorf("shouldn't allow %s", url)
			}
			url = fmt.Sprintf("%s:%d", url, 99999)
			if _, err := storage.NewImage(map[string][]string{
				"url":   []string{},
				"width": []string{"400"},
			}, o.Hosts); err == nil {
				t.Errorf("shouldn't allow %s", url)
			}
		}
	}

	// その他のホストとそのサブドメインを許可しない
	for _, host := range []string{
		"example.com",
		"foo.bar.baz.example.com",
	} {
		// どのようなプロトコルも許可しない
		for _, protocol := range []string{
			"http",
			"https",
			"ftp",
			"ssh",
			"git",
		} {
			url := fmt.Sprintf("%s://%s", protocol, host)
			if _, err := storage.NewImage(map[string][]string{
				"url":   []string{url},
				"width": []string{"400"},
			}, o.Hosts); err == nil {
				t.Errorf("shouldn't allow %s", url)
			}
			url = fmt.Sprintf("%s:%d", url, 99999)
			if _, err := storage.NewImage(map[string][]string{
				"url":   []string{url},
				"width": []string{"400"},
			}, o.Hosts); err == nil {
				t.Errorf("shouldn't allow %s", url)
			}
		}
	}
}

func TestValidateMethod(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
	}, o.Hosts); err != nil {
		t.Fatalf("fail to validate: error=%v", err)
	} else if f.ValidatedMethod != "normal" {
		t.Errorf("default method should be normal")
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"method": []string{"normal"},
	}, o.Hosts); err != nil {
		t.Fatalf("fail to validate: error=%v", err)
	} else if f.ValidatedMethod != "normal" {
		t.Errorf("format should be normal")
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"method": []string{"thumbnail"},
	}, o.Hosts); err != nil {
		t.Errorf("thumbnail format is allowed")
	} else if f.ValidatedMethod != "thumbnail" {
		t.Errorf("format should be thumbnail")
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"method": []string{"foo"},
	}, o.Hosts); err == nil {
		t.Errorf("format other than normal or thumbnail isn't allowed")
	}
}

func TestValidateSize(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"-100"},
		"height": []string{"100"},
	}, o.Hosts); err == nil {
		t.Errorf("negative width isn't allowed")
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"100"},
		"height": []string{"-100"},
	}, o.Hosts); err == nil {
		t.Errorf("negative height isn't allowed")
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"-100"},
		"height": []string{"-100"},
	}, o.Hosts); err == nil {
		t.Errorf("negative width and height isn't allowed")
	}

	func() {
		if _, err := storage.NewImage(map[string][]string{
			"url":    []string{"http://0.0.0.0"},
			"width":  []string{"0"},
			"height": []string{"0"},
		}, o.Hosts); err == nil {
			t.Errorf("zero size isn't allowed")
		}
	}()

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"0"},
		"height": []string{"100"},
	}, o.Hosts); err != nil {
		t.Errorf("zero width is allowed err=", err)
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"100"},
		"height": []string{"0"},
	}, o.Hosts); err != nil {
		t.Errorf("zero height is allowed err=", err)
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"100"},
		"height": []string{"100"},
	}, o.Hosts); err != nil {
		t.Errorf("non-zero size is allowed err=", err)
	}
}

func TestValidateFormat(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
	}, o.Hosts); err != nil {
		t.Fatalf("fail to validate: error=%v", err)
	} else if f.ValidatedFormat != "jpeg" {
		t.Errorf("default format should be jpeg")
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"format": []string{"jpeg"},
	}, o.Hosts); err != nil {
		t.Errorf("jpeg format is allowed")
	} else if f.ValidatedFormat != "jpeg" {
		t.Errorf("format should be jpeg")
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"format": []string{"png"},
	}, o.Hosts); err != nil {
		t.Errorf("png format is allowed")
	} else if f.ValidatedFormat != "png" {
		t.Errorf("format should be png")
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"format": []string{"gif"},
	}, o.Hosts); err != nil {
		t.Errorf("gif format is allowed")
	} else if f.ValidatedFormat != "gif" {
		t.Errorf("format should be gif")
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
		"format": []string{"foo"},
	}, o.Hosts); err == nil {
		t.Errorf("format is allowed only jpeg or png or gif")
	}
}

func TestValidateQuality(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	if f, err := storage.NewImage(map[string][]string{
		"url":    []string{"http://0.0.0.0"},
		"width":  []string{"400"},
		"height": []string{"0"},
	}, o.Hosts); err != nil {
		t.Errorf("fail to validate: error=%v", err)
	} else if f.ValidatedQuality != 100 {
		t.Errorf("default quality should be 100")
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"400"},
		"height":  []string{"0"},
		"quality": []string{"-1"},
	}, o.Hosts); err == nil {
		t.Errorf("negative quality shouldn't be allowed")
	}

	if _, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"400"},
		"height":  []string{"0"},
		"quality": []string{"101"},
	}, o.Hosts); err == nil {
		t.Errorf("over 100 quality shouldn't be allowed")
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"400"},
		"height":  []string{"0"},
		"quality": []string{"67"},
		"format":  []string{storage.FormatJpeg},
	}, o.Hosts); err != nil {
		t.Errorf("fail to validate: error=%v", err)
	} else if f.ValidatedQuality != 67 {
		t.Errorf("Quality with png format should be 67, but actual %d", f.ValidatedQuality)
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"400"},
		"height":  []string{"0"},
		"quality": []string{"67"},
		"format":  []string{storage.FormatPng},
	}, o.Hosts); err != nil {
		t.Errorf("fail to validate: error=%v", err)
	} else if f.ValidatedQuality != 0 {
		t.Errorf("Quality with png format should be 0, but actual %d", f.ValidatedQuality)
	}

	if f, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"400"},
		"height":  []string{"0"},
		"quality": []string{"67"},
		"format":  []string{storage.FormatGif},
	}, o.Hosts); err != nil {
		t.Errorf("fail to validate: error=%v", err)
	} else if f.ValidatedQuality != 0 {
		t.Errorf("Quality with gif format should be 0, but actual %d", f.ValidatedQuality)
	}
}

func TestSerialization(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	f1, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"400"},
		"height":  []string{"300"},
		"method":  []string{"thumbnail"},
		"format":  []string{"png"},
		"quality": []string{"30"},
	}, o.Hosts)
	if err != nil {
		t.Fatalf("should be valid image: error=%v", err)
	}

	f2, err := storage.NewImage(map[string][]string{
		"quality":  []string{"80"},
		"format":   []string{"png"},
		"method":   []string{"thumbnail"},
		"height":   []string{"300"},
		"width":    []string{"400"},
		"url":      []string{"http://0.0.0.0"},
		"filename": []string{"dummy"},
	}, o.Hosts)
	if err != nil {
		t.Fatalf("should be valid image: error=%v", err)
	}

	f3, err := storage.NewImage(map[string][]string{
		"url":     []string{"http://0.0.0.0"},
		"width":   []string{"200"},
		"height":  []string{"500"},
		"method":  []string{"thumbnail"},
		"format":  []string{"png"},
		"quality": []string{"80"},
	}, o.Hosts)
	if err != nil {
		t.Fatalf("should be valid image: error=%v", err)
	}

	if f1.ValidatedHash != f2.ValidatedHash {
		t.Errorf("RequestedHash with equal request should be equal: %+v %+v", f1, f2)
	}
	if f1.ValidatedHash == f3.ValidatedHash {
		t.Errorf("RequestedHash with different file shouldn't be equal: %+v %+v", f1, f3)
	}
	if f2.ValidatedHash == f3.ValidatedHash {
		t.Errorf("RequestedHash with different file shouldn't be equal: %+v %+v", f2, f3)
	}
}

func TestGuessSizeError(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	func() {
		f, err := storage.NewImage(map[string][]string{
			"url":    []string{"http://0.0.0.0"},
			"width":  []string{"200"},
			"height": []string{"100"},
		}, o.Hosts)
		if err != nil {
			t.Fatalf("should be valid image: error=%v", err)
		}
		if _, err := f.Normalize(image.Point{0, 300}); err == nil {
			t.Errorf("with zero source width should return error")
		}
	}()

	func() {
		f, err := storage.NewImage(map[string][]string{
			"url":    []string{"http://0.0.0.0"},
			"width":  []string{"200"},
			"height": []string{"100"},
		}, o.Hosts)
		if err != nil {
			t.Fatalf("should be valid image: error=%v", err)
		}
		if _, err := f.Normalize(image.Point{400, 0}); err == nil {
			t.Errorf("with zero source height should return error")
		}
	}()

	func() {
		f, err := storage.NewImage(map[string][]string{
			"url":    []string{"http://0.0.0.0"},
			"width":  []string{"200"},
			"height": []string{"100"},
		}, o.Hosts)
		if err != nil {
			t.Fatalf("should be valid image: error=%v", err)
		}
		if _, err := f.Normalize(image.ZP); err == nil {
			t.Errorf("source size with zero point should return error")
		}
	}()
}

func TestNormalizeWithMethodNormal(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	for _, points := range []map[string]image.Point{
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 0},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{0, 150},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{800, 0},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{0, 600},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 150},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{300, 150},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{500, 150},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 200},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 400},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{500, 450},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{600, 450},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{700, 450},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{800, 400},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{200, 100},
			"target": image.Point{400, 300},
			"dest":   image.Point{200, 100},
			"canvas": image.Point{200, 100},
		},
	} {
		source := points["source"]
		target := points["target"]
		dest := points["dest"]
		canvas := points["canvas"]
		f, err := storage.NewImage(map[string][]string{
			"url":    []string{"http://0.0.0.0"},
			"method": []string{storage.MethodNormal},
			"width":  []string{strconv.Itoa(target.X)},
			"height": []string{strconv.Itoa(target.Y)},
		}, o.Hosts)
		if err != nil {
			t.Fatalf("should be valid image: error=%v", err)
		}
		if f, err := f.Normalize(source); err != nil {
			t.Errorf("Normalize with normal method should not return error: %v", err)
		} else if f.DestWidth != dest.X {
			t.Errorf("DestWidth expected %d, but actual %d", dest.X, f.DestWidth)
		} else if f.DestHeight != dest.Y {
			t.Errorf("DestHeight expected %d, but actual %d", dest.Y, f.DestHeight)
		} else if f.CanvasWidth != canvas.X {
			t.Errorf("CanvasWidth expected %d, but actual %d", canvas.X, f.CanvasWidth)
		} else if f.CanvasHeight != canvas.Y {
			t.Errorf("CanvasHeight expected %d, but actual %d", canvas.Y, f.CanvasHeight)
		}
	}
}

func TestNormalizeWithMethodThumbnail(t *testing.T) {
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	for i, points := range []map[string]image.Point{
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 0},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{0, 150},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 150},
			"dest":   image.Point{200, 150},
			"canvas": image.Point{200, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{300, 150},
			"dest":   image.Point{300, 225},
			"canvas": image.Point{300, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{500, 150},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 150},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 200},
			"dest":   image.Point{266, 200},
			"canvas": image.Point{200, 200},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{200, 400},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{200, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{500, 450},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{600, 450},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
		map[string]image.Point{
			"source": image.Point{400, 300},
			"target": image.Point{700, 450},
			"dest":   image.Point{400, 300},
			"canvas": image.Point{400, 300},
		},
	} {
		source := points["source"]
		target := points["target"]
		dest := points["dest"]
		canvas := points["canvas"]
		f, err := storage.NewImage(map[string][]string{
			"url":    []string{"http://0.0.0.0"},
			"method": []string{"thumbnail"},
			"width":  []string{strconv.Itoa(target.X)},
			"height": []string{strconv.Itoa(target.Y)},
		}, o.Hosts)
		if err != nil {
			t.Fatalf("should be valid image: error=%v", err)
		}
		if f, err := f.Normalize(source); err != nil {
			t.Errorf("Normalize with thumbnail method should not return error: %v", err)
		} else if f.DestWidth != dest.X {
			t.Errorf("DestWidth(%d) expected %d, but actual %d", i, dest.X, f.DestWidth)
		} else if f.DestHeight != dest.Y {
			t.Errorf("DestHeight(%d) expected %d, but actual %d", i, dest.Y, f.DestHeight)
		} else if f.CanvasWidth != canvas.X {
			t.Errorf("CanvasWidth(%d) expected %d, but actual %d", i, canvas.X, f.CanvasWidth)
		} else if f.CanvasHeight != canvas.Y {
			t.Errorf("CanvasHeight(%d) expected %d, but actual %d", i, canvas.Y, f.CanvasHeight)
		}
	}
}
