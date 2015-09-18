package log

import "time"

type timer struct {
	from time.Time
}

func NewTimer() *timer {
	return &timer{time.Now()}
}

func (self *timer) Update() {
	self.from = time.Now()
}

func (self *timer) Now() time.Duration {
	t := time.Now()
	return t.Sub(self.from)
}
