package utils

import (
	"fmt"
	"strings"
	"time"
)

func PrependTimestamp(str string) string {
	var b strings.Builder
	now := time.Now()
	fmt.Fprintf(&b, "%s %s", now.Format("15:04:05"), str)
	return b.String()
}
