package test

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	var a *int
	b := 6
	a = &b
	fmt.Print(a)

}
