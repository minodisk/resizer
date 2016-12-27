package storage

func in(target string, elements []string) bool {
	for _, e := range elements {
		if target == e {
			return true
		}
	}
	return false
}
