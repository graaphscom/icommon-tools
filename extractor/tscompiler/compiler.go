package tscompiler

import (
	"os"
)

func Compile(src, dest, iconName string, resultCh chan<- TsResult) {
	contents, err := os.ReadFile(src)
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
		return
	}

	//file, err := os.OpenFile(path.Join(dest), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	//if err != nil {
	//	resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
	//	return
	//}

	resultCh <- TsResult{Success: &TsDetails{Dest: dest}}
}
