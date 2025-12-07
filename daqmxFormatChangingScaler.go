package tdms

import (
	"fmt"
	"io"
)

type DAQmxFormatChangingScaler struct {
	DataType                     DataType
	RawBufferIndex               uint32
	RawByteOffsetWithinTheStride uint32
	SampleFormatBitmap           uint32
	ScaleId                      uint32
}

func ReadDAQmxFormatChangingScaler(r io.Reader, valueReader *ValueReader) (*DAQmxFormatChangingScaler, error) {
	var scaler DAQmxFormatChangingScaler
	var err error
	scaler.DataType, err = valueReader.ReadDAQmxDataType(r)
	if err != nil {
		return nil, err
	}
	scaler.RawBufferIndex, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.RawByteOffsetWithinTheStride, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.SampleFormatBitmap, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.ScaleId, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	return &scaler, nil
}

func (scaler *DAQmxFormatChangingScaler) Scale(v any) (float64, error) {
	return 0, fmt.Errorf("DAQmxFormatChangingScaler.Scale not implemented")
}
