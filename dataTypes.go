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

func (dataType DataType) SizeOf() int {
	switch dataType {
	case DataTypeVoid:
		return 0
	case DataTypeI8:
		return 1
	case DataTypeI16:
		return 2
	case DataTypeI32:
		return 4
	case DataTypeI64:
		return 8
	case DataTypeU8:
		return 1
	case DataTypeU16:
		return 2
	case DataTypeU32:
		return 4
	case DataTypeU64:
		return 8
	case DataTypeSingleFloat:
		return 4
	case DataTypeDoubleFloat:
		return 8
	case DataTypeBoolean:
		return 1
	case DataTypeTimestamp:
		return 16
	case DataTypeComplexSingleFloat:
		return 8
	case DataTypeComplexDoubleFloat:
		return 16
	default:
		return -1
	}
}
