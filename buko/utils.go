package server

import "strings"

func makeURLSlug(name string) string {
	lowered := strings.ToLower(name)
	parts := strings.Fields(lowered)
	return strings.Join(parts, "-")
}
