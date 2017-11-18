package tpls

import "strings"

func qMarks(n int) string {
	if n == 0 {
		return ""
	}
	qMark := strings.Repeat("?, ", n)
	qMark = qMark[:len(qMark)-2] // remove last ", "
	return qMark
}
