package log

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

type Empty struct{}

const (
	EnvLogFilename = "RESIZER_LOG_FILENAME"
)

var (
	pkgPath  = reflect.TypeOf(Empty{}).PkgPath()
	basePath = func() string {
		ps := strings.Split(pkgPath, "/")
		ps = ps[:len(ps)-1]
		basePath := strings.Join(ps, "/")
		return basePath
	}()
	logger = logrus.New()
	file   *os.File
)

// Init はロガーを初期化します。
// ログ用のディレクトリが指定されない場合、Stdoutにログを出力します。
// 指定された場合、'ディレクトリ/日時.log'という名前のログファイルにログを出力します。
func Init() (string, error) {
	logFilename := os.Getenv(EnvLogFilename)
	if logFilename == "" {
		logger.Level = logrus.DebugLevel
		return "", nil
	}

	// ログ用のディレクトリが指定された場合、
	// ディレクトリを作成しログファイルをオープンする。
	if err := os.MkdirAll(filepath.Dir(logFilename), 0777); err != nil {
		return "", err
	}
	var err error
	file, err = os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}
	logger.Out = file
	logger.Formatter = new(logrus.JSONFormatter)
	return logFilename, nil
}

// Close はログファイルをオープンしている場合、ファイルをクローズします。
func Close() error {
	if file == nil {
		return nil
	}
	if err := file.Close(); err != nil {
		return err
	}
	file = nil
	return nil
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	entry(nil).Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	entry(nil).Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	entry(nil).Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	entry(nil).Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	entry(nil).Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	entry(nil).Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	entry(nil).Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	entry(nil).Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	entry(nil).Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	entry(nil).Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	entry(nil).Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	entry(nil).Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	entry(nil).Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	entry(nil).Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	entry(nil).Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	entry(nil).Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	entry(nil).Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	entry(nil).Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	entry(nil).Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	entry(nil).Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	entry(nil).Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	entry(nil).Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	entry(nil).Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	entry(nil).Fatalln(args...)
}

func Start() *timer {
	e := entry(nil)
	e.Info("start")
	return NewTimer()
}

func End(t *timer) {
	entryd(t.Now()).Println("end")
}

func entryd(d time.Duration) *logrus.Entry {
	return entry(logrus.Fields{"duration": Convert(d)})
}

// entry returns entry filled with the upper stack of caller.
func entry(fields logrus.Fields) *logrus.Entry {
	if fields == nil {
		fields = logrus.Fields{}
	}

	// スタックを遡りながらログのコール元を特定します。
	// アプリケーション内で、なるべく浅い、loggerパッケージではないパッケージを
	// コール元とし、出力のフィールドに含めます。
	for s := 1; ; s++ {
		pc, file, line, ok := runtime.Caller(s)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		name := f.Name()
		fields["file"] = file
		fields["line"] = line
		fields["name"] = name
		if strings.Index(name, pkgPath) != 0 && strings.Index(name, basePath) == 0 {
			break
		}
	}

	return logger.WithFields(fields)
}
