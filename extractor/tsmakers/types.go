package tsmakers

type Maker func(src, dest, iconName string, resultCh chan<- MakeResult)

type MakeResult struct {
	Success *MakeDetails
	Err     *MakeError
}

type MakeError struct {
	Details MakeDetails
	Err     error
}

type MakeDetails struct {
	Dest string
}
