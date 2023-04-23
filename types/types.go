package types

type B uint8
type I uint32
type L uint64
type S string
type A []byte

type Value interface {
	B | I | L | S
}
