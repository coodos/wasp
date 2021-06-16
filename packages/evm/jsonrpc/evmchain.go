package jsonrpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/contracts/native/evmchain"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/root"
)

type EVMChain struct {
	backend ChainBackend
}

func NewEVMChain(backend ChainBackend) *EVMChain {
	return &EVMChain{backend}
}

func (e *EVMChain) GasPerIota() (uint64, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetGasPerIota)
	if err != nil {
		return 0, err
	}
	gpi, _, err := codec.DecodeUint64(ret.MustGet(evmchain.FieldResult))
	return gpi, err
}

func (e *EVMChain) BlockNumber() (*big.Int, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetBlockNumber)
	if err != nil {
		return nil, err
	}

	bal := big.NewInt(0)
	bal.SetBytes(ret.MustGet(evmchain.FieldResult))
	return bal, nil
}

func (e *EVMChain) FeeColor() (ledgerstate.Color, error) {
	feeInfo, err := e.backend.CallView(root.Interface.Name, root.FuncGetFeeInfo,
		root.ParamHname, evmchain.Interface.Hname(),
	)
	if err != nil {
		return ledgerstate.Color{}, err
	}
	feeColor, _, err := codec.DecodeColor(feeInfo.MustGet(root.ParamFeeColor))
	return feeColor, err
}

func (e *EVMChain) GasLimitFee(tx *types.Transaction) (map[ledgerstate.Color]uint64, error) {
	gpi, err := e.GasPerIota()
	if err != nil {
		return nil, err
	}
	feeColor, err := e.FeeColor()
	if err != nil {
		return nil, err
	}
	return map[ledgerstate.Color]uint64{feeColor: tx.Gas() / gpi}, nil
}

func (e *EVMChain) SendTransaction(keyPair *ed25519.KeyPair, tx *types.Transaction) error {
	transfer, err := e.GasLimitFee(tx)
	if err != nil {
		return err
	}
	txdata, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	return e.backend.PostRequest(keyPair, transfer, evmchain.Interface.Name, evmchain.FuncSendTransaction, 1024,
		evmchain.FieldTransactionData, txdata,
	)
}

func paramsWithOptionalBlockNumber(blockNumber *big.Int, params ...interface{}) []interface{} {
	if blockNumber != nil {
		return append(params, evmchain.FieldBlockNumber, blockNumber.Bytes())
	}
	return params
}

func (e *EVMChain) Balance(address common.Address, blockNumber *big.Int) (*big.Int, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetBalance, paramsWithOptionalBlockNumber(blockNumber,
		evmchain.FieldAddress, address.Bytes(),
	)...)
	if err != nil {
		return nil, err
	}

	bal := big.NewInt(0)
	bal.SetBytes(ret.MustGet(evmchain.FieldResult))
	return bal, nil
}

func (e *EVMChain) Code(address common.Address, blockNumber *big.Int) ([]byte, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetCode, paramsWithOptionalBlockNumber(blockNumber,
		evmchain.FieldAddress, address.Bytes(),
	)...)
	if err != nil {
		return nil, err
	}
	return ret.MustGet(evmchain.FieldResult), nil
}

func (e *EVMChain) BlockByNumber(blockNumber *big.Int) (*types.Block, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetBlockByNumber, paramsWithOptionalBlockNumber(blockNumber)...)
	if err != nil {
		return nil, err
	}

	if !ret.MustHas(evmchain.FieldResult) {
		return nil, nil
	}

	block, err := evmchain.DecodeBlock(ret.MustGet(evmchain.FieldResult))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (e *EVMChain) getTransaction(funcName string, args ...interface{}) (tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, err error) {
	var ret dict.Dict
	ret, err = e.backend.CallView(evmchain.Interface.Name, funcName, args...)
	if err != nil {
		return
	}

	if !ret.MustHas(evmchain.FieldTransaction) {
		return
	}

	tx, err = evmchain.DecodeTransaction(ret.MustGet(evmchain.FieldTransaction))
	if err != nil {
		return
	}
	blockHash = common.BytesToHash(ret.MustGet(evmchain.FieldBlockHash))
	blockNumber, _, err = codec.DecodeUint64(ret.MustGet(evmchain.FieldBlockNumber))
	if err != nil {
		return
	}
	index, _, err = codec.DecodeUint64(ret.MustGet(evmchain.FieldTransactionIndex))
	if err != nil {
		return
	}
	return
}

func (e *EVMChain) TransactionByHash(hash common.Hash) (tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, err error) {
	return e.getTransaction(evmchain.FuncGetTransactionByHash,
		evmchain.FieldTransactionHash, hash.Bytes(),
	)
}

func (e *EVMChain) TransactionByBlockHashAndIndex(hash common.Hash, index uint64) (tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index_ uint64, err error) {
	return e.getTransaction(evmchain.FuncGetTransactionByBlockHashAndIndex,
		evmchain.FieldBlockHash, hash.Bytes(),
		evmchain.FieldTransactionIndex, codec.EncodeUint64(index),
	)
}

func (e *EVMChain) TransactionByBlockNumberAndIndex(blockNumber *big.Int, index uint64) (tx *types.Transaction, blockHash common.Hash, blockNumber_ uint64, index_ uint64, err error) {
	return e.getTransaction(evmchain.FuncGetTransactionByBlockNumberAndIndex, paramsWithOptionalBlockNumber(blockNumber,
		evmchain.FieldTransactionIndex, codec.EncodeUint64(index),
	)...)
}

func (e *EVMChain) BlockByHash(hash common.Hash) (*types.Block, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetBlockByHash,
		evmchain.FieldBlockHash, hash.Bytes(),
	)
	if err != nil {
		return nil, err
	}

	if !ret.MustHas(evmchain.FieldResult) {
		return nil, nil
	}

	block, err := evmchain.DecodeBlock(ret.MustGet(evmchain.FieldResult))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (e *EVMChain) TransactionReceipt(txHash common.Hash) (*evmchain.Receipt, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetReceipt,
		evmchain.FieldTransactionHash, txHash.Bytes(),
	)
	if err != nil {
		return nil, err
	}

	if !ret.MustHas(evmchain.FieldResult) {
		return nil, nil
	}

	receipt, err := evmchain.DecodeReceipt(ret.MustGet(evmchain.FieldResult))
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (e *EVMChain) TransactionCount(address common.Address, blockNumber *big.Int) (uint64, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetNonce, paramsWithOptionalBlockNumber(blockNumber,
		evmchain.FieldAddress, address.Bytes(),
	)...)
	if err != nil {
		return 0, err
	}

	n, _, err := codec.DecodeUint64(ret.MustGet(evmchain.FieldResult))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (e *EVMChain) CallContract(args ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncCallContract, paramsWithOptionalBlockNumber(blockNumber,
		evmchain.FieldCallMsg, evmchain.EncodeCallMsg(args),
	)...)
	if err != nil {
		return nil, err
	}
	return ret.MustGet(evmchain.FieldResult), nil
}

func (e *EVMChain) EstimateGas(args ethereum.CallMsg) (uint64, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncEstimateGas,
		evmchain.FieldCallMsg, evmchain.EncodeCallMsg(args),
	)
	if err != nil {
		return 0, err
	}
	gas, _, err := codec.DecodeUint64(ret.MustGet(evmchain.FieldResult))
	return gas, err
}

func (e *EVMChain) StorageAt(address common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetStorage, paramsWithOptionalBlockNumber(blockNumber,
		evmchain.FieldAddress, address.Bytes(),
		evmchain.FieldKey, key.Bytes(),
	)...)
	if err != nil {
		return nil, err
	}
	return ret.MustGet(evmchain.FieldResult), nil
}

func (e *EVMChain) BlockTransactionCountByHash(blockHash common.Hash) (uint64, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetBlockTransactionCountByHash,
		evmchain.FieldBlockHash, blockHash.Bytes(),
	)
	if err != nil {
		return 0, err
	}
	n, _, err := codec.DecodeUint64(ret.MustGet(evmchain.FieldResult))
	return n, err
}

func (e *EVMChain) BlockTransactionCountByNumber(blockNumber *big.Int) (uint64, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetBlockTransactionCountByNumber, paramsWithOptionalBlockNumber(blockNumber)...)
	if err != nil {
		return 0, err
	}
	n, _, err := codec.DecodeUint64(ret.MustGet(evmchain.FieldResult))
	return n, err
}

func (e *EVMChain) Logs(q *ethereum.FilterQuery) ([]*types.Log, error) {
	ret, err := e.backend.CallView(evmchain.Interface.Name, evmchain.FuncGetLogs,
		evmchain.FieldFilterQuery, evmchain.EncodeFilterQuery(q),
	)
	if err != nil {
		return nil, err
	}
	return evmchain.DecodeLogs(ret.MustGet(evmchain.FieldResult))
}
