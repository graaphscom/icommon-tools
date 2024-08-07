package unitree

import (
	"fmt"
	"github.com/graaphscom/icommon-tools/extractor/js"
	"github.com/graaphscom/icommon-tools/extractor/metadata"
	"github.com/graaphscom/icommon-tools/extractor/strcase"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
)

var treeBuilders = map[string]treeBuilder{
	"boxicons": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			tagsExtractor: func(_ metadata.Store, rawRootName, rawName string) (IconTags, error) {
				return IconTags{
					Search: []string{
						strings.ReplaceAll(strings.TrimSuffix(firstHyphenRegexp.ReplaceAllString(rawName, ""), ".svg"), "-", " "),
					},
					Visual: []string{rawRootName},
				}, nil
			},
		},
	},
	"bytesize": iconsTreeBuilder{
		iconNameConverter: iconNameKebabCaseConverter,
		tagsExtractor:     tagsExtractorKebabCase,
	},
	"fluentui": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: func(in string) string {
				return fixVarNameFirstChar(strcase.ToCamel(strings.TrimPrefix(strings.TrimSuffix(in, ".svg"), "ic_fluent"), strcase.SnakeRegexp))
			},
			treeNameConverter: treeNameSpaceConverter,
			srcSuffix:         "SVG",
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetFluentui(rawRootName)
				matches := fluentuiTagsRegexp.FindStringSubmatch(rawName)
				if os.IsNotExist(err) {
					return IconTags{
						Search: []string{strings.ReplaceAll(matches[1], "_", " ")},
						Visual: matches[2:],
					}, nil
				}
				if err != nil {
					return IconTags{}, err
				}
				return IconTags{
					Search: append(m.Metaphor, strings.ReplaceAll(matches[1], "_", " ")),
					Visual: matches[2:],
				}, nil
			},
		},
	},
	"fontawesome": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetFontawesome()
				if err != nil {
					return IconTags{}, err
				}
				rawNameTrimmed := strings.TrimSuffix(rawName, ".svg")
				if icoMeta, ok := m[rawNameTrimmed]; ok {
					return IconTags{
						Search: append(icoMeta.Search.Terms, strings.ReplaceAll(rawNameTrimmed, "-", " ")),
						Visual: []string{rawRootName},
					}, nil
				}
				return IconTags{}, nil
			},
		},
	},
	"material": categoriesTreeBuilder{
		iconsTreeBuilder: materialIconsTreeBuilder{},
	},
	"octicons": iconsTreeBuilder{
		iconNameConverter: iconNameKebabCaseConverter,
		tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
			m, err := metadata.GetOcticons()
			if err != nil {
				return IconTags{}, err
			}
			matches := octiconsTagsRegexp.FindStringSubmatch(rawName)
			visualTags := slices.DeleteFunc(matches[2:], func(s string) bool {
				return s == ""
			})
			for idx, visualTag := range visualTags {
				visualTags[idx] = strings.Trim(visualTag, "-")
			}
			searchTags := []string{strings.ReplaceAll(matches[1], "-", " ")}
			if _, ok := m[matches[1]]; ok {
				searchTags = append(searchTags, m[matches[1]]...)
			}
			return IconTags{Search: searchTags, Visual: visualTags}, nil
		},
	},
	"radixui": iconsTreeBuilder{
		iconNameConverter: iconNameKebabCaseConverter,
		tagsExtractor:     tagsExtractorKebabCase,
	},
	"remixicon": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			treeNameConverter: treeNameSpaceConverter,
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetRemixicon()
				if err != nil {
					return IconTags{}, err
				}
				matches := remixiconTagsRegexp.FindStringSubmatch(rawName)

				searchTags := []string{strings.ToLower(rawRootName), strings.ReplaceAll(matches[1], "-", " ")}
				if _, ok := m[rawRootName][matches[1]]; ok {
					searchTags = append(searchTags, strings.Split(m[rawRootName][matches[1]], ",")...)
				}
				return IconTags{
					Search: searchTags,
					Visual: []string{strings.Trim(matches[2], "-")},
				}, nil
			},
		},
	},
	"unicons": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetUnicons(rawRootName)
				var _ = m
				if err != nil {
					return IconTags{}, err
				}

				var searchTags []string
				if tags, ok := m[strings.TrimSuffix(rawName, ".svg")]; ok {
					for _, tag := range tags {
						searchTags = append(searchTags, strings.ReplaceAll(tag, "-", " "))
					}
				} else {
					searchTags = []string{strings.ReplaceAll(strings.TrimSuffix(rawName, ".svg"), "-", " ")}
				}

				return IconTags{Search: searchTags, Visual: []string{rawRootName}}, nil
			},
		},
	},
}

func (b materialIconsTreeBuilder) buildTree(_ metadata.Store, src, rootName string) (IconsTree, error) {
	srcEntries, err := os.ReadDir(src)

	if err != nil {
		return IconsTree{}, err
	}

	icons := make([]Icon, 0, len(srcEntries))

	for _, srcEntry := range srcEntries {
		if !srcEntry.IsDir() {
			continue
		}

		subSrcEntries, err := os.ReadDir(path.Join(src, srcEntry.Name()))
		if err != nil {
			return IconsTree{}, err
		}

		searchTag := strings.ReplaceAll(srcEntry.Name(), "_", " ")
		snakeCaseConverted := iconNameSnakeCaseConverter(srcEntry.Name())
		for _, subSrcEntry := range subSrcEntries {
			if !subSrcEntry.IsDir() {
				continue
			}

			subSubSrcEntries, err := os.ReadDir(path.Join(src, srcEntry.Name(), subSrcEntry.Name()))
			if err != nil {
				return IconsTree{}, err
			}

			style := strings.TrimPrefix(subSrcEntry.Name(), "materialicons")
			var styleFirstUpper string
			if len(style) > 0 {
				styleFirstUpper = strings.ToUpper(string(style[0])) + style[1:]
			}
			for _, subSubSrcEntry := range subSubSrcEntries {
				icons = append(icons, Icon{
					Name:    snakeCaseConverted + styleFirstUpper + strings.TrimSuffix(subSubSrcEntry.Name(), ".svg"),
					SrcFile: path.Join(src, srcEntry.Name(), subSrcEntry.Name(), subSubSrcEntry.Name()),
					Tags: IconTags{
						Search: []string{searchTag},
						Visual: []string{style, strings.TrimSuffix(subSubSrcEntry.Name(), "px.svg")},
					},
				})
			}
		}
	}

	for i := 0; i < len(icons); i++ {
		for j := i + 1; j < len(icons); j++ {
			if icons[i].Name == icons[j].Name {
				icons = append(icons[:i], icons[i+1:]...)
			}
		}
	}

	return IconsTree{
		Name: rootName,
		IconSet: &IconSet{
			Icons: icons,
		},
	}, nil
}

type materialIconsTreeBuilder struct{}

func iconNameKebabCaseConverter(in string) string {
	return fixReservedWords(fixVarNameFirstChar(strcase.ToCamel(strings.TrimSuffix(in, ".svg"), strcase.KebabRegexp)))
}

func iconNameSnakeCaseConverter(in string) string {
	return fixVarNameFirstChar(strcase.ToCamel(strings.TrimSuffix(in, ".svg"), strcase.SnakeRegexp))
}

func treeNameSpaceConverter(in string) string {
	return strcase.ToCamel(in, strcase.SpaceRegexp)
}

func tagsExtractorKebabCase(_ metadata.Store, _, rawName string) (IconTags, error) {
	return IconTags{Search: []string{strings.ReplaceAll(strings.TrimSuffix(rawName, ".svg"), "-", " ")}}, nil
}

var firstHyphenRegexp, _ = regexp.Compile(`^\w*-`)
var fluentuiTagsRegexp, _ = regexp.Compile(`ic_fluent_(.*?)_(\d*)?_(regular|filled|light)?(_ltr)?(_rtl)?.svg`)
var octiconsTagsRegexp, _ = regexp.Compile(`(.*?)(-circle)?(-fill)?(-\d*)?.svg`)
var remixiconTagsRegexp, _ = regexp.Compile(`(.*)(-line|-fill)?.svg`)

func fixVarNameFirstChar(varName string) string {
	notAllowedFirstCharRegexp, _ := regexp.Compile("^[^A-Za-z_]")
	if notAllowedFirstCharRegexp.MatchString(varName) {
		return fmt.Sprintf("__%s", varName)
	}

	return varName
}

func fixReservedWords(varName string) string {
	if _, ok := js.ReservedWords[varName]; ok {
		return "__" + varName
	}

	return varName
}
