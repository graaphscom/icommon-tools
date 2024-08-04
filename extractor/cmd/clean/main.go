package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/db"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	printer := message.NewPrinter(language.Polish)

	dryRun := flag.Bool("dry-run", false, "don't perform actual clean, print redis keys and files to be deleted instead")
	flag.Parse()

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

	uniTreeKeys, uniTreeFiles := flattenUniTree(tree, manifest.TsResultPath)

	obsoleteKeys, err := findObsoleteKeys(conn, uniTreeKeys)
	if err != nil {
		log.Fatalln(err)
	}

	obsoleteFiles, err := findObsoleteFiles(manifest, uniTreeFiles)
	if err != nil {
		log.Fatalln(err)
	}

	if len(obsoleteKeys) == 0 {
		fmt.Println("No redis keys to delete\n")
	} else {
		fmt.Println("Redis keys to delete:")
		for _, obsoleteKey := range obsoleteKeys {
			fmt.Println(obsoleteKey)
		}
	}

	if len(obsoleteFiles) == 0 {
		fmt.Println("No files to delete")
	} else {
		fmt.Println("Files to delete:")
		for _, obsoleteFile := range obsoleteFiles {
			fmt.Println(obsoleteFile)
		}
	}

	if *dryRun {
		fmt.Println("dry-run - neither redis keys nor files have been deleted")
		return
	}

	if len(obsoleteKeys) > 0 {
		delKeysCount, err := conn.Do(context.Background(), conn.B().Del().Key(obsoleteKeys...).Build()).AsInt64()
		if err != nil {
			log.Fatalln(err)
		}

		_, err = printer.Printf("Successfully deleted %d redis keys\n", delKeysCount)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if len(obsoleteFiles) > 0 {
		for _, obsoleteFile := range obsoleteFiles {
			err = os.Remove(obsoleteFile)
			if err != nil {
				log.Fatalln(err)
			}
		}

		_, err = printer.Printf("Successfully deleted %d files\n", len(obsoleteFiles))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func flattenUniTree(tree unitree.IconsTree, tsResultPath string) (uniTreeKeys, uniTreeFiles map[string]bool) {
	uniTreeKeys = make(map[string]bool)
	uniTreeFiles = make(map[string]bool)

	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		packagePath := js.BuildPackagePath(tsResultPath, segments)

		uniTreeFiles[js.IndexJs(packagePath)] = true
		uniTreeFiles[js.IndexTs(packagePath)] = true

		for _, icon := range iconSet.Icons {
			uniTreeFiles[js.FileJs(packagePath, icon.Name)] = true
			uniTreeFiles[js.FileTs(packagePath, icon.Name)] = true
			uniTreeKeys[db.CreateIconKey(segments, icon.Name)] = true
		}
	})

	return
}

func findObsoleteKeys(conn rueidis.Client, uniTreeKeys map[string]bool) ([]string, error) {
	obsoleteKeys := make([]string, 0)

	var cursor uint64

scan:
	scanEntry, err := conn.Do(context.Background(), conn.B().Scan().Cursor(cursor).Build()).AsScanEntry()
	if err != nil {
		return obsoleteKeys, err
	}

	cursor = scanEntry.Cursor
	for _, scannedKey := range scanEntry.Elements {
		if _, ok := uniTreeKeys[scannedKey]; !ok {
			obsoleteKeys = append(obsoleteKeys, scannedKey)
		}
	}

	if scanEntry.Cursor == 0 {
		return obsoleteKeys, nil
	}
	goto scan
}

func findObsoleteFiles(manifest js.IcoManifest, uniTreeFiles map[string]bool) ([]string, error) {
	obsoleteFiles := make([]string, 0)
	ignoredPackages := []string{"components", "eslint-config", "typescript-config"}

	err := filepath.WalkDir(manifest.TsResultPath, func(walkPath string, d fs.DirEntry, err error) error {
		for _, ignoredPackage := range ignoredPackages {
			if walkPath == path.Join(manifest.TsResultPath, ignoredPackage) {
				return fs.SkipDir
			}
		}
		if d.Name() == "node_modules" {
			return fs.SkipDir
		}
		if d.IsDir() || d.Name() == "package.json" {
			return nil
		}
		if _, ok := uniTreeFiles[walkPath]; !ok {
			obsoleteFiles = append(obsoleteFiles, walkPath)
		}
		return err
	})

	return obsoleteFiles, err
}
