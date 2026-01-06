/*
    file:           utils/path.go
    description:    Utilitas pendukung untuk path
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package utils

import (
	"os"
	"path/filepath"
)

// ResolvePath adalah fungsi untuk resolve jalur.
func ResolvePath(rel string) string {
	if rel == "" || filepath.IsAbs(rel) {
		return rel
	}

	candidates := []string{}
	if wd, err := os.Getwd(); err == nil && wd != "" {
		candidates = append(candidates, wd)
		parent := wd
		for i := 0; i < 3; i++ {
			next := filepath.Dir(parent)
			if next == "" || next == parent {
				break
			}
			candidates = append(candidates, next)
			parent = next
		}
	}
	if exe, err := os.Executable(); err == nil && exe != "" {
		exeDir := filepath.Dir(exe)
		if exeDir != "" {
			candidates = append(candidates, exeDir)
			parent := filepath.Dir(exeDir)
			if parent != "" && parent != exeDir {
				candidates = append(candidates, parent)
			}
		}
	}

	seen := map[string]struct{}{}
	for _, base := range candidates {
		if base == "" {
			continue
		}
		if _, ok := seen[base]; ok {
			continue
		}
		seen[base] = struct{}{}

		full := filepath.Join(base, rel)
		if _, err := os.Stat(full); err == nil {
			return full
		}
	}

	return rel
}
