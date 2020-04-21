package constants

import "github.com/cosmos/cosmos-sdk/types/errors"

var IncorrectMessageCode = errors.Register(ModuleName, 101, "IncorrectMessageCode")
