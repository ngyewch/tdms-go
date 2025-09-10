package tdms

type DataType uint32

const (
	DataTypeVoid DataType = iota
	DataTypeI8
	DataTypeI16
	DataTypeI32
	DataTypeI64
	DataTypeU8
	DataTypeU16
	DataTypeU32
	DataTypeU64
	DataTypeSingleFloat
	DataTypeDoubleFloat
	DataTypeExtendedFloat
	DataTypeSingleFloatWithUnit   = 0x19
	DataTypeDoubleFloatWithUnit   = 0x1A
	DataTypeExtendedFloatWithUnit = 0x1C
	DataTypeString                = 0x20
	DataTypeBoolean               = 0x21
	DataTypeTimestamp             = 0x44
	DataTypeFixedPoint            = 0x4f
	DataTypeComplexSingleFloat    = 0x08000c
	DataTypeComplexDoubleFloat    = 0x10000d
	DataTypeDAQmxRawData          = 0xffffffff
)
