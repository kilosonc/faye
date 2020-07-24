package faye

import "path"

func fileAddr(rawPath, name string) string {
	res := path.Join(rawPath, name)
	return res
}
