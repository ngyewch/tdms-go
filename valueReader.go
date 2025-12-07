package tdms

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

var (
	LittleEndianValueReader = NewValueReader(binary.LittleEndian)
	BigEndianValueReader    = NewValueReader(binary.BigEndian)
)

type ValueReader struct {
	byteOrder binary.ByteOrder
}

func NewValueReader(byteOrder binary.ByteOrder) *ValueReader {
	return &ValueReader{
		byteOrder: byteOrder,
	}
}

type VoidType struct{}

func (vr *ValueReader) ReadVoid(r io.Reader) (VoidType, error) {
	return VoidType{}, nil
}

func (vr *ValueReader) ReadI8(r io.Reader) (int8, error) {
	var v int8
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadI16(r io.Reader) (int16, error) {
	var v int16
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadI32(r io.Reader) (int32, error) {
	var v int32
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadI64(r io.Reader) (int64, error) {
	var v int64
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadU8(r io.Reader) (uint8, error) {
	var v uint8
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadU16(r io.Reader) (uint16, error) {
	var v uint16
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadU32(r io.Reader) (uint32, error) {
	var v uint32
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadU64(r io.Reader) (uint64, error) {
	var v uint64
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadSingleFloat(r io.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadDoubleFloat(r io.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadString(r io.Reader) (string, error) {
	stringLengthInBytes, err := vr.ReadU32(r)
	if err != nil {
		return "", err
	}
	stringBytes := make([]byte, stringLengthInBytes)
	_, err = io.ReadFull(r, stringBytes)
	if err != nil {
		return "", err
	}
	return string(stringBytes), nil
}

func (vr *ValueReader) ReadBoolean(r io.Reader) (bool, error) {
	var v bool
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return false, err
	}
	return v, nil
}

var (
	tdmsEpoch = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
)

const (
	fractionsPerNanosecond = (10 ^ -9) / (2 ^ -64)
)

func (vr *ValueReader) ReadTimestamp(r io.Reader) (time.Time, error) {
	var seconds int64
	var fractionalSeconds uint64
	var err error

	switch vr.byteOrder {
	case binary.LittleEndian:
		fractionalSeconds, err = vr.ReadU64(r)
		if err != nil {
			return time.Time{}, err
		}
		seconds, err = vr.ReadI64(r)
		if err != nil {
			return time.Time{}, err
		}
	case binary.BigEndian:
		seconds, err = vr.ReadI64(r)
		if err != nil {
			return time.Time{}, err
		}
		fractionalSeconds, err = vr.ReadU64(r)
		if err != nil {
			return time.Time{}, err
		}
	default:
		return time.Time{}, fmt.Errorf("unknown byte order")
	}
	nanoSeconds := float64(fractionalSeconds) * fractionsPerNanosecond
	return time.Unix(seconds+tdmsEpoch.Unix()-unixEpoch.Unix(), int64(math.Round(nanoSeconds))), nil
}

func (vr *ValueReader) ReadComplexSingleFloat(r io.Reader) (complex64, error) {
	var v complex64
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadComplexDoubleFloat(r io.Reader) (complex128, error) {
	var v complex128
	err := binary.Read(r, vr.byteOrder, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (vr *ValueReader) ReadDataType(r io.Reader) (DataType, error) {
	v, err := vr.ReadU32(r)
	if err != nil {
		return 0, err
	}
	return DataType(v), nil
}

// ReadDAQmxDataType reads a DAQmx data type.
// See https://github.com/adamreeve/npTDMS/blob/207463192d9fe59af43a41a3b9dc93d522384941/nptdms/daqmx.py#L321-L334
func (vr *ValueReader) ReadDAQmxDataType(r io.Reader) (DataType, error) {
	v, err := vr.ReadU32(r)
	if err != nil {
		return 0, err
	}
	switch v {
	case 0:
		return DataTypeU8, nil
	case 1:
		return DataTypeI8, nil
	case 2:
		return DataTypeU16, nil
	case 3:
		return DataTypeI16, nil
	case 4:
		return DataTypeU32, nil
	case 5:
		return DataTypeI32, nil
	case 6:
		return DataTypeU64, nil
	case 7:
		return DataTypeI64, nil
	case 8:
		return DataTypeSingleFloat, nil
	case 9:
		return DataTypeDoubleFloat, nil
	case 0xffffffff:
		return DataTypeTimestamp, nil
	default:
		return DataTypeVoid, fmt.Errorf("unknown DAQmx data type")
	}
}

func (vr *ValueReader) ReadValue(r io.Reader) (any, error) {
	dataType, err := vr.ReadU32(r)
	if err != nil {
		return nil, err
	}
	switch DataType(dataType) {
	case DataTypeVoid:
		return vr.ReadVoid(r)
	case DataTypeI8:
		return vr.ReadI8(r)
	case DataTypeI16:
		return vr.ReadI16(r)
	case DataTypeI32:
		return vr.ReadI32(r)
	case DataTypeI64:
		return vr.ReadI64(r)
	case DataTypeU8:
		return vr.ReadU8(r)
	case DataTypeU16:
		return vr.ReadU16(r)
	case DataTypeU32:
		return vr.ReadU32(r)
	case DataTypeU64:
		return vr.ReadU64(r)
	case DataTypeSingleFloat:
		return vr.ReadSingleFloat(r)
	case DataTypeDoubleFloat:
		return vr.ReadDoubleFloat(r)
	//case DataTypeExtendedFloat:
	//case DataTypeSingleFloatWithUnit:
	//case DataTypeDoubleFloatWithUnit:
	//case DataTypeExtendedFloatWithUnit:
	case DataTypeString:
		return vr.ReadString(r)
	case DataTypeBoolean:
		return vr.ReadBoolean(r)
	case DataTypeTimestamp:
		return vr.ReadTimestamp(r)
	//case DataTypeFixedPoint:
	case DataTypeComplexSingleFloat:
		return vr.ReadComplexSingleFloat(r)
	case DataTypeComplexDoubleFloat:
		return vr.ReadComplexDoubleFloat(r)
	//case DataTypeDAQmxRawData:
	default:
		return nil, fmt.Errorf("unsupported data type %d", dataType)
	}
}
