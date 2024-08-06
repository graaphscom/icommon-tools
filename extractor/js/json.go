package js

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type IcoManifest struct {
	RedisUrl         string
	VendorsClonePath string
	Vendors          map[string]VendorManifest
	TsResultPath     string
}

type VendorManifest struct {
	RepoUrl      string
	IconsPath    string
	MetadataPath string
	License      string
	LicenseUrl   string
	Homepage     string
	Funding      string
}

func ReadJsonOmitProp[T any](path, prop string) (T, error) {
	contents, err := os.ReadFile(path)

	var result T

	if err != nil {
		err = fmt.Errorf("%w, %s", err, path)
		return result, err
	}

	propNodeRegexp, err := regexp.Compile(`"` + prop + `":.*,`)
	if err != nil {
		return result, err
	}
	noProp := propNodeRegexp.ReplaceAll(contents, []byte(""))

	err = json.Unmarshal(noProp, &result)
	if err != nil {
		err = fmt.Errorf("%w, %s", err, path)
	}

	return result, err
}

func ReadJson[T any](path string) (T, error) {
	contents, err := os.ReadFile(path)

	var result T

	if err != nil {
		err = fmt.Errorf("%w, %s", err, path)
		return result, err
	}

	err = json.Unmarshal(contents, &result)
	if err != nil {
		err = fmt.Errorf("%w, %s", err, path)
	}

	return result, err
}
