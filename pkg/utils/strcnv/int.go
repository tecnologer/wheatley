package strcnv

import "strconv"

func ParseIntWithDefault(str string, def int) int {
	if str == "" {
		return def
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return def
	}

	return i
}
