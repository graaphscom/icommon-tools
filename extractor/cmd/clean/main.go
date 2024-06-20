package main

import (
	"github.com/graaphscom/icommon-tools/extractor/json"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"log"
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

}
