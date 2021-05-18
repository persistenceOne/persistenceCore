// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AbiABI is the input ABI used to generate the binding from.
const AbiABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"GenerateUTokens\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"SetUTokensContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toAtomAddress\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"WithdrawUTokens\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BRIDGE_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"generateUTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"uAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"bridgeAdminAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pauserAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"uAddress\",\"type\":\"address\"}],\"name\":\"setUTokensContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"toAtomAddress\",\"type\":\"string\"}],\"name\":\"withdrawUTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AbiBin is the compiled bytecode used for deploying new contracts.
var AbiBin = "0x608060405234801561001057600080fd5b506120f3806100206000396000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c80635c975abb116100a2578063a217fddf11610071578063a217fddf146104a9578063c0c53b8b146104c7578063ca15c8731461054b578063d547741f1461058d578063e63ab1e9146105db5761010b565b80635c975abb146103a35780638456cb59146103c35780639010d07c146103e357806391d14854146104455761010b565b80632f2ff15d116100de5780632f2ff15d1461020257806336568abe146102505780633f4ba83a1461029e578063428bee9e146102be5761010b565b8063118c38c71461011057806321bdf9e51461012e57806322dd9bc314610172578063248a9ca3146101c0575b600080fd5b6101186105f9565b6040518082815260200191505060405180910390f35b6101706004803603602081101561014457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061061d565b005b6101be6004803603604081101561018857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061070d565b005b6101ec600480360360208110156101d657600080fd5b8101908080359060200190929190505050610990565b6040518082815260200191505060405180910390f35b61024e6004803603604081101561021857600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506109b0565b005b61029c6004803603604081101561026657600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a3a565b005b6102a6610ad3565b60405180821515815260200191505060405180910390f35b6103a1600480360360608110156102d457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561031b57600080fd5b82018360208201111561032d57600080fd5b8035906020019184600183028401116401000000008311171561034f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610b6a565b005b6103ab610f82565b60405180821515815260200191505060405180910390f35b6103cb610f99565b60405180821515815260200191505060405180910390f35b610419600480360360408110156103f957600080fd5b810190808035906020019092919080359060200190929190505050611030565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6104916004803603604081101561045b57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611062565b60405180821515815260200191505060405180910390f35b6104b1611094565b6040518082815260200191505060405180910390f35b610549600480360360608110156104dd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061109b565b005b6105776004803603602081101561056157600080fd5b810190808035906020019092919050505061121d565b6040518082815260200191505060405180910390f35b6105d9600480360360408110156105a357600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611244565b005b6105e36112ce565b6040518082815260200191505060405180910390f35b7f751b795d24b92e3d92d1d0d8f2885f4e9c9c269da350af36ae6b49069babf4bf81565b6106316000801b61062c6112f2565b611062565b610686576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526038815260200180611f7c6038913960400191505060405180910390fd5b80609760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167ff745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c60405160405180910390a250565b610715610f82565b15610788576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b600081116107e1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526037815260200180611f176037913960400191505060405180910390fd5b6108127f751b795d24b92e3d92d1d0d8f2885f4e9c9c269da350af36ae6b49069babf4bf61080d6112f2565b611062565b610867576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603e815260200180611e7c603e913960400191505060405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff167f06d08f8705b74e3172df8733fc5da269157556d015544994a9759be4693d3ff58242604051808381526020018281526020019250505060405180910390a2609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f1983836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561095057600080fd5b505af1158015610964573d6000803e3d6000fd5b505050506040513d602081101561097a57600080fd5b8101908080519060200190929190505050505050565b600060656000838152602001908152602001600020600201549050919050565b6109d760656000848152602001908152602001600020600201546109d26112f2565b611062565b610a2c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180611e4d602f913960400191505060405180910390fd5b610a3682826112fa565b5050565b610a426112f2565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610ac5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f81526020018061205b602f913960400191505060405180910390fd5b610acf828261138e565b5050565b6000610b067f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a610b016112f2565b611062565b610b5b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260368152602001806120256036913960400191505060405180910390fd5b610b63611422565b6001905090565b610b72610f82565b15610be5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b60008211610c3e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526040815260200180611fe56040913960400191505060405180910390fd5b6000609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231856040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015610cc957600080fd5b505afa158015610cdd573d6000803e3d6000fd5b505050506040513d6020811015610cf357600080fd5b8101908080519060200190929190505050905082811015610d5f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602d815260200180611eba602d913960400191505060405180910390fd5b610d676112f2565b73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614610dea576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526031815260200180611fb46031913960400191505060405180910390fd5b609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639dc29fac85856040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015610e7d57600080fd5b505af1158015610e91573d6000803e3d6000fd5b505050506040513d6020811015610ea757600080fd5b8101908080519060200190929190505050508373ffffffffffffffffffffffffffffffffffffffff167fc2db0b30181b3532965e53ba1bdf883b207dccb658a9589865c5a9c91e28b80b8484426040518084815260200180602001838152602001828103825284818151815260200191508051906020019080838360005b83811015610f40578082015181840152602081019050610f25565b50505050905090810190601f168015610f6d5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a250505050565b6000603360009054906101000a900460ff16905090565b6000610fcc7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a610fc76112f2565b611062565b611021576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603481526020018061208a6034913960400191505060405180910390fd5b61102961150d565b6001905090565b600061105a82606560008681526020019081526020016000206000016115f990919063ffffffff16565b905092915050565b600061108c826065600086815260200190815260200160002060000161161390919063ffffffff16565b905092915050565b6000801b81565b600060019054906101000a900460ff16806110ba57506110b9611643565b5b806110d0575060008054906101000a900460ff16155b611125576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611f4e602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015611175576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b61117d611654565b611185611762565b6111996000801b6111946112f2565b611870565b6111c37f751b795d24b92e3d92d1d0d8f2885f4e9c9c269da350af36ae6b49069babf4bf84611870565b6111ed7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a83611870565b6111f68461061d565b80156112175760008060016101000a81548160ff0219169083151502179055505b50505050565b600061123d6065600084815260200190815260200160002060000161187e565b9050919050565b61126b60656000848152602001908152602001600020600201546112666112f2565b611062565b6112c0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526030815260200180611ee76030913960400191505060405180910390fd5b6112ca828261138e565b5050565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b600033905090565b611322816065600085815260200190815260200160002060000161189390919063ffffffff16565b1561138a5761132f6112f2565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b6113b681606560008581526020019081526020016000206000016118c390919063ffffffff16565b1561141e576113c36112f2565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a45b5050565b61142a610f82565b61149c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f5061757361626c653a206e6f742070617573656400000000000000000000000081525060200191505060405180910390fd5b6000603360006101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6114e06112f2565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1565b611515610f82565b15611588576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b6001603360006101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586115cc6112f2565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1565b600061160883600001836118f3565b60001c905092915050565b600061163b836000018373ffffffffffffffffffffffffffffffffffffffff1660001b611976565b905092915050565b600061164e30611999565b15905090565b600060019054906101000a900460ff16806116735750611672611643565b5b80611689575060008054906101000a900460ff16155b6116de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611f4e602e913960400191505060405180910390fd5b60008060019054906101000a900460ff16159050801561172e576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b6117366119ac565b61173e611aaa565b801561175f5760008060016101000a81548160ff0219169083151502179055505b50565b600060019054906101000a900460ff16806117815750611780611643565b5b80611797575060008054906101000a900460ff16155b6117ec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611f4e602e913960400191505060405180910390fd5b60008060019054906101000a900460ff16159050801561183c576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b6118446119ac565b61184c611ba8565b801561186d5760008060016101000a81548160ff0219169083151502179055505b50565b61187a82826112fa565b5050565b600061188c82600001611cc1565b9050919050565b60006118bb836000018373ffffffffffffffffffffffffffffffffffffffff1660001b611cd2565b905092915050565b60006118eb836000018373ffffffffffffffffffffffffffffffffffffffff1660001b611d42565b905092915050565b600081836000018054905011611954576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180611e2b6022913960400191505060405180910390fd5b82600001828154811061196357fe5b9060005260206000200154905092915050565b600080836001016000848152602001908152602001600020541415905092915050565b600080823b905060008111915050919050565b600060019054906101000a900460ff16806119cb57506119ca611643565b5b806119e1575060008054906101000a900460ff16155b611a36576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611f4e602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015611a86576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b8015611aa75760008060016101000a81548160ff0219169083151502179055505b50565b600060019054906101000a900460ff1680611ac95750611ac8611643565b5b80611adf575060008054906101000a900460ff16155b611b34576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611f4e602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015611b84576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b8015611ba55760008060016101000a81548160ff0219169083151502179055505b50565b600060019054906101000a900460ff1680611bc75750611bc6611643565b5b80611bdd575060008054906101000a900460ff16155b611c32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611f4e602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015611c82576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b6000603360006101000a81548160ff0219169083151502179055508015611cbe5760008060016101000a81548160ff0219169083151502179055505b50565b600081600001805490509050919050565b6000611cde8383611976565b611d37578260000182908060018154018082558091505060019003906000526020600020016000909190919091505582600001805490508360010160008481526020019081526020016000208190555060019050611d3c565b600090505b92915050565b60008083600101600084815260200190815260200160002054905060008114611e1e5760006001820390506000600186600001805490500390506000866000018281548110611d8d57fe5b9060005260206000200154905080876000018481548110611daa57fe5b9060005260206000200181905550600183018760010160008381526020019081526020016000208190555086600001805480611de257fe5b60019003818190600052602060002001600090559055866001016000878152602001908152602001600020600090556001945050505050611e24565b60009150505b9291505056fe456e756d657261626c655365743a20696e646578206f7574206f6620626f756e6473416363657373436f6e74726f6c3a2073656e646572206d75737420626520616e2061646d696e20746f206772616e74546f6b656e577261707065723a204f6e6c79206272696467652061646d696e2063616e206d696e74206e657720746f6b656e7320666f7220612075736572546f6b656e577261707065723a20496e737566666369656e742062616c616e636520666f72206163636f756e74416363657373436f6e74726f6c3a2073656e646572206d75737420626520616e2061646d696e20746f207265766f6b65546f6b656e577261707065723a204e756d626572206f6620746f6b656e732073686f756c642062652067726561746572207468616e2030496e697469616c697a61626c653a20636f6e747261637420697320616c726561647920696e697469616c697a6564546f6b656e577261707065723a2055736572206e6f7420617574686f726973656420746f207365742055546f6b656e20636f6e7472616374546f6b656e577261707065723a2057697468647261772063616e206f6e6c7920626520646f6e65206279205374616b6572546f6b656e577261707065723a204e756d626572206f6620756e7374616b656420746f6b656e732073686f756c642062652067726561746572207468616e2030546f6b656e577261707065723a2055736572206e6f7420617574686f726973656420746f20756e706175736520636f6e747261637473416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636520726f6c657320666f722073656c66546f6b656e577261707065723a2055736572206e6f7420617574686f726973656420746f20706175736520636f6e747261637473a2646970667358221220efd88b535d38d2629e51bb36548062831c649997a603cbd79495ac38bb92396f64736f6c63430007060033"

// DeployAbi deploys a new Ethereum contract, binding an instance of Abi to it.
func DeployAbi(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Abi, error) {
	parsed, err := abi.JSON(strings.NewReader(AbiABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AbiBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Abi{AbiCaller: AbiCaller{contract: contract}, AbiTransactor: AbiTransactor{contract: contract}, AbiFilterer: AbiFilterer{contract: contract}}, nil
}

// Abi is an auto generated Go binding around an Ethereum contract.
type Abi struct {
	AbiCaller     // Read-only binding to the contract
	AbiTransactor // Write-only binding to the contract
	AbiFilterer   // Log filterer for contract events
}

// AbiCaller is an auto generated read-only Go binding around an Ethereum contract.
type AbiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AbiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AbiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AbiSession struct {
	Contract     *Abi              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AbiCallerSession struct {
	Contract *AbiCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AbiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AbiTransactorSession struct {
	Contract     *AbiTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbiRaw is an auto generated low-level Go binding around an Ethereum contract.
type AbiRaw struct {
	Contract *Abi // Generic contract binding to access the raw methods on
}

// AbiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AbiCallerRaw struct {
	Contract *AbiCaller // Generic read-only contract binding to access the raw methods on
}

// AbiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AbiTransactorRaw struct {
	Contract *AbiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAbi creates a new instance of Abi, bound to a specific deployed contract.
func NewAbi(address common.Address, backend bind.ContractBackend) (*Abi, error) {
	contract, err := bindAbi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Abi{AbiCaller: AbiCaller{contract: contract}, AbiTransactor: AbiTransactor{contract: contract}, AbiFilterer: AbiFilterer{contract: contract}}, nil
}

// NewAbiCaller creates a new read-only instance of Abi, bound to a specific deployed contract.
func NewAbiCaller(address common.Address, caller bind.ContractCaller) (*AbiCaller, error) {
	contract, err := bindAbi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbiCaller{contract: contract}, nil
}

// NewAbiTransactor creates a new write-only instance of Abi, bound to a specific deployed contract.
func NewAbiTransactor(address common.Address, transactor bind.ContractTransactor) (*AbiTransactor, error) {
	contract, err := bindAbi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbiTransactor{contract: contract}, nil
}

// NewAbiFilterer creates a new log filterer instance of Abi, bound to a specific deployed contract.
func NewAbiFilterer(address common.Address, filterer bind.ContractFilterer) (*AbiFilterer, error) {
	contract, err := bindAbi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbiFilterer{contract: contract}, nil
}

// bindAbi binds a generic wrapper to an already deployed contract.
func bindAbi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AbiABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abi *AbiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abi.Contract.AbiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abi *AbiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.Contract.AbiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abi *AbiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abi.Contract.AbiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abi *AbiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abi.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abi *AbiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abi *AbiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abi.Contract.contract.Transact(opts, method, params...)
}

// BRIDGEADMINROLE is a free data retrieval call binding the contract method 0x118c38c7.
//
// Solidity: function BRIDGE_ADMIN_ROLE() view returns(bytes32)
func (_Abi *AbiCaller) BRIDGEADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "BRIDGE_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BRIDGEADMINROLE is a free data retrieval call binding the contract method 0x118c38c7.
//
// Solidity: function BRIDGE_ADMIN_ROLE() view returns(bytes32)
func (_Abi *AbiSession) BRIDGEADMINROLE() ([32]byte, error) {
	return _Abi.Contract.BRIDGEADMINROLE(&_Abi.CallOpts)
}

// BRIDGEADMINROLE is a free data retrieval call binding the contract method 0x118c38c7.
//
// Solidity: function BRIDGE_ADMIN_ROLE() view returns(bytes32)
func (_Abi *AbiCallerSession) BRIDGEADMINROLE() ([32]byte, error) {
	return _Abi.Contract.BRIDGEADMINROLE(&_Abi.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Abi *AbiCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Abi *AbiSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Abi.Contract.DEFAULTADMINROLE(&_Abi.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Abi *AbiCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Abi.Contract.DEFAULTADMINROLE(&_Abi.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_Abi *AbiCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_Abi *AbiSession) PAUSERROLE() ([32]byte, error) {
	return _Abi.Contract.PAUSERROLE(&_Abi.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_Abi *AbiCallerSession) PAUSERROLE() ([32]byte, error) {
	return _Abi.Contract.PAUSERROLE(&_Abi.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Abi *AbiCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Abi *AbiSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Abi.Contract.GetRoleAdmin(&_Abi.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Abi *AbiCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Abi.Contract.GetRoleAdmin(&_Abi.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Abi *AbiCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Abi *AbiSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Abi.Contract.GetRoleMember(&_Abi.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Abi *AbiCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Abi.Contract.GetRoleMember(&_Abi.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Abi *AbiCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Abi *AbiSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Abi.Contract.GetRoleMemberCount(&_Abi.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Abi *AbiCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Abi.Contract.GetRoleMemberCount(&_Abi.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Abi *AbiCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Abi *AbiSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Abi.Contract.HasRole(&_Abi.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Abi *AbiCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Abi.Contract.HasRole(&_Abi.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Abi *AbiCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Abi.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Abi *AbiSession) Paused() (bool, error) {
	return _Abi.Contract.Paused(&_Abi.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Abi *AbiCallerSession) Paused() (bool, error) {
	return _Abi.Contract.Paused(&_Abi.CallOpts)
}

// GenerateUTokens is a paid mutator transaction binding the contract method 0x22dd9bc3.
//
// Solidity: function generateUTokens(address to, uint256 amount) returns()
func (_Abi *AbiTransactor) GenerateUTokens(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "generateUTokens", to, amount)
}

// GenerateUTokens is a paid mutator transaction binding the contract method 0x22dd9bc3.
//
// Solidity: function generateUTokens(address to, uint256 amount) returns()
func (_Abi *AbiSession) GenerateUTokens(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Abi.Contract.GenerateUTokens(&_Abi.TransactOpts, to, amount)
}

// GenerateUTokens is a paid mutator transaction binding the contract method 0x22dd9bc3.
//
// Solidity: function generateUTokens(address to, uint256 amount) returns()
func (_Abi *AbiTransactorSession) GenerateUTokens(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Abi.Contract.GenerateUTokens(&_Abi.TransactOpts, to, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Abi *AbiTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Abi *AbiSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.Contract.GrantRole(&_Abi.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Abi *AbiTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.Contract.GrantRole(&_Abi.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address uAddress, address bridgeAdminAddress, address pauserAddress) returns()
func (_Abi *AbiTransactor) Initialize(opts *bind.TransactOpts, uAddress common.Address, bridgeAdminAddress common.Address, pauserAddress common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "initialize", uAddress, bridgeAdminAddress, pauserAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address uAddress, address bridgeAdminAddress, address pauserAddress) returns()
func (_Abi *AbiSession) Initialize(uAddress common.Address, bridgeAdminAddress common.Address, pauserAddress common.Address) (*types.Transaction, error) {
	return _Abi.Contract.Initialize(&_Abi.TransactOpts, uAddress, bridgeAdminAddress, pauserAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address uAddress, address bridgeAdminAddress, address pauserAddress) returns()
func (_Abi *AbiTransactorSession) Initialize(uAddress common.Address, bridgeAdminAddress common.Address, pauserAddress common.Address) (*types.Transaction, error) {
	return _Abi.Contract.Initialize(&_Abi.TransactOpts, uAddress, bridgeAdminAddress, pauserAddress)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool success)
func (_Abi *AbiTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool success)
func (_Abi *AbiSession) Pause() (*types.Transaction, error) {
	return _Abi.Contract.Pause(&_Abi.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool success)
func (_Abi *AbiTransactorSession) Pause() (*types.Transaction, error) {
	return _Abi.Contract.Pause(&_Abi.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Abi *AbiTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Abi *AbiSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.Contract.RenounceRole(&_Abi.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Abi *AbiTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.Contract.RenounceRole(&_Abi.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Abi *AbiTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Abi *AbiSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.Contract.RevokeRole(&_Abi.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Abi *AbiTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Abi.Contract.RevokeRole(&_Abi.TransactOpts, role, account)
}

// SetUTokensContract is a paid mutator transaction binding the contract method 0x21bdf9e5.
//
// Solidity: function setUTokensContract(address uAddress) returns()
func (_Abi *AbiTransactor) SetUTokensContract(opts *bind.TransactOpts, uAddress common.Address) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "setUTokensContract", uAddress)
}

// SetUTokensContract is a paid mutator transaction binding the contract method 0x21bdf9e5.
//
// Solidity: function setUTokensContract(address uAddress) returns()
func (_Abi *AbiSession) SetUTokensContract(uAddress common.Address) (*types.Transaction, error) {
	return _Abi.Contract.SetUTokensContract(&_Abi.TransactOpts, uAddress)
}

// SetUTokensContract is a paid mutator transaction binding the contract method 0x21bdf9e5.
//
// Solidity: function setUTokensContract(address uAddress) returns()
func (_Abi *AbiTransactorSession) SetUTokensContract(uAddress common.Address) (*types.Transaction, error) {
	return _Abi.Contract.SetUTokensContract(&_Abi.TransactOpts, uAddress)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool success)
func (_Abi *AbiTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool success)
func (_Abi *AbiSession) Unpause() (*types.Transaction, error) {
	return _Abi.Contract.Unpause(&_Abi.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool success)
func (_Abi *AbiTransactorSession) Unpause() (*types.Transaction, error) {
	return _Abi.Contract.Unpause(&_Abi.TransactOpts)
}

// WithdrawUTokens is a paid mutator transaction binding the contract method 0x428bee9e.
//
// Solidity: function withdrawUTokens(address from, uint256 tokens, string toAtomAddress) returns()
func (_Abi *AbiTransactor) WithdrawUTokens(opts *bind.TransactOpts, from common.Address, tokens *big.Int, toAtomAddress string) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "withdrawUTokens", from, tokens, toAtomAddress)
}

// WithdrawUTokens is a paid mutator transaction binding the contract method 0x428bee9e.
//
// Solidity: function withdrawUTokens(address from, uint256 tokens, string toAtomAddress) returns()
func (_Abi *AbiSession) WithdrawUTokens(from common.Address, tokens *big.Int, toAtomAddress string) (*types.Transaction, error) {
	return _Abi.Contract.WithdrawUTokens(&_Abi.TransactOpts, from, tokens, toAtomAddress)
}

// WithdrawUTokens is a paid mutator transaction binding the contract method 0x428bee9e.
//
// Solidity: function withdrawUTokens(address from, uint256 tokens, string toAtomAddress) returns()
func (_Abi *AbiTransactorSession) WithdrawUTokens(from common.Address, tokens *big.Int, toAtomAddress string) (*types.Transaction, error) {
	return _Abi.Contract.WithdrawUTokens(&_Abi.TransactOpts, from, tokens, toAtomAddress)
}

// AbiGenerateUTokensIterator is returned from FilterGenerateUTokens and is used to iterate over the raw logs and unpacked data for GenerateUTokens events raised by the Abi contract.
type AbiGenerateUTokensIterator struct {
	Event *AbiGenerateUTokens // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiGenerateUTokensIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiGenerateUTokens)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiGenerateUTokens)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiGenerateUTokensIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiGenerateUTokensIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiGenerateUTokens represents a GenerateUTokens event raised by the Abi contract.
type AbiGenerateUTokens struct {
	AccountAddress common.Address
	Tokens         *big.Int
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterGenerateUTokens is a free log retrieval operation binding the contract event 0x06d08f8705b74e3172df8733fc5da269157556d015544994a9759be4693d3ff5.
//
// Solidity: event GenerateUTokens(address indexed accountAddress, uint256 tokens, uint256 timestamp)
func (_Abi *AbiFilterer) FilterGenerateUTokens(opts *bind.FilterOpts, accountAddress []common.Address) (*AbiGenerateUTokensIterator, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "GenerateUTokens", accountAddressRule)
	if err != nil {
		return nil, err
	}
	return &AbiGenerateUTokensIterator{contract: _Abi.contract, event: "GenerateUTokens", logs: logs, sub: sub}, nil
}

// WatchGenerateUTokens is a free log subscription operation binding the contract event 0x06d08f8705b74e3172df8733fc5da269157556d015544994a9759be4693d3ff5.
//
// Solidity: event GenerateUTokens(address indexed accountAddress, uint256 tokens, uint256 timestamp)
func (_Abi *AbiFilterer) WatchGenerateUTokens(opts *bind.WatchOpts, sink chan<- *AbiGenerateUTokens, accountAddress []common.Address) (event.Subscription, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "GenerateUTokens", accountAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiGenerateUTokens)
				if err := _Abi.contract.UnpackLog(event, "GenerateUTokens", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGenerateUTokens is a log parse operation binding the contract event 0x06d08f8705b74e3172df8733fc5da269157556d015544994a9759be4693d3ff5.
//
// Solidity: event GenerateUTokens(address indexed accountAddress, uint256 tokens, uint256 timestamp)
func (_Abi *AbiFilterer) ParseGenerateUTokens(log types.Log) (*AbiGenerateUTokens, error) {
	event := new(AbiGenerateUTokens)
	if err := _Abi.contract.UnpackLog(event, "GenerateUTokens", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Abi contract.
type AbiPausedIterator struct {
	Event *AbiPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPaused represents a Paused event raised by the Abi contract.
type AbiPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Abi *AbiFilterer) FilterPaused(opts *bind.FilterOpts) (*AbiPausedIterator, error) {

	logs, sub, err := _Abi.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AbiPausedIterator{contract: _Abi.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Abi *AbiFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AbiPaused) (event.Subscription, error) {

	logs, sub, err := _Abi.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPaused)
				if err := _Abi.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Abi *AbiFilterer) ParsePaused(log types.Log) (*AbiPaused, error) {
	event := new(AbiPaused)
	if err := _Abi.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Abi contract.
type AbiRoleAdminChangedIterator struct {
	Event *AbiRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiRoleAdminChanged represents a RoleAdminChanged event raised by the Abi contract.
type AbiRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Abi *AbiFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AbiRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AbiRoleAdminChangedIterator{contract: _Abi.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Abi *AbiFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AbiRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiRoleAdminChanged)
				if err := _Abi.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Abi *AbiFilterer) ParseRoleAdminChanged(log types.Log) (*AbiRoleAdminChanged, error) {
	event := new(AbiRoleAdminChanged)
	if err := _Abi.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Abi contract.
type AbiRoleGrantedIterator struct {
	Event *AbiRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiRoleGranted represents a RoleGranted event raised by the Abi contract.
type AbiRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Abi *AbiFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AbiRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AbiRoleGrantedIterator{contract: _Abi.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Abi *AbiFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AbiRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiRoleGranted)
				if err := _Abi.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Abi *AbiFilterer) ParseRoleGranted(log types.Log) (*AbiRoleGranted, error) {
	event := new(AbiRoleGranted)
	if err := _Abi.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Abi contract.
type AbiRoleRevokedIterator struct {
	Event *AbiRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiRoleRevoked represents a RoleRevoked event raised by the Abi contract.
type AbiRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Abi *AbiFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AbiRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AbiRoleRevokedIterator{contract: _Abi.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Abi *AbiFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AbiRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiRoleRevoked)
				if err := _Abi.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Abi *AbiFilterer) ParseRoleRevoked(log types.Log) (*AbiRoleRevoked, error) {
	event := new(AbiRoleRevoked)
	if err := _Abi.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiSetUTokensContractIterator is returned from FilterSetUTokensContract and is used to iterate over the raw logs and unpacked data for SetUTokensContract events raised by the Abi contract.
type AbiSetUTokensContractIterator struct {
	Event *AbiSetUTokensContract // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiSetUTokensContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiSetUTokensContract)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiSetUTokensContract)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiSetUTokensContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiSetUTokensContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiSetUTokensContract represents a SetUTokensContract event raised by the Abi contract.
type AbiSetUTokensContract struct {
	Contract common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetUTokensContract is a free log retrieval operation binding the contract event 0xf745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c.
//
// Solidity: event SetUTokensContract(address indexed _contract)
func (_Abi *AbiFilterer) FilterSetUTokensContract(opts *bind.FilterOpts, _contract []common.Address) (*AbiSetUTokensContractIterator, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "SetUTokensContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return &AbiSetUTokensContractIterator{contract: _Abi.contract, event: "SetUTokensContract", logs: logs, sub: sub}, nil
}

// WatchSetUTokensContract is a free log subscription operation binding the contract event 0xf745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c.
//
// Solidity: event SetUTokensContract(address indexed _contract)
func (_Abi *AbiFilterer) WatchSetUTokensContract(opts *bind.WatchOpts, sink chan<- *AbiSetUTokensContract, _contract []common.Address) (event.Subscription, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "SetUTokensContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiSetUTokensContract)
				if err := _Abi.contract.UnpackLog(event, "SetUTokensContract", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetUTokensContract is a log parse operation binding the contract event 0xf745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c.
//
// Solidity: event SetUTokensContract(address indexed _contract)
func (_Abi *AbiFilterer) ParseSetUTokensContract(log types.Log) (*AbiSetUTokensContract, error) {
	event := new(AbiSetUTokensContract)
	if err := _Abi.contract.UnpackLog(event, "SetUTokensContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Abi contract.
type AbiUnpausedIterator struct {
	Event *AbiUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiUnpaused represents a Unpaused event raised by the Abi contract.
type AbiUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Abi *AbiFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AbiUnpausedIterator, error) {

	logs, sub, err := _Abi.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AbiUnpausedIterator{contract: _Abi.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Abi *AbiFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AbiUnpaused) (event.Subscription, error) {

	logs, sub, err := _Abi.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiUnpaused)
				if err := _Abi.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Abi *AbiFilterer) ParseUnpaused(log types.Log) (*AbiUnpaused, error) {
	event := new(AbiUnpaused)
	if err := _Abi.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiWithdrawUTokensIterator is returned from FilterWithdrawUTokens and is used to iterate over the raw logs and unpacked data for WithdrawUTokens events raised by the Abi contract.
type AbiWithdrawUTokensIterator struct {
	Event *AbiWithdrawUTokens // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbiWithdrawUTokensIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiWithdrawUTokens)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbiWithdrawUTokens)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbiWithdrawUTokensIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiWithdrawUTokensIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiWithdrawUTokens represents a WithdrawUTokens event raised by the Abi contract.
type AbiWithdrawUTokens struct {
	AccountAddress common.Address
	Tokens         *big.Int
	ToAtomAddress  string
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdrawUTokens is a free log retrieval operation binding the contract event 0xc2db0b30181b3532965e53ba1bdf883b207dccb658a9589865c5a9c91e28b80b.
//
// Solidity: event WithdrawUTokens(address indexed accountAddress, uint256 tokens, string toAtomAddress, uint256 timestamp)
func (_Abi *AbiFilterer) FilterWithdrawUTokens(opts *bind.FilterOpts, accountAddress []common.Address) (*AbiWithdrawUTokensIterator, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "WithdrawUTokens", accountAddressRule)
	if err != nil {
		return nil, err
	}
	return &AbiWithdrawUTokensIterator{contract: _Abi.contract, event: "WithdrawUTokens", logs: logs, sub: sub}, nil
}

// WatchWithdrawUTokens is a free log subscription operation binding the contract event 0xc2db0b30181b3532965e53ba1bdf883b207dccb658a9589865c5a9c91e28b80b.
//
// Solidity: event WithdrawUTokens(address indexed accountAddress, uint256 tokens, string toAtomAddress, uint256 timestamp)
func (_Abi *AbiFilterer) WatchWithdrawUTokens(opts *bind.WatchOpts, sink chan<- *AbiWithdrawUTokens, accountAddress []common.Address) (event.Subscription, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "WithdrawUTokens", accountAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiWithdrawUTokens)
				if err := _Abi.contract.UnpackLog(event, "WithdrawUTokens", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawUTokens is a log parse operation binding the contract event 0xc2db0b30181b3532965e53ba1bdf883b207dccb658a9589865c5a9c91e28b80b.
//
// Solidity: event WithdrawUTokens(address indexed accountAddress, uint256 tokens, string toAtomAddress, uint256 timestamp)
func (_Abi *AbiFilterer) ParseWithdrawUTokens(log types.Log) (*AbiWithdrawUTokens, error) {
	event := new(AbiWithdrawUTokens)
	if err := _Abi.contract.UnpackLog(event, "WithdrawUTokens", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
