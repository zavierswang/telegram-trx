package grid

import (
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/pkg/errors"
)

type tronAddress struct {
	address string
}

func (t tronAddress) String() string {
	return t.address
}

func (t *tronAddress) Set(s string) error {
	_, err := address.Base58ToAddress(s)
	if err != nil {
		return errors.Wrap(err, "not a valid address")
	}
	t.address = s
	return nil
}

func (t *tronAddress) GetAddress() address.Address {
	addr, err := address.Base58ToAddress(t.address)
	if err != nil {
		return nil
	}
	return addr
}

func (t tronAddress) Type() string {
	return "tron-address"
}
