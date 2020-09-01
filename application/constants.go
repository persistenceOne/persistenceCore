package application

import "os"

const Name = "AssetMantle"

var DefaultClientHome = os.ExpandEnv("$HOME/.assetClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.assetNode")
