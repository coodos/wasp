package vmcontext

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	valuetransaction "github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/sctransaction"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/builtinvm/root"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"github.com/iotaledger/wasp/packages/vm/statetxbuilder"
)

// context for one request
type VMContext struct {
	// same for the block
	chainID       coretypes.ChainID
	processors    *processors.ProcessorCache
	accrueFeesTo  coretypes.AgentID
	balances      map[valuetransaction.ID][]*balance.Balance
	txBuilder     *statetxbuilder.Builder // mutated
	saveTxBuilder *statetxbuilder.Builder // for rollback
	virtualState  state.VirtualState      // mutated
	log           *logger.Logger
	// request context
	entropy        hashing.HashValue // mutates with each request
	reqRef         sctransaction.RequestRef
	reqHname       coretypes.Hname
	contractRecord *root.ContractRecord
	timestamp      int64
	stateUpdate    state.StateUpdate // mutated
	callStack      []*callContext
}

type callContext struct {
	isRequestContext bool                      // is called from the request (true) or from another SC (false)
	caller           coretypes.AgentID         // calling agent
	contract         coretypes.Hname           // called contract
	params           codec.ImmutableCodec      // params passed
	transfer         coretypes.ColoredBalances // transfer passed
}

// NewVMContext:
// - creates state block in the tx builder, including moving the SC token
// - handles request tokens by moving them either to the
// reward address or sending it back to the requester
// All request tokens are handled for the whole block
func NewVMContext(task *vm.VMTask, txb *statetxbuilder.Builder) (*VMContext, error) {
	ret := &VMContext{
		processors:   task.Processors,
		chainID:      task.ChainID,
		balances:     task.Balances,
		txBuilder:    txb,
		virtualState: task.VirtualState.Clone(),
		log:          task.Log,
		entropy:      task.Entropy,
		callStack:    make([]*callContext, 0),
	}
	return ret, nil
}