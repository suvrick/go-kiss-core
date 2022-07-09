package until

type NumberFloat interface {
	float32 | float64
}

type NumberPositive interface {
	int8 | int16 | int | int32 | int64
}

type NumberNegative interface {
	uint8 | uint16 | uint | uint32 | uint64
}

type Number interface {
	NumberNegative | NumberPositive | NumberFloat
}
