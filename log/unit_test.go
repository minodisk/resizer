package log_test

import (
	"testing"
	"time"

	"github.com/go-microservices/resizer/log"
)

func TestConvert(t *testing.T) {
	for arg, expected := range map[string]string{
		"1ns":    "1.0ns",
		"999ns":  "999.0ns",
		"1000ns": "1.0μs",

		"1us":        "1.0μs",
		"1us999ns":   "1.9μs",
		"999us999ns": "999.9μs",
		"1000us":     "1.0ms",

		"1ms":             "1.0ms",
		"1ms999us999ns":   "1.9ms",
		"999ms999us999ns": "999.9ms",
		"1000ms":          "1.0s",

		"1s":                 "1.0s",
		"1s999ms999us999ns":  "1.9s",
		"59s999ms999us999ns": "59.9s",
		"60s":                "1.0m",

		"1m": "1.0m",
		"1m59s999ms999us999ns":  "1.9m",
		"59m59s999ms999us999ns": "59.9m",
		"60m": "1.0h",

		"1h": "1.0h",
		"1h59m59s999ms999us999ns":  "1.9h",
		"23h59m59s999ms999us999ns": "23.9h",
		"24h": "1.0d",

		"24h59m59s999ms999us999ns": "1.0d",
		"25h": "1.0d",
		"26h": "1.0d",
		"27h": "1.1d",
		"167h59m59s999ms999us999ns": "6.9d",

		"168h":  "1.0w",
		"252h":  "1.5w",
		"1680h": "10.0w",
	} {
		d, err := time.ParseDuration(arg)
		if err != nil {
			t.Fatalf("fail to parse as duration: arg=%s, error=%v", arg, err)
		}
		if actual := log.Convert(d); actual != expected {
			t.Fatalf("fail to convert from %s: expected %s, but actual %s", arg, expected, actual)
		}
	}
}
