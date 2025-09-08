package tdms

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	ErrTDSMTagNotFound = errors.New("TDSM tag not found")
)

type Reader struct {
	r io.Reader
}

type LeadIn struct {
	ToC               TableOfContents
	VersionNumber     uint32
	NextSegmentOffset uint64
	RawDataOffset     uint64
}

type TableOfContents struct {
	MetaData        bool
	RawData         bool
	DAQmxRawData    bool
	InterleavedData bool
	BigEndian       bool
	NewObjList      bool
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

func (reader *Reader) ReadLeadIn() (*LeadIn, error) {
	tdsmTag := make([]byte, 4)
	_, err := io.ReadFull(reader.r, tdsmTag)
	if err != nil {
		return nil, err
	}
	if bytes.Compare(tdsmTag, []byte("TDSm")) != 0 {
		return nil, ErrTDSMTagNotFound
	}

	leadIn := LeadIn{
		ToC: TableOfContents{},
	}

	var tocMask uint32
	err = binary.Read(reader.r, binary.LittleEndian, &tocMask)
	if err != nil {
		return nil, err
	}
	fmt.Printf("tocMask: %+v\n", tocMask)
	leadIn.ToC.MetaData = tocMask&(1<<1) != 0
	leadIn.ToC.RawData = tocMask&(1<<3) != 0
	leadIn.ToC.DAQmxRawData = tocMask&(1<<7) != 0
	leadIn.ToC.InterleavedData = tocMask&(1<<5) != 0
	leadIn.ToC.BigEndian = tocMask&(1<<6) != 0
	leadIn.ToC.NewObjList = tocMask&(1<<2) != 0

	var byteOrder binary.ByteOrder = binary.LittleEndian
	if leadIn.ToC.BigEndian {
		byteOrder = binary.BigEndian
	}

	err = binary.Read(reader.r, byteOrder, &leadIn.VersionNumber)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader.r, byteOrder, &leadIn.NextSegmentOffset)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader.r, byteOrder, &leadIn.RawDataOffset)
	if err != nil {
		return nil, err
	}

	return &leadIn, nil
}

func (reader *Reader) ReadMetaData(leadIn *LeadIn) error {
	var byteOrder binary.ByteOrder = binary.LittleEndian
	if leadIn.ToC.BigEndian {
		byteOrder = binary.BigEndian
	}

	r := io.LimitReader(reader.r, int64(leadIn.RawDataOffset))
	var numberOfObjects uint32
	err := binary.Read(r, byteOrder, &numberOfObjects)
	if err != nil {
		return err
	}

	fmt.Println("----")
	for i := 0; i < int(numberOfObjects); i++ {
		objectPath, err := readString(r, byteOrder)
		if err != nil {
			return err
		}
		fmt.Printf("objectPath: %s\n", objectPath)
		var rawDataIndex uint32
		err = binary.Read(r, byteOrder, &rawDataIndex)
		if err != nil {
			return err
		}
		fmt.Printf("- rawDataIndex: %d\n", rawDataIndex)
		if rawDataIndex == 0xffffffff {
			// no raw data assigned
		} else if leadIn.ToC.DAQmxRawData {
			if rawDataIndex == 0x1269 { // DAQmx Format Changing scaler
				// TODO
			} else if rawDataIndex == 0x126a { // DAQmx Digital Line scaler
				// TODO
			}
			var dataType DataType
			err = binary.Read(r, byteOrder, &dataType)
			if err != nil {
				return err
			}
			fmt.Println("dataType", dataType)
			var arrayDimension uint32
			err = binary.Read(r, byteOrder, &arrayDimension)
			if err != nil {
				return err
			}
			var numberOfValues uint64
			err = binary.Read(r, byteOrder, &numberOfValues)
			if err != nil {
				return err
			}
			if rawDataIndex == 0x1269 {
				var vectorSize uint32
				err = binary.Read(r, byteOrder, &vectorSize)
				if err != nil {
					return err
				}
				for j := 0; j < int(vectorSize); j++ {
					var daqmxDataType DataType
					err = binary.Read(r, byteOrder, &daqmxDataType)
					if err != nil {
						return err
					}
					fmt.Println("daqmxDataType", daqmxDataType)
					var rawBufferIndex uint32
					err = binary.Read(r, byteOrder, &rawBufferIndex)
					if err != nil {
						return err
					}
					var rawByteOffsetWithinTheStride uint32
					err = binary.Read(r, byteOrder, &rawByteOffsetWithinTheStride)
					if err != nil {
						return err
					}
					var sampleFormatBitmap uint32
					err = binary.Read(r, byteOrder, &sampleFormatBitmap)
					if err != nil {
						return err
					}
					var scaleId uint32
					err = binary.Read(r, byteOrder, &scaleId)
					if err != nil {
						return err
					}
				}
			}
			var vectorSize uint32
			for j := 0; j < int(vectorSize); j++ {
				var v uint32 // TODO
				err = binary.Read(r, byteOrder, &v)
				if err != nil {
					return err
				}
			}
		} else {
			// TODO
		}

		var numberOfProperties uint32
		err = binary.Read(r, byteOrder, &numberOfProperties)
		if numberOfProperties > 0 {
			fmt.Println("- properties:")
			for j := 0; j < int(numberOfProperties); j++ {
				propertyName, err := readString(r, byteOrder)
				if err != nil {
					return err
				}
				propertyValue, err := readValue(r, byteOrder)
				if err != nil {
					return err
				}
				fmt.Printf("  - %s: %v\n", propertyName, propertyValue)
			}
		}
	}

	return nil
}

func readString(r io.Reader, byteOrder binary.ByteOrder) (string, error) {
	var stringLength uint32
	err := binary.Read(r, byteOrder, &stringLength)
	if err != nil {
		return "", err
	}
	buf := make([]byte, stringLength)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

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

func readValue(r io.Reader, byteOrder binary.ByteOrder) (any, error) {
	var dataType DataType
	err := binary.Read(r, byteOrder, &dataType)
	if err != nil {
		return nil, err
	}
	switch dataType {
	case DataTypeString:
		return readString(r, byteOrder)
	default:
		return nil, fmt.Errorf("unsupported data type %d", dataType)
	}
}
