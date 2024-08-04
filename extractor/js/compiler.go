package js

import (
	"encoding/json"
	"encoding/xml"
	"github.com/graaphscom/icommon-tools/extractor/strcase"
	"io"
	"os"
	"strings"
	"text/template"
)

func Compile(src, destDir, iconName string, openFlag int, resultCh chan<- CompileResult) {
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		resultCh <- CompileResult{Err: &CompileError{Details: CompileDetails{Dest: destDir}, Err: err}}
		return
	}

	decoder := xml.NewDecoder(srcFile)

	var curr *IcommonNode
	var chain []*IcommonNode

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			resultCh <- CompileResult{Err: &CompileError{Details: CompileDetails{Dest: destDir}, Err: err}}
		}

		if startElement, ok := token.(xml.StartElement); ok {
			newNode := &IcommonNode{Attributes: make(map[string]string)}
			newNode.Name = startElement.Name.Local
			for _, attr := range startElement.Attr {
				if strings.ToLower(attr.Name.Local) == "class" || strings.HasPrefix(attr.Name.Local, "data") {
					continue
				}
				newNode.Attributes[strcase.ToCamel(attr.Name.Local, strcase.KebabRegexp)] = attr.Value
			}

			if curr == nil {
				curr = newNode
			} else {
				curr.Children = append(curr.Children, newNode)
				curr = newNode
			}
			chain = append(chain, curr)
		}

		if _, ok := token.(xml.EndElement); ok {
			chain = chain[:len(chain)-1]
			if len(chain) > 0 {
				curr = chain[len(chain)-1]
			}
		}
	}

	jsonEncoded, err := json.MarshalIndent(curr, "", "  ")
	if err != nil {
		resultCh <- CompileResult{Err: &CompileError{Details: CompileDetails{Dest: destDir}, Err: err}}
	}

	writeDest(
		FileJs(destDir, iconName),
		openFlag,
		append([]byte(`export var `+iconName+` = `), jsonEncoded...),
		resultCh,
	)
	writeDest(
		FileTs(destDir, iconName),
		openFlag,
		[]byte(`import { IcommonNode } from "@icommon/components/types";
export declare const `+iconName+`: IcommonNode;`),
		resultCh,
	)

	resultCh <- CompileResult{Success: &CompileDetails{Dest: destDir}}
}

func writeDest(dest string, openFlag int, toWrite []byte, resultCh chan<- CompileResult) {
	destFile, err := os.OpenFile(dest, openFlag, 0666)
	if os.IsExist(err) {
		return
	}
	defer destFile.Close()
	if err != nil {
		resultCh <- CompileResult{Err: &CompileError{Details: CompileDetails{Dest: dest}, Err: err}}
		return
	}

	_, err = destFile.Write(toWrite)
	if err != nil {
		resultCh <- CompileResult{Err: &CompileError{Details: CompileDetails{Dest: dest}, Err: err}}
		return
	}
}

var indexFileTmpl *template.Template

func CompileIndex(dest string, openFlag int, iconsName []struct{ Name string }) error {
	var err error
	if indexFileTmpl == nil {
		indexFileTmpl, err = template.New("tsIndexFileTpl").Parse(`{{range .}}export { {{.Name}} } from "./{{.Name}}";
{{end}}`)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(dest, openFlag, 0666)
	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return indexFileTmpl.Execute(file, iconsName)
}
