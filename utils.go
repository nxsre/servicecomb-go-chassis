package servicecomb

import (
	"strings"
)

func Path(path string, m map[string]string) string {
	tmps := strings.Split(path, "/")
	for _, s := range tmps {
		if len(s) > 1 && s[0] == byte(':') {
			k := string(s[1:])
			path = strings.Replace(path, s, m[k], -1)
		}
	}
	return path
}
