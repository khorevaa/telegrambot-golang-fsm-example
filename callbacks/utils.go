// some helpers
package callbacks

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func randInt(start int, stop int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return start + rand.Intn(int(math.Abs(float64(start-stop))))
}


func HtmlFmt(s string, format string) string {
	return fmt.Sprintf("<%s>%s</%s>", format, s, format)
}
