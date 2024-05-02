package unitree

import (
	"errors"
	"github.com/graaphscom/icommon/extractor/metadata"
)

type treeBuilder interface {
	buildTree(metadata metadata.Store, src, rootName string) (IconsTree, error)
}

type IconsTree struct {
	Name     string
	SubTrees *[]IconsTree
	IconSet  *IconSet
}

func (tree IconsTree) Traverse(
	segments []string,
	iconSetFn func(segments []string, iconSet IconSet) error,
	subTreeFn func(segments []string) error,
) error {
	var err error

	if tree.IconSet != nil {
		err = errors.Join(err, iconSetFn(append(segments, tree.Name), *tree.IconSet))
	}

	if tree.SubTrees != nil {
		for _, subTree := range *tree.SubTrees {
			appendedSegment := append(segments, tree.Name)
			if subTreeFn != nil {
				err = errors.Join(err, subTreeFn(appendedSegment))
			}
			err = errors.Join(err, subTree.Traverse(appendedSegment, iconSetFn, subTreeFn))
		}
	}

	return err
}

func (tree IconsTree) MustTraverse(segments []string, iconSetFn func(segments []string, iconSet IconSet)) {
	if tree.IconSet != nil {
		iconSetFn(append(segments, tree.Name), *tree.IconSet)
	}
	if tree.SubTrees != nil {
		for _, subTree := range *tree.SubTrees {
			appendedSegment := append(segments, tree.Name)
			subTree.MustTraverse(appendedSegment, iconSetFn)
		}
	}
}

type IconSet struct {
	Icons []Icon
}

type Icon struct {
	Name    string
	SrcFile string
	Tags    IconTags
}

type IconTags struct {
	Search []string
	Visual []string
}
