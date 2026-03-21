package service

import (
	"github.com/microcosm-cc/bluemonday"
)

var (
	pStrict = bluemonday.StrictPolicy()
)

func strictSanitization(input string) string {
	return pStrict.Sanitize(input)
}
