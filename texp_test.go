package texp

import (
	"time"
	"fmt"
	"testing"
	"github.com/qjsoftcn/texp"
)

func Test_Parser(t *testing.T) {
	parser := texp.NewParser(time.Now(), time.Now())
	td, err := parser.Parse("d-1d")
	fmt.Println(td, err)
}
