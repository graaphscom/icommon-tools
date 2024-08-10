package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/db"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"log"
	"os"
	"strings"
)

func main() {
	manifestPath := flag.String("manifest", "", "path to the icons manifest file")
	noTrunc := flag.Bool("no-trunc", false, "don't write a file when it already exists")
	flag.Parse()

	manifest, err := js.ReadJson[js.IcoManifest](*manifestPath)
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

	caseInsensitiveOccurrences := make(map[string][]string)
	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		for _, icon := range iconSet.Icons {
			iconLowerNameFullPath := strings.Join(append(segments, strings.ToLower(icon.Name)), ":")
			caseInsensitiveOccurrences[iconLowerNameFullPath] = append(caseInsensitiveOccurrences[iconLowerNameFullPath], icon.Name)
		}
	})
	var duplicates [][]string
	for _, cases := range caseInsensitiveOccurrences {
		if len(cases) > 1 {
			duplicates = append(duplicates, cases)
		}
	}

	ctx := context.Background()
	commands[commandsIdx] = db.CreateSearchIndex(conn.B())

	conn.DoMulti(ctx, commands...)

	if err := os.Mkdir(manifest.TsResultPath, 0750); err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatalln(err)
	}

	openFlag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if *noTrunc {
		openFlag = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	err = tree.Traverse(
		[]string{},
		func(segments []string, iconSet unitree.IconSet) error {
			packagePath := js.BuildPackagePath(manifest.TsResultPath, segments)

			if err := os.MkdirAll(packagePath, 0750); err != nil && !os.IsExist(err) {
				return err
			}

			var iconsName []struct{ Name string }
			for _, icon := range iconSet.Icons {
				iconsName = append(iconsName, struct{ Name string }{icon.Name})
			}

			err = js.CompileIndex(js.IndexTs(packagePath), openFlag, iconsName)
			if err != nil {
				return err
			}

			err = js.CompileIndex(js.IndexJs(packagePath), openFlag, iconsName)
			if err != nil {
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
			js.Compile(icon.SrcFile, packagePath, icon.Name, openFlag, resultCh)
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
