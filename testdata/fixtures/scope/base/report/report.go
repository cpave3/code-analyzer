package report

import "strings"

// Header renders a title.
func Header(title string) string { return "== " + strings.ToUpper(title) + " ==" }

// Footer renders a page footer.
func Footer(page, of int) string { return strings.Repeat("-", 10) + pageOf(page, of) }

func pageOf(page, of int) string { return " p" + itoa(page) + "/" + itoa(of) }

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
