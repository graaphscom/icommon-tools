package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/db"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/redis/rueidis"
	"log"
)

func main() {
	manifestPath := flag.String("manifest", "", "path to the icons manifest file")
	flag.Parse()

	manifest, err := js.ReadJson[js.IcoManifest](*manifestPath)

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

	responses := conn.DoMulti(ctx, commands...)
	err = nil
	for _, resp := range responses {
		errors.Join(err, resp.Error())
	}
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Vendors manifests have been saved in the db")
}
