package stringutils

import (
	"fmt"
	"testing"
)

func Test_RemoveSliceElement(t *testing.T) {
	s := make([]string, 0, 10)
	s = append(s, "a")
	s = append(s, "b")
	s = append(s, "c")
	s = append(s, "d")
	s = RemoveSliceElement(s, 1).([]string)

	fmt.Println(s)
}

func Test_Substr(t *testing.T) {
	s := "asldfjasl;djf;"
	s = Substr(s, 0, len(s)-4)
	fmt.Println(s)
}
