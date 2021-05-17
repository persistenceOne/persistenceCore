package ethereum

import (
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum/contracts"
)

func init() {
	contracts.STokens.SetABI(constants.STokensABI)
	contracts.LiquidStaking.SetABI(constants.LiquidStakingABI)
}
