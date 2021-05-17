package constants

import "os"

const (
	FlagTimeOut               = "timeout"
	FlagCoinType              = "coinType"
	FlagPStakeHome            = "pStakeHome"
	FlagEthereumEndPoint      = "ethEndPoint"
	FlagTendermintSleepTime   = "tmSleepTime"
	FlagEthereumSleepTime     = "ethSleepTime"
	FlagTendermintStartHeight = "tmStart"
	FlagEthereumStartHeight   = "ethStart"
	FlagDenom                 = "denom"
)

var (
	DefaultTimeout               = "10s"
	DefaultCoinType              = uint32(118)
	DefaultEthereumEndPoint      = "wss://goerli.infura.io/ws/v3/e2549c9ec9764e46a7768cc7619a1939"
	DefaultTendermintSleepTime   = 3000     //ms
	DefaultEthereumSleepTime     = 4500     //ms
	DefaultTendermintStartHeight = int64(0) // 0 will not change the db at start
	DefaultEthereumStartHeight   = int64(0) // 0 will not change the db at start
	DefaultPStakeHome            = os.ExpandEnv("$HOME/.pStake")
	DefaultDenom                 = Denom
)
