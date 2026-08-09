package flags

import "regexp"

func init() {
	Register("V", &Detector{
		Defaults: []string{},
		Rules: []Rule{
			{regexp.MustCompile(`#flag\s+-l`), []string{}},
			{regexp.MustCompile(`\[live\]`), []string{"-live"}},
		},
	})
}
