package accountsc

import (
	"fmt"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/vm/contract"
	"github.com/iotaledger/wasp/packages/vm/vmtypes"
)

const (
	Name        = "accounts"
	Version     = "0.1"
	description = "Chain account ledger contract"
	fullName    = Name + "-" + Version
)

var (
	Interface = &contract.ContractInterface{
		Name:        fullName,
		Description: description,
		ProgramHash: *hashing.HashStrings(fullName),
	}
	TotalAssetsAccountID = coretypes.NewAgentIDFromContractID(coretypes.NewContractID(coretypes.ChainID{}, Interface.Hname()))
)

func init() {
	Interface.WithFunctions(initialize, []contract.ContractFunctionInterface{
		contract.ViewFunc(FuncBalance, getBalance),
		contract.ViewFunc(FuncAccounts, getAccounts),
		contract.Func(FuncDeposit, deposit),
		contract.Func(FuncMove, move),
		contract.Func(FuncAllow, allow),
		contract.Func(FuncWithdraw, withdraw),
	})
}

const (
	FuncBalance  = "balance"
	FuncDeposit  = "deposit"
	FuncMove     = "move"
	FuncAllow    = "allow"
	FuncWithdraw = "withdraw"
	FuncAccounts = "accounts"

	VarStateInitialized = "i"
	VarStateAllAccounts = "a"
	VarStateAllowances  = "l"

	ParamAgentID = "a"
	ParamColor   = "c"
	ParamAmount  = "t"
	ParamChainID = "i"
)

var (
	ErrParamWrongOrNotFound = fmt.Errorf("wrong parameters: agent ID is wrong or not found")
)

func GetProcessor() vmtypes.Processor {
	return Interface
}

func ChainOwnerAgentID(chainID coretypes.ChainID) coretypes.AgentID {
	return coretypes.NewAgentIDFromContractID(coretypes.NewContractID(chainID, Interface.Hname()))
}