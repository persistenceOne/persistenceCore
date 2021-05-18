package ethereum

import (
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum/contracts"
)

func init() {
	contracts.LiquidStaking.SetABI(constants.LiquidStakingABI)
}
