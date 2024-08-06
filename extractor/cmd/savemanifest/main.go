package main

import (
	"context"
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/db"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/redis/rueidis"
	"log"
)

func main() {
	manifest, err := js.ReadJson[js.IcoManifest]("testdata/ico_manifest_downloads.json")

	opt, err := rueidis.ParseURL(manifest.RedisUrl)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := rueidis.NewClient(opt)
	if err != nil {
		log.Fatalln(err)
	}

	var commands []rueidis.Completed

	for vendorName, vendorManifest := range manifest.Vendors {
		commands = append(commands, db.CreateManifestEntry(conn.B(), vendorName, vendorManifest))
	}

	ctx := context.Background()

	conn.DoMulti(ctx, commands...)

	fmt.Println("Vendors manifests have been saved in the db")
}
