package main

import (
	"context"
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/db"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"log"
	"os"
	"text/template"
)

func main() {
	manifest, err := js.ReadJson[js.IcoManifest]("testdata/ico_manifest_downloads.json")
	tree, err := unitree.BuildRootTree(manifest)

	if err != nil {
		log.Fatalln(err)
	}

	opt, err := rueidis.ParseURL(manifest.RedisUrl)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := rueidis.NewClient(opt)
	if err != nil {
		log.Fatalln(err)
	}

	iconsCount := 0
	tree.MustTraverse([]string{}, func(_ []string, iconSet unitree.IconSet) {
		for range iconSet.Icons {
			iconsCount++
		}
	})

	commands := make([]rueidis.Completed, iconsCount+1)
	commandsIdx := 0
	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		for _, icon := range iconSet.Icons {
			commands[commandsIdx] = db.CreateIconEntry(conn.B(), segments, icon)
			commandsIdx++
		}
	})

	ctx := context.Background()
	commands[commandsIdx] = db.CreateSearchIndex(conn.B())

	conn.DoMulti(ctx, commands...)

	if err := os.Mkdir(manifest.TsResultPath, 0750); err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}

	tsIndexFileTmpl, err := template.New("tsIndexFileTpl").Parse(`{{range .}}export { {{.Name}} } from "./{{.Name}}";
{{end}}`)

	if err != nil {
		log.Fatalln(err)
	}

	err = tree.Traverse(
		[]string{},
		func(segments []string, iconSet unitree.IconSet) error {
			packagePath := js.BuildPackagePath(manifest.TsResultPath, segments)

			if err := os.MkdirAll(packagePath, 0750); err != nil && !os.IsExist(err) {
				return err
			}

			dTsFile, err := os.OpenFile(js.IndexTs(packagePath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}

			jsFile, err := os.OpenFile(js.IndexJs(packagePath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}

			if err := tsIndexFileTmpl.Execute(dTsFile, iconSet.Icons); err != nil {
				return err
			}
			if err := tsIndexFileTmpl.Execute(jsFile, iconSet.Icons); err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		log.Fatalln(err)
	}

	resultCh := make(chan js.CompileResult, iconsCount)

	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		packagePath := js.BuildPackagePath(manifest.TsResultPath, segments)
		for _, icon := range iconSet.Icons {
			js.Compile(icon.SrcFile, packagePath, icon.Name, resultCh)
		}
	})

	successCount := 0
	errCount := 0

	for range iconsCount {
		writeResult := <-resultCh
		if writeResult.Success != nil {
			successCount++
		}
		if writeResult.Err != nil {
			fmt.Println(writeResult.Err)
			errCount++
		}
	}

	fmt.Printf("Total icons count: %d\ntscompiler success: %d\ntscompiler errors: %d", iconsCount, successCount, errCount)
}
