package main

import (
	"context"
	"github.com/graaphscom/icommon-tools/extractor/json"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"log"
	"strings"
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

	existingKeys := make(map[string]bool)
	tree.MustTraverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		for _, icon := range iconSet.Icons {
			existingKeys[strings.Join(append(segments, icon.Name), ":")] = true
		}
	})

	obsoleteKeys, err := findObsoleteKeys(conn, existingKeys)
	if err != nil {
		log.Fatalln(err)
	}
	var _ = obsoleteKeys

}

func findObsoleteKeys(conn rueidis.Client, existingKeys map[string]bool) ([]string, error) {
	obsoleteKeys := make([]string, 0)

	var cursor uint64

scan:
	scanEntry, err := conn.Do(context.Background(), conn.B().Scan().Cursor(cursor).Build()).AsScanEntry()
	if err != nil {
		return obsoleteKeys, err
	}

	cursor = scanEntry.Cursor
	for _, scannedKey := range scanEntry.Elements {
		if _, ok := existingKeys[scannedKey]; !ok {
			obsoleteKeys = append(obsoleteKeys, scannedKey)
		}
	}

	if scanEntry.Cursor == 0 {
		return obsoleteKeys, nil
	}
	goto scan
}
