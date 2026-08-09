package flags

import "regexp"

func init() {
	Register("C", &Detector{
		Defaults: []string{"-Wall", "-Wextra"},
		Rules: []Rule{
			// Threading
			{regexp.MustCompile(`pthread_create|pthread_t|pthread_join|#include\s*<pthread\.h>`), []string{"-pthread"}},
			// Math library
			{regexp.MustCompile(`#include\s*<math\.h>|#include\s*<cmath>`), []string{"-lm"}},
			// POSIX real-time extensions
			{regexp.MustCompile(`#include\s*<aio\.h>|#include\s*<mqueue\.h>|#include\s*<semaphore\.h>`), []string{"-lrt"}},
			// Dynamic loading
			{regexp.MustCompile(`dlopen|dlsym|dlerror|dlclose|#include\s*<dlfcn\.h>`), []string{"-ldl"}},
			// Socket networking (not always needed, but common)
			{regexp.MustCompile(`#include\s*<sys/socket\.h>|#include\s*<netinet/|#include\s*<arpa/`), []string{}},
			// C11 atomics
			{regexp.MustCompile(`#include\s*<stdatomic\.h>|_Atomic\s`), []string{"-std=c11"}},
			// C99 VLA or designated init patterns
			{regexp.MustCompile(`#include\s*<stdbool\.h>|#include\s*<stdint\.h>|#include\s*<stdio\.h>`), []string{}},
		},
	})
}
