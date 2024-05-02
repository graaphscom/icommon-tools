package tscompiler

type TsResult struct {
	Success *TsDetails
	Err     *TsError
}

type TsError struct {
	Details TsDetails
	Err     error
}

type TsDetails struct {
	Dest string
}
