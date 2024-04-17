package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	flag.Usage = func() {}
}

func cloneVendor() error {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	cloneCmd := exec.Command("git", "clone", "-n", "--depth=1", "--filter=tree:0", "git@github.com:google/material-design-icons.git")
	cloneCmd.Dir = home
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	err = cloneCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	sparseCheckoutCmd := exec.Command("git", "sparse-checkout", "set", "--no-cone", "src")
	repoDir := path.Join(home, "material-design-icons")
	sparseCheckoutCmd.Dir = repoDir
	sparseCheckoutCmd.Stdout = os.Stdout
	sparseCheckoutCmd.Stderr = os.Stderr
	err = sparseCheckoutCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	checkoutCmd := exec.Command("git", "checkout")
	checkoutCmd.Dir = repoDir
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr
	err = checkoutCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Cloned material-design-icons.")

	return nil
}
