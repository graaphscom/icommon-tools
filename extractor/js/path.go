package js

import "path"

func BuildPackagePath(base string, segments []string) string {
	if len(segments) < 2 {
		return base
	}
	if len(segments) == 2 {
		return path.Join(base, segments[1])
	}
	return path.Join(base, segments[1], path.Join(segments[2:]...))
}

func IndexJs(base string) string {
	return path.Join(base, "index.js")
}

func IndexTs(base string) string {
	return path.Join(base, "index.d.ts")
}

func FileJs(base, filename string) string {
	return path.Join(base, filename+".js")
}

func FileTs(base, filename string) string {
	return path.Join(base, filename+".d.ts")
}
