package tscompiler

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path"
)

func Compile(src, destDir, iconName string, resultCh chan<- TsResult) {
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
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
			resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
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

	jsonEncoded, err := json.MarshalIndent(curr.node, "", "  ")
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
	}

	destJsFile, err := os.OpenFile(path.Join(destDir, iconName+".js"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer destJsFile.Close()
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
		return
	}

	_, err = destJsFile.Write(append([]byte(`export var `+iconName+` = `), jsonEncoded...))
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
		return
	}

	destTsFile, err := os.OpenFile(path.Join(destDir, iconName+".d.ts"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer destTsFile.Close()
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
		return
	}

	_, err = destTsFile.Write([]byte(`import { IcommonNode } from "@icommon/components/types";
export declare const ` + iconName + `: IcommonNode;`))
	if err != nil {
		resultCh <- TsResult{Err: &TsError{Details: TsDetails{Dest: destDir}, Err: err}}
		return
	}

	resultCh <- TsResult{Success: &TsDetails{Dest: destDir}}
}
