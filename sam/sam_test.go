package sam

import (
	"fmt"
	"testing"
)

func TestNewSAM(t *testing.T) {
	sam := NewSAM("abcbc")
	fmt.Println(sam.SubStrCount())
}
