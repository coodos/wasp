package codec

import (
	"encoding/hex"
	"fmt"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"sort"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/datatypes"
	"github.com/iotaledger/wasp/packages/util"
)

// MutableCodec is an interface that offers easy conversions between []byte
// and other types (including collections) when manipulating a KVStore
type MutableCodec interface {
	ImmutableCodec
	wCodec
	GetArray(kv.Key) (*datatypes.Array, error)
	GetMap(kv.Key) (*datatypes.Map, error)
	GetTimestampedLog(kv.Key) (*datatypes.TimestampedLog, error)
}

// MutableMustCodec is a MutableCodec that automatically panics on error
type MutableMustCodec interface {
	ImmutableMustCodec
	wCodec
}

// ImmutableCodec is an interface that offers easy conversions between []byte and other types when
// manipulating a read-only KVStore
// TODO replace returned pointers to values
type ImmutableCodec interface {
	KVStore() kv.KVStore

	Has(key kv.Key) (bool, error)
	Get(key kv.Key) ([]byte, error)
	GetString(key kv.Key) (string, bool, error)
	GetInt64(key kv.Key) (int64, bool, error)
	GetHname(key kv.Key) (coretypes.Hname, bool, error)
	GetAddress(key kv.Key) (*address.Address, bool, error)
	GetHashValue(key kv.Key) (*hashing.HashValue, bool, error)
	GetChainID(key kv.Key) (*coretypes.ChainID, bool, error)
	GetAgentID(key kv.Key) (*coretypes.AgentID, bool, error)
	GetContractID(key kv.Key) (*coretypes.ContractID, bool, error)
	GetColor(key kv.Key) (*balance.Color, bool, error)
	Iterate(prefix kv.Key, f func(key kv.Key, value []byte) bool) error
	IterateKeys(prefix kv.Key, f func(key kv.Key) bool) error
	Keys() ([]kv.Key, error)
	KeysSorted() ([]kv.Key, error)
}

// ImmutableMustCodec is an ImmutableCodec that automatically panics on error
type ImmutableMustCodec interface {
	KVStore() kv.KVStore

	Has(key kv.Key) bool
	Get(key kv.Key) []byte
	GetString(key kv.Key) (string, bool)
	GetInt64(key kv.Key) (int64, bool)
	GetHname(key kv.Key) (coretypes.Hname, bool)
	GetAddress(key kv.Key) (*address.Address, bool)
	GetHashValue(key kv.Key) (*hashing.HashValue, bool)
	GetChainID(key kv.Key) (*coretypes.ChainID, bool)
	GetAgentID(key kv.Key) (*coretypes.AgentID, bool)
	GetContractID(key kv.Key) (*coretypes.ContractID, bool)
	GetColor(key kv.Key) (*balance.Color, bool)
	Iterate(prefix kv.Key, f func(key kv.Key, value []byte) bool)
	IterateKeys(prefix kv.Key, f func(key kv.Key) bool)
	Keys() []kv.Key
	KeysSorted() []kv.Key
	GetArray(kv.Key) *datatypes.MustArray
	GetMap(kv.Key) *datatypes.MustMap
	GetTimestampedLog(kv.Key) *datatypes.MustTimestampedLog
}

// wCodec is an interface that offers easy conversions between []byte and other types when
// manipulating a writable KVStore
type wCodec interface {
	Del(key kv.Key)
	Set(key kv.Key, value []byte)
	SetString(key kv.Key, value string)
	SetInt64(key kv.Key, value int64)
	SetHname(key kv.Key, value coretypes.Hname)
	SetAddress(key kv.Key, value *address.Address)
	SetHashValue(key kv.Key, value *hashing.HashValue)
	SetChainID(key kv.Key, value *coretypes.ChainID)
	SetAgentID(key kv.Key, value *coretypes.AgentID)
	SetContractID(key kv.Key, value *coretypes.ContractID)
	SetColor(key kv.Key, value *balance.Color)
	Append(from ImmutableCodec) error
}

type codec struct {
	kv kv.KVStore
}

type mustcodec struct {
	codec
}

func NewCodec(kv kv.KVStore) MutableCodec {
	return codec{kv: kv}
}

func NewMustCodec(kv kv.KVStore) MutableMustCodec {
	return mustcodec{codec{kv: kv}}
}

func (c codec) KVStore() kv.KVStore {
	return c.kv
}

func (c codec) GetArray(key kv.Key) (*datatypes.Array, error) {
	return datatypes.NewArray(c, string(key))
}

func (c mustcodec) GetArray(key kv.Key) *datatypes.MustArray {
	array, err := c.codec.GetArray(key)
	if err != nil {
		panic(err)
	}
	return datatypes.NewMustArray(array)
}

func (c codec) GetMap(key kv.Key) (*datatypes.Map, error) {
	return datatypes.NewMap(c, string(key))
}

func (c mustcodec) GetMap(key kv.Key) *datatypes.MustMap {
	d, err := c.codec.GetMap(key)
	if err != nil {
		panic(err)
	}
	return datatypes.NewMustMap(d)
}

func (c codec) GetTimestampedLog(key kv.Key) (*datatypes.TimestampedLog, error) {
	return datatypes.NewTimestampedLog(c, key)
}

func (c mustcodec) GetTimestampedLog(key kv.Key) *datatypes.MustTimestampedLog {
	tlog, err := c.codec.GetTimestampedLog(key)
	if err != nil {
		panic(err)
	}
	return datatypes.NewMustTimestampedLog(tlog)
}

func (c codec) Has(key kv.Key) (bool, error) {
	return c.kv.Has(key)
}

func (c codec) Iterate(prefix kv.Key, f func(key kv.Key, value []byte) bool) error {
	return c.kv.Iterate(prefix, f)
}

func (c mustcodec) Iterate(prefix kv.Key, f func(key kv.Key, value []byte) bool) {
	err := c.kv.Iterate(prefix, f)
	if err != nil {
		panic(err)
	}
}

func (c codec) IterateKeys(prefix kv.Key, f func(key kv.Key) bool) error {
	return c.kv.IterateKeys(prefix, f)
}

func (c mustcodec) IterateKeys(prefix kv.Key, f func(key kv.Key) bool) {
	err := c.kv.IterateKeys(prefix, f)
	if err != nil {
		panic(err)
	}
}

func (c codec) Keys() ([]kv.Key, error) {
	ret := make([]kv.Key, 0)
	err := c.IterateKeys("", func(key kv.Key) bool {
		ret = append(ret, key)
		return true
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c mustcodec) Keys() []kv.Key {
	ret, err := c.codec.Keys()
	if err != nil {
		panic(err)
	}
	return ret
}

func (c codec) KeysSorted() ([]kv.Key, error) {
	k, err := c.Keys()
	if err != nil {
		return nil, err
	}
	sort.Slice(k, func(i, j int) bool {
		return k[i] < k[j]
	})
	return k, nil
}

func (c mustcodec) KeysSorted() []kv.Key {
	ret, err := c.codec.KeysSorted()
	if err != nil {
		panic(err)
	}
	return ret
}

func (c mustcodec) Has(key kv.Key) bool {
	ret, err := c.codec.Has(key)
	if err != nil {
		panic(err)
	}
	return ret
}

func (c codec) Get(key kv.Key) ([]byte, error) {
	return c.kv.Get(key)
}

func (c mustcodec) Get(key kv.Key) []byte {
	ret, err := c.codec.Get(key)
	if err != nil {
		panic(err)
	}
	return ret
}

func (c codec) Del(key kv.Key) {
	c.kv.Del(key)
}

func (c codec) Set(key kv.Key, value []byte) {
	c.kv.Set(key, value)
}

func DecodeString(b []byte) string {
	return string(b)
}

func EncodeString(value string) []byte {
	return []byte(value)
}

func (c codec) GetString(key kv.Key) (string, bool, error) {
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return "", false, err
	}
	return DecodeString(b), true, nil
}

func (c mustcodec) GetString(key kv.Key) (string, bool) {
	ret, ok, err := c.codec.GetString(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetString(key kv.Key, value string) {
	c.kv.Set(key, EncodeString(value))
}

func DecodeInt64(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, fmt.Errorf("value %s is not an int64", hex.EncodeToString(b))
	}
	return int64(util.MustUint64From8Bytes(b)), nil
}

func EncodeInt64(value int64) []byte {
	return util.Uint64To8Bytes(uint64(value))
}

func (c codec) GetInt64(key kv.Key) (int64, bool, error) {
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return 0, false, err
	}
	n, err := DecodeInt64(b)
	return n, err == nil, err
}

func (c mustcodec) GetInt64(key kv.Key) (int64, bool) {
	ret, ok, err := c.codec.GetInt64(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) GetHname(key kv.Key) (coretypes.Hname, bool, error) {
	b, err := c.Get(key)
	if err != nil || b == nil {
		return 0, false, err
	}
	ret, err := coretypes.NewHnameFromBytes(b)
	if err != nil {
		return 0, false, err
	}
	return ret, true, nil
}

func (c mustcodec) GetHname(key kv.Key) (coretypes.Hname, bool) {
	ret, ok, err := c.codec.GetHname(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetInt64(key kv.Key, value int64) {
	c.kv.Set(key, util.Uint64To8Bytes(uint64(value)))
}

func (c codec) SetHname(key kv.Key, value coretypes.Hname) {
	c.kv.Set(key, value.Bytes())
}

func (c codec) GetAddress(key kv.Key) (*address.Address, bool, error) {
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return nil, false, err
	}
	ret, _, err := address.FromBytes(b)
	if err != nil {
		return nil, false, err
	}
	return &ret, true, nil
}

func (c mustcodec) GetAddress(key kv.Key) (*address.Address, bool) {
	ret, ok, err := c.codec.GetAddress(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetAddress(key kv.Key, addr *address.Address) {
	c.kv.Set(key, addr.Bytes())
}

func (c codec) GetHashValue(key kv.Key) (*hashing.HashValue, bool, error) {
	var b []byte
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return nil, false, err
	}
	ret, err := hashing.HashValueFromBytes(b)
	return &ret, err == nil, err
}

func (c mustcodec) GetHashValue(key kv.Key) (*hashing.HashValue, bool) {
	ret, ok, err := c.codec.GetHashValue(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetHashValue(key kv.Key, h *hashing.HashValue) {
	c.kv.Set(key, h[:])
}

func (c codec) GetChainID(key kv.Key) (*coretypes.ChainID, bool, error) {
	var b []byte
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return nil, false, err
	}
	ret, err := coretypes.NewChainIDFromBytes(b)
	return &ret, err == nil, err
}

func (c mustcodec) GetChainID(key kv.Key) (*coretypes.ChainID, bool) {
	ret, ok, err := c.codec.GetChainID(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetChainID(key kv.Key, chid *coretypes.ChainID) {
	c.kv.Set(key, chid[:])
}

func (c codec) GetAgentID(key kv.Key) (*coretypes.AgentID, bool, error) {
	var b []byte
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return nil, false, err
	}
	ret, err := coretypes.NewAgentIDFromBytes(b)
	return &ret, err == nil, err
}

func (c mustcodec) GetAgentID(key kv.Key) (*coretypes.AgentID, bool) {
	ret, ok, err := c.codec.GetAgentID(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetAgentID(key kv.Key, aid *coretypes.AgentID) {
	c.kv.Set(key, aid[:])
}

func (c codec) GetContractID(key kv.Key) (*coretypes.ContractID, bool, error) {
	var b []byte
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return nil, false, err
	}
	ret, err := coretypes.NewContractIDFromBytes(b)
	return &ret, err == nil, err
}

func (c mustcodec) GetContractID(key kv.Key) (*coretypes.ContractID, bool) {
	ret, ok, err := c.codec.GetContractID(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetContractID(key kv.Key, cid *coretypes.ContractID) {
	c.kv.Set(key, cid[:])
}

func (c codec) GetColor(key kv.Key) (*balance.Color, bool, error) {
	var b []byte
	b, err := c.kv.Get(key)
	if err != nil || b == nil {
		return nil, false, err
	}
	ret, _, err := balance.ColorFromBytes(b)
	return &ret, err == nil, err
}

func (c mustcodec) GetColor(key kv.Key) (*balance.Color, bool) {
	ret, ok, err := c.codec.GetColor(key)
	if err != nil {
		panic(err)
	}
	return ret, ok
}

func (c codec) SetColor(key kv.Key, col *balance.Color) {
	c.kv.Set(key, col[:])
}

func (c codec) Append(from ImmutableCodec) error {
	return from.Iterate("", func(key kv.Key, value []byte) bool {
		c.Set(key, value)
		return true
	})
}

func EncodeDictFromMap(vars map[string]interface{}) dict.Dict {
	ret := dict.New()
	c := NewCodec(ret)
	for k, v := range vars {
		key := kv.Key(k)
		switch vt := v.(type) {
		case int:
			c.SetInt64(key, int64(vt))
		case byte:
			c.SetInt64(key, int64(vt))
		case int16:
			c.SetInt64(key, int64(vt))
		case int32:
			c.SetInt64(key, int64(vt))
		case int64:
			c.SetInt64(key, vt)
		case uint16:
			c.SetInt64(key, int64(vt))
		case uint32:
			c.SetInt64(key, int64(vt))
		case uint64:
			c.SetInt64(key, int64(vt))
		case string:
			c.SetString(key, vt)
		case []byte:
			c.Set(key, vt)
		case *hashing.HashValue:
			c.SetHashValue(key, vt)
		case hashing.HashValue:
			c.SetHashValue(key, &vt)
		case *address.Address:
			c.Set(key, vt.Bytes())
		case *balance.Color:
			c.SetColor(key, vt)
		case *coretypes.ChainID:
			c.SetChainID(key, vt)
		case *coretypes.ContractID:
			c.SetContractID(key, vt)
		case *coretypes.AgentID:
			c.SetAgentID(key, vt)
		default:
			return nil
		}
	}
	return ret
}