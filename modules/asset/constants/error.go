package constants

import "github.com/cosmos/cosmos-sdk/types/errors"

var UnknownMessageCode = errors.Register(ModuleName, 001, "UnknownMessageCode")
var IncorrectMessageCode = errors.Register(ModuleName, 002, "IncorrectMessageCode")
var UnknownQueryCode = errors.Register(ModuleName, 101, "UnknownQueryCode")
var IncorrectQueryCode = errors.Register(ModuleName, 102, "IncorrectQueryCode")
var EntityNotFoundCode = errors.Register(ModuleName, 103, "EntityNotFoundCode")
