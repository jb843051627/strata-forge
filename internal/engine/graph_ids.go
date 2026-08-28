package engine

import "strconv"

func fmtInt(value int64) string {
	return strconv.FormatInt(value, 10)
}
