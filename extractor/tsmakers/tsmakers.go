package tsmakers

import (
	"errors"
	"os"
	"path"
	"regexp"
	"text/template"
)

var Boxicons Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Bytesize Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Fluentui Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Fontawesome Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	contents, err := os.ReadFile(src)
	if err != nil {
		resultCh <- MakeResult{Err: &MakeError{Details: MakeDetails{Dest: dest}, Err: err}}
		return
	}
	svgViewBoxMatch := svgViewBoxRegexp.FindSubmatch(contents)
	if len(svgViewBoxMatch) != 2 {
		resultCh <- MakeResult{Err: &MakeError{Details: MakeDetails{Dest: dest}, Err: errors.New("invalid svg viewbox regexp match")}}
		return
	}
	svgPathMatch := svgPathRegexp.FindSubmatch(contents)
	if len(svgPathMatch) != 2 {
		resultCh <- MakeResult{Err: &MakeError{Details: MakeDetails{Dest: dest}, Err: errors.New("invalid svg path regexp match")}}
		return
	}
	copyrightNoteMatch := copyrightRegexp.FindSubmatch(contents)
	if len(copyrightNoteMatch) != 2 {
		resultCh <- MakeResult{Err: &MakeError{Details: MakeDetails{Dest: dest}, Err: errors.New("invalid copyright note regexp match")}}
		return
	}

	file, err := os.OpenFile(path.Join(dest), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		resultCh <- MakeResult{Err: &MakeError{Details: MakeDetails{Dest: dest}, Err: err}}
		return
	}

	if err := tsIconFileTmpl.Execute(file, struct {
		CopyrightNote string
		IconName      string
		SvgPath       string
		SvgViewBox    string
	}{
		string(copyrightNoteMatch[1]),
		iconName,
		string(svgPathMatch[1]),
		string(svgViewBoxMatch[1]),
	}); err != nil {
		resultCh <- MakeResult{Err: &MakeError{Details: MakeDetails{Dest: dest}, Err: err}}
		return
	}

	resultCh <- MakeResult{Success: &MakeDetails{Dest: dest}}
}

var svgViewBoxRegexp, _ = regexp.Compile("viewBox=\"([^\"]+)\"")
var svgPathRegexp, _ = regexp.Compile("d=\"([^\"]+)\"")
var copyrightRegexp, _ = regexp.Compile("<!--! (.*) -->")
var tsIconFileTmpl, _ = template.New("tsIconFileTpl").Parse(`// {{.CopyrightNote}}
export const {{.IconName}} = {
    d: "{{.SvgPath}}",
    viewBox: "{{.SvgViewBox}}",
};
`)

var Material Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Octicons Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Radixui Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Remixicon Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}

var Unicons Maker = func(src, dest, iconName string, resultCh chan<- MakeResult) {
	resultCh <- MakeResult{Err: &MakeError{Err: errors.New("not implemented"), Details: MakeDetails{Dest: dest}}}
}
