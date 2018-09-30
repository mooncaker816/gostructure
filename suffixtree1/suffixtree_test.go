package suftree

import "testing"

func TestBuild(t *testing.T) {
	st := BuildSufTree("abaaba")
	st.Print()
}
