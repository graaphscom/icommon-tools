package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/cmd"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	manifestPath := flag.String("manifest", "", "path to the icons manifest file")
	flag.Parse()

	manifest, err := js.ReadJson[js.IcoManifest](*manifestPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Mkdir(manifest.VendorsClonePath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalln(err)
	}

	cloneResults := make(chan error)

	for vendorName, vendorManifest := range manifest.Vendors {
		go cloneVendor(manifest.VendorsClonePath, vendorName, vendorManifest, cloneResults)
	}

	joinedErrors := make([]any, 0, len(manifest.VendorsClonePath))

	for i := 0; i < len(manifest.Vendors); i++ {
		cloneResult := <-cloneResults
		if cloneResult != nil {
			joinedErrors = append(joinedErrors, cloneResult)
		}
	}
	for _, err := range joinedErrors {
		fmt.Println(err)
	}
	if len(joinedErrors) > 0 {
		os.Exit(1)
	}
}

func cloneVendor(destPath, cloneDir string, vendorManifest js.VendorManifest, cloneResult chan<- error) {
	prefix := "[" + cloneDir + "] "
	prefixedOut := cmd.NewPrefixedWriter(prefix, os.Stdout)
	prefixedErr := cmd.NewPrefixedWriter(prefix, os.Stderr)

	cloneCmd := exec.Command("git", "clone", "--no-tags", "-n", "--depth=1", "--filter=tree:0", vendorManifest.RepoUrl, cloneDir)
	cloneCmd.Dir = destPath
	cloneCmd.Stdout = prefixedOut
	cloneCmd.Stderr = prefixedErr
	err := cloneCmd.Run()
	if err != nil {
		cloneResult <- cmd.NewPrefixedError(prefix, err)
		return
	}

	metadataPath := vendorManifest.MetadataPath
	if metadataPath != "" {
		metadataPath = "/" + metadataPath
	}
	repoDir := path.Join(destPath, cloneDir)

	sparseCheckoutCmd := exec.Command("git", "sparse-checkout", "set", "--no-cone", "/"+vendorManifest.IconsPath, metadataPath)
	sparseCheckoutCmd.Dir = repoDir
	sparseCheckoutCmd.Stdout = prefixedOut
	sparseCheckoutCmd.Stderr = prefixedErr
	err = sparseCheckoutCmd.Run()
	if err != nil {
		cloneResult <- cmd.NewPrefixedError(prefix, err)
		return
	}

	checkoutCmd := exec.Command("git", "checkout")
	checkoutCmd.Dir = repoDir
	checkoutCmd.Stdout = prefixedOut
	checkoutCmd.Stderr = prefixedErr
	err = checkoutCmd.Run()
	if err != nil {
		cloneResult <- cmd.NewPrefixedError(prefix, err)
		return
	}

	cloneResult <- nil
}
