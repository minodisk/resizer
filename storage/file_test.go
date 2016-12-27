package storage_test

// func TestSerialization(t *testing.T) {
// 	o, err := option.New(os.Args[1:])
// 	if err != nil {
// 		t.Fatalf("fail to create options: error=%v", err)
// 	}
// 	f1, err := storage.NewImage(map[string][]string{
// 		"url":     []string{"http://0.0.0.0"},
// 		"width":   []string{"400"},
// 		"height":  []string{"300"},
// 		"method":  []string{"thumbnail"},
// 		"format":  []string{"png"},
// 		"quality": []string{"30"},
// 	}, o.Hosts)
// 	if err != nil {
// 		t.Fatalf("should be valid image: error=%v", err)
// 	}
//
// 	f2, err := storage.NewImage(map[string][]string{
// 		"quality":  []string{"80"},
// 		"format":   []string{"png"},
// 		"method":   []string{"thumbnail"},
// 		"height":   []string{"300"},
// 		"width":    []string{"400"},
// 		"url":      []string{"http://0.0.0.0"},
// 		"filename": []string{"dummy"},
// 	}, o.Hosts)
// 	if err != nil {
// 		t.Fatalf("should be valid image: error=%v", err)
// 	}
//
// 	f3, err := storage.NewImage(map[string][]string{
// 		"url":     []string{"http://0.0.0.0"},
// 		"width":   []string{"200"},
// 		"height":  []string{"500"},
// 		"method":  []string{"thumbnail"},
// 		"format":  []string{"png"},
// 		"quality": []string{"80"},
// 	}, o.Hosts)
// 	if err != nil {
// 		t.Fatalf("should be valid image: error=%v", err)
// 	}
//
// 	if f1.ValidatedHash != f2.ValidatedHash {
// 		t.Errorf("RequestedHash with equal request should be equal: %+v %+v", f1, f2)
// 	}
// 	if f1.ValidatedHash == f3.ValidatedHash {
// 		t.Errorf("RequestedHash with different file shouldn't be equal: %+v %+v", f1, f3)
// 	}
// 	if f2.ValidatedHash == f3.ValidatedHash {
// 		t.Errorf("RequestedHash with different file shouldn't be equal: %+v %+v", f2, f3)
// 	}
// }
//
// func TestGuessSizeError(t *testing.T) {
// 	o, err := option.New(os.Args[1:])
// 	if err != nil {
// 		t.Fatalf("fail to create options: error=%v", err)
// 	}
// 	func() {
// 		f, err := storage.NewImage(map[string][]string{
// 			"url":    []string{"http://0.0.0.0"},
// 			"width":  []string{"200"},
// 			"height": []string{"100"},
// 		}, o.Hosts)
// 		if err != nil {
// 			t.Fatalf("should be valid image: error=%v", err)
// 		}
// 		if _, err := f.Normalize(image.Point{0, 300}); err == nil {
// 			t.Errorf("with zero source width should return error")
// 		}
// 	}()
//
// 	func() {
// 		f, err := storage.NewImage(map[string][]string{
// 			"url":    []string{"http://0.0.0.0"},
// 			"width":  []string{"200"},
// 			"height": []string{"100"},
// 		}, o.Hosts)
// 		if err != nil {
// 			t.Fatalf("should be valid image: error=%v", err)
// 		}
// 		if _, err := f.Normalize(image.Point{400, 0}); err == nil {
// 			t.Errorf("with zero source height should return error")
// 		}
// 	}()
//
// 	func() {
// 		f, err := storage.NewImage(map[string][]string{
// 			"url":    []string{"http://0.0.0.0"},
// 			"width":  []string{"200"},
// 			"height": []string{"100"},
// 		}, o.Hosts)
// 		if err != nil {
// 			t.Fatalf("should be valid image: error=%v", err)
// 		}
// 		if _, err := f.Normalize(image.ZP); err == nil {
// 			t.Errorf("source size with zero point should return error")
// 		}
// 	}()
// }
//
// func TestNormalizeWithMethodNormal(t *testing.T) {
// 	o, err := option.New(os.Args[1:])
// 	if err != nil {
// 		t.Fatalf("fail to create options: error=%v", err)
// 	}
// 	for _, points := range []map[string]image.Point{
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 0},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{0, 150},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{800, 0},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{0, 600},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 150},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{300, 150},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{500, 150},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 200},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 400},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{500, 450},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{600, 450},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{700, 450},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{800, 400},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{200, 100},
// 			"target": image.Point{400, 300},
// 			"dest":   image.Point{200, 100},
// 			"canvas": image.Point{200, 100},
// 		},
// 	} {
// 		source := points["source"]
// 		target := points["target"]
// 		dest := points["dest"]
// 		canvas := points["canvas"]
// 		f, err := storage.NewImage(map[string][]string{
// 			"url":    []string{"http://0.0.0.0"},
// 			"method": []string{storage.MethodNormal},
// 			"width":  []string{strconv.Itoa(target.X)},
// 			"height": []string{strconv.Itoa(target.Y)},
// 		}, o.Hosts)
// 		if err != nil {
// 			t.Fatalf("should be valid image: error=%v", err)
// 		}
// 		if f, err := f.Normalize(source); err != nil {
// 			t.Errorf("Normalize with normal method should not return error: %v", err)
// 		} else if f.DestWidth != dest.X {
// 			t.Errorf("DestWidth expected %d, but actual %d", dest.X, f.DestWidth)
// 		} else if f.DestHeight != dest.Y {
// 			t.Errorf("DestHeight expected %d, but actual %d", dest.Y, f.DestHeight)
// 		} else if f.CanvasWidth != canvas.X {
// 			t.Errorf("CanvasWidth expected %d, but actual %d", canvas.X, f.CanvasWidth)
// 		} else if f.CanvasHeight != canvas.Y {
// 			t.Errorf("CanvasHeight expected %d, but actual %d", canvas.Y, f.CanvasHeight)
// 		}
// 	}
// }
//
// func TestNormalizeWithMethodThumbnail(t *testing.T) {
// 	o, err := option.New(os.Args[1:])
// 	if err != nil {
// 		t.Fatalf("fail to create options: error=%v", err)
// 	}
// 	for i, points := range []map[string]image.Point{
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 0},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{0, 150},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 150},
// 			"dest":   image.Point{200, 150},
// 			"canvas": image.Point{200, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{300, 150},
// 			"dest":   image.Point{300, 225},
// 			"canvas": image.Point{300, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{500, 150},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 150},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 200},
// 			"dest":   image.Point{266, 200},
// 			"canvas": image.Point{200, 200},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{200, 400},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{200, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{500, 450},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{600, 450},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 		map[string]image.Point{
// 			"source": image.Point{400, 300},
// 			"target": image.Point{700, 450},
// 			"dest":   image.Point{400, 300},
// 			"canvas": image.Point{400, 300},
// 		},
// 	} {
// 		source := points["source"]
// 		target := points["target"]
// 		dest := points["dest"]
// 		canvas := points["canvas"]
// 		f, err := storage.NewImage(map[string][]string{
// 			"url":    []string{"http://0.0.0.0"},
// 			"method": []string{"thumbnail"},
// 			"width":  []string{strconv.Itoa(target.X)},
// 			"height": []string{strconv.Itoa(target.Y)},
// 		}, o.Hosts)
// 		if err != nil {
// 			t.Fatalf("should be valid image: error=%v", err)
// 		}
// 		if f, err := f.Normalize(source); err != nil {
// 			t.Errorf("Normalize with thumbnail method should not return error: %v", err)
// 		} else if f.DestWidth != dest.X {
// 			t.Errorf("DestWidth(%d) expected %d, but actual %d", i, dest.X, f.DestWidth)
// 		} else if f.DestHeight != dest.Y {
// 			t.Errorf("DestHeight(%d) expected %d, but actual %d", i, dest.Y, f.DestHeight)
// 		} else if f.CanvasWidth != canvas.X {
// 			t.Errorf("CanvasWidth(%d) expected %d, but actual %d", i, canvas.X, f.CanvasWidth)
// 		} else if f.CanvasHeight != canvas.Y {
// 			t.Errorf("CanvasHeight(%d) expected %d, but actual %d", i, canvas.Y, f.CanvasHeight)
// 		}
// 	}
// }
