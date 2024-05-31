package metadata

import (
	"errors"
	"github.com/graaphscom/icommon-tools/extractor/json"
	"os"
	"path"
	"strings"
)

func (S Store) GetFluentui(asset string) (Fluentui, error) {
	if S.metadata.fluentui == nil {
		S.metadata.fluentui = make(map[string]Fluentui)
	}

	if v, ok := S.metadata.fluentui[asset]; ok {
		return v, nil
	}

	result, err := json.ReadJson[Fluentui](
		path.Join(S.manifest.VendorsClonePath, "fluentui", S.manifest.Vendors["fluentui"].IconsPath, asset, "metadata.json"),
	)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Fluentui{}, err
	}

	S.metadata.fluentui[asset] = result

	return result, nil
}

func (S Store) GetFontawesome() (Fontawesome, error) {
	return singleFile[Fontawesome](&S.metadata.fontawesome, S.manifest, "fontawesome")
}

func (S Store) GetOcticons() (Octicons, error) {
	return singleFile[Octicons](&S.metadata.octicons, S.manifest, "octicons")
}

func (S Store) GetRemixicon() (Remixicon, error) {
	if S.metadata.remixicon != nil {
		return S.metadata.remixicon, nil
	}

	result, err := json.ReadJsonOmitProp[Remixicon](path.Join(S.manifest.VendorsClonePath, "remixicon", S.manifest.Vendors["remixicon"].MetadataPath), "_comment")

	if err != nil {
		var empty Remixicon
		return empty, err
	}

	S.metadata.remixicon = result

	return result, nil
}

func (S Store) GetUnicons(variant string) (UniconsQuickAccess, error) {
	if S.metadata.unicons == nil {
		S.metadata.unicons = make(map[string]UniconsQuickAccess)
	}

	if v, ok := S.metadata.unicons[variant]; ok {
		return v, nil
	}

	jsonContents, err := json.ReadJson[unicons](
		path.Join(S.manifest.VendorsClonePath, "unicons", S.manifest.Vendors["unicons"].MetadataPath, strings.Join([]string{variant, ".json"}, "")),
	)

	if err != nil {
		return nil, err
	}

	result := make(map[string][]string, len(jsonContents))
	for _, entry := range jsonContents {
		result[entry.Name] = entry.Tags
	}

	S.metadata.unicons[variant] = result

	return result, nil
}

func singleFile[T Fontawesome | Octicons](cacheEntry *T, manifest json.IcoManifest, vendor string) (T, error) {
	if *cacheEntry != nil {
		return *cacheEntry, nil
	}

	result, err := json.ReadJson[T](path.Join(manifest.VendorsClonePath, vendor, manifest.Vendors[vendor].MetadataPath))

	if err != nil {
		var empty T
		return empty, err
	}

	*cacheEntry = result

	return result, nil
}

func NewStore(manifest json.IcoManifest) Store {
	return Store{manifest: manifest, metadata: &metadata{}}
}

type Store struct {
	metadata *metadata
	manifest json.IcoManifest
}

type metadata struct {
	fluentui    map[string]Fluentui
	fontawesome Fontawesome
	octicons    Octicons
	remixicon   Remixicon
	unicons     map[string]UniconsQuickAccess
}

type Fluentui struct {
	Metaphor []string
}

type Fontawesome map[string]struct{ Search struct{ Terms []string } }

type Octicons map[string][]string

type Remixicon map[string]map[string]string

type unicons []struct {
	Name string
	Tags []string
}

type UniconsQuickAccess map[string][]string
