package js

import "encoding/json"

type CompileResult struct {
	Success *CompileDetails
	Err     *CompileError
}

type CompileError struct {
	Details CompileDetails
	Err     error
}

type CompileDetails struct {
	Dest string
}

type IcommonNode struct {
	Name       string
	Attributes map[string]string
	Children   []*IcommonNode
}

func (i IcommonNode) MarshalJSON() ([]byte, error) {
	attributes, err := json.Marshal(i.Attributes)
	if err != nil {
		return nil, err
	}

	if i.Children == nil {
		return []byte(`["` + i.Name + `", ` + string(attributes) + `]`), nil
	}

	children, err := json.Marshal(i.Children)

	return []byte(`["` + i.Name + `", ` + string(attributes) + `, ` + string(children) + `]`), nil
}
