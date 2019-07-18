package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tendermintTypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	freeTokensPerAccount    = sdk.TokensFromTendermintPower(150)
	defaultBondDenomination = sdk.DefaultBondDenom
)

type GenesisState struct {
	Accounts            []GenesisAccount          `json:"accounts"`
	AuthData            auth.GenesisState         `json:"auth"`
	BankData            bank.GenesisState         `json:"bank"`
	StakingData         staking.GenesisState      `json:"staking"`
	MintData            mint.GenesisState         `json:"mint"`
	DistributionData    distribution.GenesisState `json:"distribution"`
	GovernmentData      gov.GenesisState          `json:"government"`
	CrisisData          crisis.GenesisState       `json:"crisis"`
	SlashingData        slashing.GenesisState     `json:"slashing"`
	GenesisTransactions []json.RawMessage         `json:"genesisTransactions"`
}

func NewGenesisState(
	accounts []GenesisAccount,
	authData auth.GenesisState,
	bankData bank.GenesisState,
	stakingData staking.GenesisState,
	mintData mint.GenesisState,
	distributionData distribution.GenesisState,
	governmentData gov.GenesisState,
	crisisData crisis.GenesisState,
	slashingData slashing.GenesisState,
) GenesisState {
	return GenesisState{
		Accounts:         accounts,
		AuthData:         authData,
		BankData:         bankData,
		StakingData:      stakingData,
		MintData:         mintData,
		DistributionData: distributionData,
		GovernmentData:   governmentData,
		CrisisData:       crisisData,
		SlashingData:     slashingData,
	}
}
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts:            nil,
		AuthData:            auth.DefaultGenesisState(),
		BankData:            bank.DefaultGenesisState(),
		StakingData:         staking.DefaultGenesisState(),
		MintData:            mint.DefaultGenesisState(),
		DistributionData:    distribution.DefaultGenesisState(),
		GovernmentData:      gov.DefaultGenesisState(),
		CrisisData:          crisis.DefaultGenesisState(),
		SlashingData:        slashing.DefaultGenesisState(),
		GenesisTransactions: nil,
	}
}
func ValidateGenesisState(genesisState GenesisState) error {
	if err := validateGenesisStateAccounts(genesisState.Accounts); err != nil {
		return err
	}
	if len(genesisState.GenesisTransactions) > 0 {
		return nil
	}

	if err := auth.ValidateGenesis(genesisState.AuthData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
		return err
	}
	if err := distribution.ValidateGenesis(genesisState.DistributionData); err != nil {
		return err
	}
	if err := gov.ValidateGenesis(genesisState.GovernmentData); err != nil {
		return err
	}
	if err := crisis.ValidateGenesis(genesisState.CrisisData); err != nil {
		return err
	}

	return slashing.ValidateGenesis(genesisState.SlashingData)
}
func validateGenesisStateAccounts(accs []GenesisAccount) error {
	addrMap := make(map[string]bool, len(accs))
	for _, acc := range accs {
		addrStr := acc.Address.String()

		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		if !acc.OriginalVesting.IsZero() {
			if acc.EndTime == 0 {
				return fmt.Errorf("missing end time for vesting account; address: %s", addrStr)
			}

			if acc.StartTime >= acc.EndTime {
				return fmt.Errorf(
					"vesting start time must before end time; address: %s, start: %s, end: %s",
					addrStr,
					time.Unix(acc.StartTime, 0).UTC().Format(time.RFC3339),
					time.Unix(acc.EndTime, 0).UTC().Format(time.RFC3339),
				)
			}
		}

		addrMap[addrStr] = true
	}

	return nil
}
func (genesisState GenesisState) Sanitize() {
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})

	for _, acc := range genesisState.Accounts {
		acc.Coins = acc.Coins.Sort()
	}
}

type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`

	OriginalVesting  sdk.Coins `json:"original_vesting"`
	DelegatedFree    sdk.Coins `json:"delegated_free"`
	DelegatedVesting sdk.Coins `json:"delegated_vesting"`
	StartTime        int64     `json:"start_time"`
	EndTime          int64     `json:"end_time"`
}

func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}
func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	vacc, ok := acc.(auth.VestingAccount)
	if ok {
		gacc.OriginalVesting = vacc.GetOriginalVesting()
		gacc.DelegatedFree = vacc.GetDelegatedFree()
		gacc.DelegatedVesting = vacc.GetDelegatedVesting()
		gacc.StartTime = vacc.GetStartTime()
		gacc.EndTime = vacc.GetEndTime()
	}

	return gacc
}
func NewDefaultGenesisAccount(addr sdk.AccAddress) GenesisAccount {
	accAuth := auth.NewBaseAccountWithAddress(addr)
	coins := sdk.Coins{
		sdk.NewCoin("footoken", sdk.NewInt(1000)),
		sdk.NewCoin(defaultBondDenomination, freeTokensPerAccount),
	}

	coins.Sort()

	accAuth.Coins = coins
	return NewGenesisAccount(&accAuth)
}
func (genesisAccount *GenesisAccount) ToAccount() auth.Account {
	baseAccount := &auth.BaseAccount{
		Address:       genesisAccount.Address,
		Coins:         genesisAccount.Coins.Sort(),
		AccountNumber: genesisAccount.AccountNumber,
		Sequence:      genesisAccount.Sequence,
	}

	if !genesisAccount.OriginalVesting.IsZero() {
		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:      baseAccount,
			OriginalVesting:  genesisAccount.OriginalVesting,
			DelegatedFree:    genesisAccount.DelegatedFree,
			DelegatedVesting: genesisAccount.DelegatedVesting,
			EndTime:          genesisAccount.EndTime,
		}

		if genesisAccount.StartTime != 0 && genesisAccount.EndTime != 0 {
			return &auth.ContinuousVestingAccount{
				BaseVestingAccount: baseVestingAcc,
				StartTime:          genesisAccount.StartTime,
			}
		} else if genesisAccount.EndTime != 0 {
			return &auth.DelayedVestingAccount{
				BaseVestingAccount: baseVestingAcc,
			}
		} else {
			panic(fmt.Sprintf("invalid genesis vesting account: %+v", genesisAccount))
		}
	}

	return baseAccount
}

func CommitHubApplicationGenesisState(cdc *codec.Codec, genesisDoc tendermintTypes.GenesisDoc, applicationGenesisTransactions []json.RawMessage) (
	genesisState GenesisState, err error) {

	if err = cdc.UnmarshalJSON(genesisDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}

	if len(applicationGenesisTransactions) == 0 {
		return genesisState, errors.New("there must be at least one genesis tx")
	}

	stakingData := genesisState.StakingData
	for i, applicationGenesisTransaction := range applicationGenesisTransactions {
		var tx auth.StdTx
		if err := cdc.UnmarshalJSON(applicationGenesisTransaction, &tx); err != nil {
			return genesisState, err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return genesisState, errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		if _, ok := msgs[0].(staking.MsgCreateValidator); !ok {
			return genesisState, fmt.Errorf(
				"genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}

	for _, acc := range genesisState.Accounts {
		for _, coin := range acc.Coins {
			if coin.Denom == genesisState.StakingData.Params.BondDenom {
				stakingData.Pool.NotBondedTokens = stakingData.Pool.NotBondedTokens.
					Add(coin.Amount)
			}
		}
	}

	genesisState.StakingData = stakingData
	genesisState.GenesisTransactions = applicationGenesisTransactions

	return genesisState, nil
}
func CommitHubApplicationGenesiStateJSON(cdc *codec.Codec, genDoc tendermintTypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {
	genesisState, err := CommitHubApplicationGenesisState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(cdc, genesisState)
}

func CollectStandardTransacrions(cdc *codec.Codec, moniker string, genTxsDir string, genDoc tendermintTypes.GenesisDoc) (
	appGenTxs []auth.StdTx, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return appGenTxs, persistentPeers, err
	}

	// prepare a map of all accounts in genesis state to then validate
	// against the validators addresses
	var appState GenesisState
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
		return appGenTxs, persistentPeers, err
	}

	addrMap := make(map[string]GenesisAccount, len(appState.Accounts))
	for i := 0; i < len(appState.Accounts); i++ {
		acc := appState.Accounts[i]
		addrMap[acc.Address.String()] = acc
	}

	// addresses and IPs (and port) validator server info
	var addressesIPs []string

	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			return appGenTxs, persistentPeers, err
		}
		var genStdTx auth.StdTx
		if err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx); err != nil {
			return appGenTxs, persistentPeers, err
		}
		appGenTxs = append(appGenTxs, genStdTx)

		// the memo flag is used to store
		// the ip and node-id, for example this may be:
		// "528fd3df22b31f4969b05652bfe8f0fe921321d5@192.168.2.37:26656"
		nodeAddrIP := genStdTx.GetMemo()
		if len(nodeAddrIP) == 0 {
			return appGenTxs, persistentPeers, fmt.Errorf(
				"couldn't find node's address and IP in %s", fo.Name())
		}

		// genesis transactions must be single-message
		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {

			return appGenTxs, persistentPeers, errors.New(
				"each genesis transaction must provide a single genesis message")
		}

		msg := msgs[0].(staking.MsgCreateValidator)
		// validate delegator and validator addresses and funds against the accounts in the state
		delAddr := msg.DelegatorAddress.String()
		valAddr := sdk.AccAddress(msg.ValidatorAddress).String()

		delAcc, delOk := addrMap[delAddr]
		_, valOk := addrMap[valAddr]

		accsNotInGenesis := []string{}
		if !delOk {
			accsNotInGenesis = append(accsNotInGenesis, delAddr)
		}
		if !valOk {
			accsNotInGenesis = append(accsNotInGenesis, valAddr)
		}
		if len(accsNotInGenesis) != 0 {
			return appGenTxs, persistentPeers, fmt.Errorf(
				"account(s) %v not in genesis.json: %+v", strings.Join(accsNotInGenesis, " "), addrMap)
		}

		if delAcc.Coins.AmountOf(msg.Value.Denom).LT(msg.Value.Amount) {
			return appGenTxs, persistentPeers, fmt.Errorf(
				"insufficient fund for delegation %v: %v < %v",
				delAcc.Address, delAcc.Coins.AmountOf(msg.Value.Denom), msg.Value.Amount,
			)
		}

		// exclude itself from persistent peers
		if msg.Description.Moniker != moniker {
			addressesIPs = append(addressesIPs, nodeAddrIP)
		}
	}

	sort.Strings(addressesIPs)
	persistentPeers = strings.Join(addressesIPs, ",")

	return appGenTxs, persistentPeers, nil
}
