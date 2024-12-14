package config

const (
	// MinUnsignedInteger minimum unsigned integer
	MinUnsignedInteger uint = 0
	// MaxUnsignedInteger maximum unsigned integer
	MaxUnsignedInteger = ^MinUnsignedInteger
	// MaxInteger maximum integer
	MaxInteger = int(MaxUnsignedInteger >> 1)
	// MinInteger minimum integer
	// MinInteger = ^MaxInteger
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
)
