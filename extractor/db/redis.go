package db

import (
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
