package options

import "strings"

type Hosts []string

func (hs *Hosts) String() string {
	return strings.Join(*hs, ", ")
}

func (hs *Hosts) Set(host string) error {
	*hs = append(*hs, host)
	return nil
}

func (hs Hosts) Contains(host string) bool {
	if len(hs) == 0 {
		return true
	}
	for _, h := range hs {
		if h == host {
			return true
		}
	}
	return false
}
