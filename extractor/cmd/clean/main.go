package main

import (
	"context"
	"github.com/graaphscom/icommon-tools/extractor/db"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"log"
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

	uniTreeKeys := make(map[string]bool)
	uniTreeFiles := make(map[string]bool)
	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		packagePath := js.BuildPackagePath(manifest.TsResultPath, segments)

		uniTreeFiles[js.IndexJs(packagePath)] = true
		uniTreeFiles[js.IndexTs(packagePath)] = true

		for _, icon := range iconSet.Icons {
			uniTreeFiles[js.FileJs(packagePath, icon.Name)] = true
			uniTreeFiles[js.FileTs(packagePath, icon.Name)] = true
			uniTreeKeys[db.CreateIconKey(segments, icon.Name)] = true
		}
	})

	obsoleteKeys, err := findObsoleteKeys(conn, uniTreeKeys)
	if err != nil {
		log.Fatalln(err)
	}
	var _ = obsoleteKeys

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
