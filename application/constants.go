/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import "os"

const Name = "AssetMantle"

var DefaultClientHome = os.ExpandEnv("$HOME/.assetClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.assetNode")
