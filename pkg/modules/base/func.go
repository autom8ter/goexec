package base

import "github.com/autom8ter/fsctl"

func NewAssetFunc() fsctl.AssetFunc {
	return func(s string) (bytes []byte, e error) {
		return Asset(s)
	}
}

func NewAssetDirFunc() fsctl.AssetDirFunc {
	return func(s string) (strings []string, e error) {
		return AssetDir(s)
	}
}
