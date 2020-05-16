package consensus

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	valuetransaction "github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/iotaledger/wasp/packages/committee"
	"github.com/iotaledger/wasp/packages/hashing"
	"time"
)

type runCalculationsParams struct {
	req             *request
	ts              time.Time
	balances        map[valuetransaction.ID][]*balance.Balance
	rewardAddress   *address.Address
	leaderPeerIndex uint16
}

// runs the VM for the request and posts result to committee's queue
func (op *operator) processRequest(par runCalculationsParams) {
	result := op.processor.Run(&runtimeContext{
		address:         op.committee.Address(),
		color:           op.committee.Color(),
		stateTx:         op.stateTx,
		balances:        par.balances,
		rewardAddress:   par.rewardAddress,
		leaderPeerIndex: par.leaderPeerIndex,
		reqMsg:          []*committee.RequestMsg{par.req.reqMsg},
		timestamp:       par.ts,
		variableState:   op.variableState,
		log:             par.req.log,
	})

	op.committee.ReceiveMessage(result)
}

func (op *operator) sendResultToTheLeader(result *committee.VMOutput) {
	reqId := result.Inputs.RequestMsg()[0].Id()
	ctx := result.Inputs.(*runtimeContext)
	log.Debugw("sendResultToTheLeader",
		"req", reqId.Short(),
		"ts", result.Inputs.Timestamp(),
		"leader", ctx.leaderPeerIndex,
	)

	sigShare, err := op.dkshare.SignShare(result.ResultTransaction.EssenceBytes())
	if err != nil {
		log.Errorf("error while signing transaction %v", err)
		return
	}
	msgData := hashing.MustBytes(&committee.SignedHashMsg{
		PeerMsgHeader: committee.PeerMsgHeader{
			StateIndex: op.stateIndex(),
		},
		RequestId:     reqId,
		OrigTimestamp: ctx.Timestamp(),
		EssenceHash:   hashing.HashData(result.ResultTransaction.EssenceBytes()),
		SigShare:      sigShare,
	})

	if err := op.committee.SendMsg(ctx.leaderPeerIndex, committee.MsgSignedHash, msgData); err != nil {
		log.Error(err)
	}
}

func (op *operator) saveOwnResult(result *committee.VMOutput) {
	reqId := result.Inputs.RequestMsg()[0].Id()
	log.Debugw("saveOwnResult",
		"req", reqId.Short(),
		"ts", result.Inputs.Timestamp(),
	)

	sigShare, err := op.dkshare.SignShare(result.ResultTransaction.EssenceBytes())
	if err != nil {
		log.Errorf("error while signing transaction %v", err)
		return
	}
	op.leaderStatus.resultTx = result.ResultTransaction
	op.leaderStatus.signedResults[op.committee.OwnPeerIndex()] = signedResult{
		essenceHash: hashing.HashData(result.ResultTransaction.EssenceBytes()),
		sigShare:    sigShare,
	}
}

func (op *operator) aggregateSigShares(sigShares [][]byte) error {
	resTx := op.leaderStatus.resultTx

	finalSignature, err := op.dkshare.RecoverFullSignature(sigShares, resTx.EssenceBytes())
	if err != nil {
		return err
	}
	finalSignature = finalSignature
	// if err := resTx.PutSignature(finalSignature); err != nil{
	// 		return fmt.Errorf("something wrong while aggregating final signature: %v", err)
	// }
	return nil
}