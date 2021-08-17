package util

import (
	"path/filepath"
	"strings"
)

func ListMatchFiles(root string, pattern string, ignores []string) []string {
	files := []string{}
	pattern = filepath.Join(root, pattern)
	if !strings.HasPrefix(pattern, root) {
		return files
	}
	tmp, _ := filepath.Glob(pattern)
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
		if !match {
			files = append(files, rel)
		}
	}
	return files
}

func ListMultiMatchFiles(root string, patterns []string, ignores []string) []string {
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
			if !match {
				files = append(files, rel)
			}
		}
	}
	return files
}

func FilenameNoExt(file string) string {
	return file[:len(file)-len(filepath.Ext(file))]
}
