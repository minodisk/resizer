package log

import (
	"fmt"
	"time"
)

type Gorm struct{}

func (self *Gorm) Print(values ...interface{}) {
	level := values[0]
	if level == "sql" {
		entryd(values[2].(time.Duration)).Debug(fmt.Sprintf("%s %s", values[3], values[4]))
	} else {
		entry(nil).Debug(values[2:]...)
	}
}
