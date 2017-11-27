package common

import (
	"strings"
)

// QMarks is a helper function for concatenating question mark for an SQL statement
func QMarks(n int) string {
	if n == 0 {
		return ""
	}
	qMark := strings.Repeat("?, ", n)
	qMark = qMark[:len(qMark)-2] // remove last ", "
	return qMark
}
