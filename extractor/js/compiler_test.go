package js

import (
	"os"
	"testing"
)

func TestCompile(t *testing.T) {
	resultCh := make(chan<- CompileResult, 1)
	Compile("./testdata/moveDownSharp24px.svg", "./testdata", "moveDownSharp24", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, resultCh)
}
