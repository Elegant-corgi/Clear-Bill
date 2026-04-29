package bll

import "regexp"

var chineseCharacterPattern = regexp.MustCompile(`[\p{Han}]`)

func containsChineseCharacters(value string) bool {
	return chineseCharacterPattern.MatchString(value)
}
