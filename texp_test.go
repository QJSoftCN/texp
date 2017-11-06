package texp

import (
	"time"
	"fmt"
	"testing"
)

func Test_Parser(t *testing.T) {
	parser := NewParser(time.Now(), time.Now())
	td, err := parser.Parse("d-1d")
	fmt.Println(td, err)


}
