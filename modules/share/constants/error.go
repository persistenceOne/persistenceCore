package constants

import "github.com/cosmos/cosmos-sdk/types/errors"

var IncorrectMessageCode = errors.Register(ModuleName, 101, "IncorrectMessageCode")
var ShareNotFoundCode = errors.Register(ModuleName, 201, "ShareNotFoundCode")
var IncorrectQueryCode = errors.Register(ModuleName, 301, "IncorrectQueryCode")
