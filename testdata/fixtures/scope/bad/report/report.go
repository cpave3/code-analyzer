package report

import "strings"

// Header renders a report title in upper case, framed by equals signs.
func Header(title string) string {
	return "== " + strings.ToUpper(title) + " =="
}

// Footer renders a page footer with a separator rule.
func Footer(page, of int) string {
	return strings.Repeat("-", 10) + pageOf(page, of)
}

// pageOf renders " pN/M".
func pageOf(page, of int) string {
	return " p" + itoa(page) + "/" + itoa(of)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
