package faye

import "path/filepath"

func fileAddr(rawPath, name string) string {
	res := filepath.Join(rawPath, name)
	return res
}
