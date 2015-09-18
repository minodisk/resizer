package log_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/go-microservices/resizer/log"
)

var (
	file = "log/test.log"
)

func TestMain(m *testing.M) {
	code := m.Run()
	defer os.Exit(code)

	log.Close()
	os.RemoveAll(file)
}

func TestStartAndEnd(t *testing.T) {
	os.Setenv(log.EnvLogFilename, file)
	filename, err := log.Init()
	if err != nil {
		t.Fatalf("fail to init: error=%v", err)
	}
	if err := os.Unsetenv(log.EnvLogFilename); err != nil {
		t.Fatalf("fail to unset env: error=%v", err)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	timer := log.Start()

	log.Debug("debug")
	log.Print("print")
	log.Info("info")
	log.Warn("warn")
	log.Warning("warning")
	log.Error("error")

	log.Debugf("debug%s", "f")
	log.Printf("print%s", "f")
	log.Infof("info%s", "f")
	log.Warnf("warn%s", "f")
	log.Warningf("warning%s", "f")
	log.Errorf("error%s", "f")

	log.Debugln("debugln")
	log.Println("println")
	log.Infoln("infoln")
	log.Warnln("warnln")
	log.Warningln("warningln")
	log.Errorln("errorln")

	log.End(timer)

	expected := []map[string]string{
		map[string]string{"level": "info", "msg": "start"},

		map[string]string{"level": "info", "msg": "print"},
		map[string]string{"level": "info", "msg": "info"},
		map[string]string{"level": "warning", "msg": "warn"},
		map[string]string{"level": "warning", "msg": "warning"},
		map[string]string{"level": "error", "msg": "error"},

		map[string]string{"level": "info", "msg": "printf"},
		map[string]string{"level": "info", "msg": "infof"},
		map[string]string{"level": "warning", "msg": "warnf"},
		map[string]string{"level": "warning", "msg": "warningf"},
		map[string]string{"level": "error", "msg": "errorf"},

		map[string]string{"level": "info", "msg": "println"},
		map[string]string{"level": "info", "msg": "infoln"},
		map[string]string{"level": "warning", "msg": "warnln"},
		map[string]string{"level": "warning", "msg": "warningln"},
		map[string]string{"level": "error", "msg": "errorln"},

		map[string]string{"level": "info", "msg": "end"},
	}

	bufs, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("fail to read file: error=%v", err)
	}
	lines := strings.Split(string(bufs), "\n")

	l := len(lines)
	for i, line := range lines {
		// fmt.Println(i, line)

		if i == l-1 {
			break
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Fatalf("fail to decode json")
		}

		e := expected[i]
		// fmt.Println(i, e)
		if obj["level"] != e["level"] {
			t.Fatalf("wrong level: expected %s, but actual %s", e["level"], obj["level"])
		}
		if obj["msg"] != e["msg"] {
			t.Fatalf("wrong msg: expected %s, but actual %s", e["msg"], obj["msg"])
		}
	}

	if err := log.Close(); err != nil {
		t.Fatalf("fail to close log: error=%v", err)
	}
}
