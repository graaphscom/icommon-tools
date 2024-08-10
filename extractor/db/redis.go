package db

import (
	"context"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/graaphscom/icommon-tools/extractor/unitree"
	"github.com/redis/rueidis"
	"strings"
)

func CreateIconEntry(builder rueidis.Builder, segments []string, icon unitree.Icon) rueidis.Completed {
	return builder.Hset().
		Key(CreateIconKey(segments, icon.Name)).
		FieldValue().
		FieldValue("searchTags", strings.Join(icon.Tags.Search, ",")).
		FieldValue("visualTags", strings.Join(icon.Tags.Visual, ",")).
		Build()
}

func CreateManifestEntry(builder rueidis.Builder, vendorName string, manifest js.VendorManifest) rueidis.Completed {
	return builder.Hset().
		Key("icommon-manifest:"+vendorName).
		FieldValue().
		FieldValue("name", vendorName).
		FieldValue("funding", manifest.Funding).
		FieldValue("homepage", manifest.Homepage).
		FieldValue("license", manifest.License).
		FieldValue("licenseUrl", manifest.LicenseUrl).
		Build()
}

func CreateIconKey(segments []string, iconName string) string {
	return strings.Join(append(segments, iconName), ":")
}

func CreateSearchIndex(builder rueidis.Builder) rueidis.Completed {
	return builder.FtCreate().
		Index("icommon").
		Prefix(1).Prefix("icommon:").
		Schema().
		FieldName("searchTags").Text().
		FieldName("visualTags").Tag().
		Build()
}

func FindObsoleteKeys(conn rueidis.Client, uniTreeKeys map[string]bool) ([]string, error) {
	obsoleteKeys := make([]string, 0)

	var cursor uint64

scan:
	scanEntry, err := conn.Do(context.Background(), conn.B().Scan().Cursor(cursor).Match("icommon:*").Build()).AsScanEntry()
	if err != nil {
		return obsoleteKeys, err
	}

	cursor = scanEntry.Cursor
	for _, scannedKey := range scanEntry.Elements {
		if _, ok := uniTreeKeys[scannedKey]; !ok {
			obsoleteKeys = append(obsoleteKeys, scannedKey)
		}
	}

	if scanEntry.Cursor == 0 {
		return obsoleteKeys, nil
	}
	goto scan
}
