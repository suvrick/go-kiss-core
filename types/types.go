package types

type B byte
type I int32
type L uint64
type S string

type Value interface {
	B | I | L | S
}
