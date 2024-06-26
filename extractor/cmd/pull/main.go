package main

import (
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/cmd"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	manifest, err := js.ReadJson[js.IcoManifest]("testdata/ico_manifest_downloads.json")
	if err != nil {
		log.Fatalln(err)
	}

	pullResults := make(chan error)

	for vendorName := range manifest.Vendors {
		go pullVendor(manifest.VendorsClonePath, vendorName, pullResults)
	}

	joinedErrors := make([]any, 0, len(manifest.VendorsClonePath))

	for i := 0; i < len(manifest.Vendors); i++ {
		pullResult := <-pullResults
		if pullResult != nil {
			joinedErrors = append(joinedErrors, pullResult)
		}
	}
	for _, err := range joinedErrors {
		fmt.Println(err)
	}
	if len(joinedErrors) > 0 {
		os.Exit(1)
	}
}

func pullVendor(destPath, cloneDir string, cloneResult chan<- error) {
	pullCmd := exec.Command("git", "pull", "--depth=1", "--rebase=true")
	pullCmd.Dir = path.Join(destPath, cloneDir)
	prefix := "[" + cloneDir + "] "
	pullCmd.Stdout = cmd.NewPrefixedWriter(prefix, os.Stdout)
	pullCmd.Stderr = cmd.NewPrefixedWriter(prefix, os.Stderr)
	err := pullCmd.Run()
	if err != nil {
		cloneResult <- cmd.NewPrefixedError(prefix, err)
		return
	}

	cloneResult <- nil
}
