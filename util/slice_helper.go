package util

func UniqueUint(UintSlice []uint) []uint {
	keys := make(map[uint]bool)
	list := []uint{}
	for _, entry := range UintSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func UniqueString(StringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range StringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
