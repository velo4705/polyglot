package flags

import "regexp"

func init() {
	Register("C++", &Detector{
		Defaults: []string{"-Wall", "-Wextra"},
		Rules: []Rule{
			// C++17 features
			{regexp.MustCompile(`std::(optional|variant|any|make_optional|make_variant|get_if|holds_alternative)`), []string{"-std=c++17"}},
			{regexp.MustCompile(`std::filesystem|std::fs::`), []string{"-std=c++17"}},
			{regexp.MustCompile(`std::(string_view|clamp|apply|invoke|sample)`), []string{"-std=c++17"}},
			{regexp.MustCompile(`if\s*\(\s*auto|if\s+constexpr`), []string{"-std=c++17"}},
			{regexp.MustCompile(`structured_bindings|std::(tie|tuple_size|tuple_element)`), []string{"-std=c++17"}},

			// C++20 features
			{regexp.MustCompile(`std::(format|ranges|views|ranges::views)`), []string{"-std=c++20"}},
			{regexp.MustCompile(`concept\s+\w+|requires\s*\(|requires\s*\{`), []string{"-std=c++20"}},
			{regexp.MustCompile(`co_(await|yield|return)|co_await\s`), []string{"-std=c++20"}},
			{regexp.MustCompile(`std::(span|source_location|coroutine_handle|jthread)`), []string{"-std=c++20"}},
			{regexp.MustCompile(`operator\s+<=>|<=>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<format>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<ranges>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<concepts>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<coroutine>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<span>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<source_location>`), []string{"-std=c++20"}},
			{regexp.MustCompile(`#include\s*<bit>`), []string{"-std=c++20"}},

			// C++23 features
			{regexp.MustCompile(`import\s+std;|std::(expected|mdspan|flat_map|flat_set|stacktrace)`), []string{"-std=c++23"}},
			{regexp.MustCompile(`#include\s*<expected>`), []string{"-std=c++23"}},
			{regexp.MustCompile(`#include\s*<mdspan>`), []string{"-std=c++23"}},
			{regexp.MustCompile(`#include\s*<print>`), []string{"-std=c++23"}},
			{regexp.MustCompile(`std::print|std::println`), []string{"-std=c++23"}},

			// Threading
			{regexp.MustCompile(`std::(thread|mutex|condition_variable|future|promise|async|atomic)`), []string{"-pthread"}},
			{regexp.MustCompile(`#include\s*<thread>|#include\s*<mutex>|#include\s*<future>`), []string{"-pthread"}},

			// Math library
			{regexp.MustCompile(`#include\s*<cmath>|std::(sin|cos|sqrt|pow|log|exp|floor|ceil|round)\b`), []string{"-lm"}},

			// Dynamic loading
			{regexp.MustCompile(`dlopen|dlsym|#include\s*<dlfcn\.h>`), []string{"-ldl"}},

			// POSIX real-time
			{regexp.MustCompile(`#include\s*<aio\.h>|#include\s*<semaphore\.h>`), []string{"-lrt"}},

			// SIMD / SSE / AVX
			{regexp.MustCompile(`#include\s*<(emmintrin|immintrin|xmmintrin|pmmintrin|smmintrin|nmmintrin|tmmintrin|ammintrin|wmmintrin)\.h>`), []string{"-msse4.1"}},
			{regexp.MustCompile(`_mm\d{2,}_|_mm256_|_mm512_`), []string{"-mavx2"}},
			{regexp.MustCompile(`_mm512_`), []string{"-mavx512f"}},

			// OpenMP
			{regexp.MustCompile(`#pragma\s+omp|omp_get_thread_num|omp_get_num_threads`), []string{"-fopenmp"}},
		},
	})
}
