package model

type AddressSpace [65536]Data8

func NewAddressSpace() *AddressSpace {
	return &AddressSpace{}
}
