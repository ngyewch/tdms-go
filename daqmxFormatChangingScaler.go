package tdms

import (
	"bytes"
	"fmt"
	"io"
)

type DAQmxFormatChangingScaler struct {
	dataType                     DataType
	rawBufferIndex               uint32
	rawByteOffsetWithinTheStride uint32
	sampleFormatBitmap           uint32
	scaleId                      uint32
}

func ReadDAQmxFormatChangingScaler(r io.Reader, valueReader *ValueReader) (*DAQmxFormatChangingScaler, error) {
	var scaler DAQmxFormatChangingScaler
	var err error
	scaler.dataType, err = valueReader.ReadDAQmxDataType(r)
	if err != nil {
		return nil, err
	}
	scaler.rawBufferIndex, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.rawByteOffsetWithinTheStride, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.sampleFormatBitmap, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.scaleId, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	return &scaler, nil
}

func (scaler *DAQmxFormatChangingScaler) ScaleId() uint32 {
	return scaler.scaleId
}

func (scaler *DAQmxFormatChangingScaler) Scale(v any) (float64, error) {
	return 0, fmt.Errorf("DAQmxFormatChangingScaler.Scale not implemented")
}

func (scaler *DAQmxFormatChangingScaler) ReadFromBuffer(vr *ValueReader, buffers [][]byte) (any, error) {
	if scaler.rawBufferIndex >= uint32(len(buffers)) {
		return nil, fmt.Errorf("buffer index out of range")
	}
	buffer := buffers[scaler.rawBufferIndex]
	startOffset := int(scaler.rawByteOffsetWithinTheStride)
	endOffset := startOffset + scaler.dataType.SizeInBytes()
	if (startOffset < 0) || (startOffset >= len(buffer)) {
		return nil, fmt.Errorf("start byte offset out of range")
	}
	if (endOffset < 0) || (endOffset > len(buffer)) {
		return nil, fmt.Errorf("end byte offset out of range")
	}
	b := buffer[startOffset:endOffset]
	return vr.ReadValueForDataType(bytes.NewReader(b), scaler.dataType)
}
