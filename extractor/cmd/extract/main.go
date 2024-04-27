package main

import (
	"context"
	"fmt"
	"github.com/graaphscom/icommon/extractor/json"
	"github.com/graaphscom/icommon/extractor/tsmakers"
	"github.com/graaphscom/icommon/extractor/unitree"
	"github.com/redis/rueidis"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

func main() {
	manifest, err := json.ReadJson[json.IcoManifest]("testdata/ico_manifest_downloads.json")
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
			commands[commandsIdx] = conn.B().Hset().
				Key(strings.Join(append(segments, icon.Name), ":")).
				FieldValue().
				FieldValue("searchTags", strings.Join(icon.Tags.Search, ",")).
				FieldValue("visualTags", strings.Join(icon.Tags.Visual, ",")).
				Build()
			commandsIdx++
		}
	})

	ctx := context.Background()
	commands[commandsIdx] = conn.B().FtCreate().
		Index("icommon").
		Prefix(1).Prefix("icommon:").
		Schema().
		FieldName("searchTags").Text().
		FieldName("visualTags").Tag().
		Build()

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
			joinedPath := path.Join(manifest.TsResultPath, path.Join(segments...))

			if err := os.Mkdir(joinedPath, 0750); err != nil && !os.IsExist(err) {
				return err
			}

			file, err := os.OpenFile(path.Join(joinedPath, "index.ts"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}

			if err := tsIndexFileTmpl.Execute(file, iconSet.Icons); err != nil {
				return err
			}
			return nil
		},
		func(segments []string) error {
			if err := os.Mkdir(path.Join(manifest.TsResultPath, path.Join(segments...)), 0750); err != nil && !os.IsExist(err) {
				return err
			}
			return nil
		},
	)

	if err != nil {
		log.Fatalln(err)
	}

	resultCh := make(chan tsmakers.MakeResult, iconsCount)

	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		for _, icon := range iconSet.Icons {
			iconSet.TsMaker(icon.SrcFile, path.Join(manifest.TsResultPath, path.Join(segments...), icon.Name+".ts"), icon.Name, resultCh)
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

	fmt.Printf("Total icons count: %d\ntsmakers success: %d\ntsmakers errors: %d", iconsCount, successCount, errCount)
}
