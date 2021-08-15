package util

import (
	"path/filepath"
	"strings"
)

func ListMatchFiles(root string, patterns []string, ignores []string) []string {
	files := []string{}
	for _, pat := range patterns {
		pat = filepath.Join(root, pat)
		if !strings.HasPrefix(pat, root) {
			continue
		}
		tmp, _ := filepath.Glob(pat)
		for _, file := range tmp {
			match := false
			rel, _ := filepath.Rel(root, file)
			base := filepath.Base(rel)
			for _, ign := range ignores {
				if ok, _ := filepath.Match(ign, rel); ok {
					match = true
					break
				}
				if ok, _ := filepath.Match(ign, base); ok {
					match = true
					break
				}
			}
			// fmt.Printf("file=%v, rel=%v, base=%v, match=%v\n", file, rel, base, match)
			if !match {
				files = append(files, rel)
			}
		}
	}
	return files
}
