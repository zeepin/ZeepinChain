package common

import (
	"bytes"
	"testing"

	"github.com/zeepin/ZeepinChain/common"
	cstates "github.com/zeepin/ZeepinChain/smartcontract/states"
)

func TestSerialize(t *testing.T) {
	contract := &cstates.Contract{}
	contract.Address = common.ADDRESS_EMPTY
	contract.Method = "1"
	contract.Version = byte('1')
	contract.Args = []byte("012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")
	bf := bytes.NewBuffer(nil)
	contract.Serialize(bf)
	t.Errorf("len: %d\n", len(bf.Bytes()))
	t.Errorf("%+v", bf.Bytes())
}
