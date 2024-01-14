package server

import "strings"

func makeFilename(name string) string {
	lowered := strings.ToLower(name)
	parts := strings.Fields(lowered)
	return strings.Join(parts, "-")
}
