package messages

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholder = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)

func Format(text string, vars map[string]string) string {
	if len(vars) == 0 {
		return text
	}

	result := placeholder.ReplaceAllStringFunc(text, func(match string) string {

		key := strings.Trim(match, "{}")

		if val, ok := vars[key]; ok {
			return val
		}

		return match
	})

	return result
}

func Mention(telegramID int64, firstName string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, telegramID, firstName)
}
