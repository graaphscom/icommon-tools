package main

import (
	"github.com/graaphscom/icommon/extractor/json"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	manifest, err := json.ReadJson[json.IcoManifest]("testdata/ico_manifest_downloads.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Mkdir(manifest.VendorsClonePath, 0750)
	if err != nil {
		log.Fatalln(err)
	}

	cloneResults := make(chan error)

	for vendorName, vendorManifest := range manifest.Vendors {
		go cloneVendor(manifest.VendorsClonePath, vendorName, vendorManifest, cloneResults)
	}

	for i := 0; i < len(manifest.Vendors); i++ {
		cloneErr := <-cloneResults
		if cloneErr != nil {
			log.Fatalln(cloneErr)
		}
	}
}

func cloneVendor(destPath, cloneDir string, vendorManifest json.VendorManifest, cloneResult chan<- error) {
	cloneCmd := exec.Command("git", "clone", "-n", "--depth=1", "--filter=tree:0", vendorManifest.RepoUrl, cloneDir)
	cloneCmd.Dir = destPath
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	err := cloneCmd.Run()
	if err != nil {
		cloneResult <- err
		return
	}

	metadataPath := vendorManifest.MetadataPath
	if metadataPath != "" {
		metadataPath = "/" + metadataPath
	}
	repoDir := path.Join(destPath, cloneDir)

	sparseCheckoutCmd := exec.Command("git", "sparse-checkout", "set", "--no-cone", "/"+vendorManifest.IconsPath, metadataPath)
	sparseCheckoutCmd.Dir = repoDir
	sparseCheckoutCmd.Stdout = os.Stdout
	sparseCheckoutCmd.Stderr = os.Stderr
	err = sparseCheckoutCmd.Run()
	if err != nil {
		cloneResult <- err
		return
	}

	checkoutCmd := exec.Command("git", "checkout")
	checkoutCmd.Dir = repoDir
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr
	err = checkoutCmd.Run()
	if err != nil {
		cloneResult <- err
		return
	}

	cloneResult <- nil
}
