package parseerror

import (
	"fmt"
	"os"
)

var HadError bool

func Error(line int, message string) {
	HadError = true
	fmt.Fprintf(os.Stderr, "[line %v] Error : %v\n", line, message)
}
