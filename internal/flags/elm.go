package flags

import "regexp"

func init() {
	Register("Elm", &Detector{
		Defaults: []string{},
		Rules: []Rule{
			{regexp.MustCompile(`Browser\.application`), []string{}},
			{regexp.MustCompile(`Browser\.element`), []string{}},
		},
	})
}
