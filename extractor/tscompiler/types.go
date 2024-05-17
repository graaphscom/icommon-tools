package tscompiler

import "encoding/json"

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
