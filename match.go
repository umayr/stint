package stint

import "strings"

func matchTitle(title, name string) bool {
	return strings.Contains(strings.ToLower(title), strings.ToLower(name))
}

func matchQuality(title string, quality int) bool {
	title = strings.ToLower(title)

	switch quality {
	case Normal:
		return !strings.Contains(title, "720p") && !strings.Contains(title, "1080p")
	case Medium:
		return strings.Contains(title, "720p")
	case High:
		return strings.Contains(title, "1080p")
	}

	return false
}
