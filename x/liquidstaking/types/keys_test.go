package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type keysTestSuite struct {
	suite.Suite
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(keysTestSuite))
}

