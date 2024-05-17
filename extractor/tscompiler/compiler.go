package tscompiler

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path"
)

func Compile(src, dest, iconName string, resultCh chan<- TsResult) {
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
		return
	}

	decoder := xml.NewDecoder(srcFile)

	curr := struct {
		node   *IcommonNode
		parent *IcommonNode
	}{}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
		}

		if startElement, ok := token.(xml.StartElement); ok {
			if curr.node == nil {
				curr.node = &IcommonNode{Attributes: make(map[string]string)}
			} else {
				newNode := &IcommonNode{Attributes: make(map[string]string)}
				curr.node.Children = append(curr.node.Children, newNode)
				curr.parent = curr.node
				curr.node = newNode
			}

			curr.node.Name = startElement.Name.Local
			for _, attr := range startElement.Attr {
				curr.node.Attributes[attr.Name.Local] = attr.Value
			}
		}

		if _, ok := token.(xml.EndElement); ok {
			curr.node = curr.parent
		}
	}

	jsonEncoded, err := json.MarshalIndent(curr.node, "", "    ")
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
	}

	destFile, err := os.OpenFile(path.Join(dest), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer destFile.Close()
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
		return
	}

	_, err = destFile.Write(append([]byte(`export const `+iconName+` = `), jsonEncoded...))
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: dest}, Err: err}}
		return
	}

	resultCh <- TsResult{Success: &TsDetails{Dest: dest}}
}
