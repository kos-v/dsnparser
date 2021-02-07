package dsnparser

func isEscaped(pos int, target []rune) bool {
	return pos > 0 && target[pos-1] == '\\'
}

func searchRuneInArray(target []rune, needle rune) int {
	for i, item := range target {
		if needle == item {
			return i
		}
	}

	return -1
}

func unEscape(needs []rune, target []rune) []rune {
	var unescaped []rune

	for symbolPos, symbol := range target {
		if symbol == '\\' {
			if symbolPos+1 < len(target) {
				if searchRuneInArray(needs, target[symbolPos+1]) != -1 {
					continue
				}
			}
			unescaped = append(unescaped, '\\')
		} else {
			unescaped = append(unescaped, symbol)
		}
	}

	return unescaped
}
