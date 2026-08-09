package flags

import "regexp"

func init() {
	Register("Fortran", &Detector{
		Defaults: []string{"-Wall"},
		Rules: []Rule{
			{regexp.MustCompile(`open\s*\(|read\s*\(|write\s*\(|print\s*\*`), []string{"-std=f2008"}},
			{regexp.MustCompile(`use\s+iso_c_binding|use\s+iso_fortran_env`), []string{"-std=f2008"}},
			{regexp.MustCompile(`submodule\s*\(`), []string{"-std=f2008"}},
			{regexp.MustCompile(`function\s+\w+\s*\(`), []string{}},
		},
	})
}
