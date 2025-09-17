package tdms

import (
	"fmt"
	"io"
)

const (
	leadInByteLength = 28
)

type Reader struct {
	r io.ReadSeeker
}

func NewReader(r io.ReadSeeker) *Reader {
	return &Reader{
		r: r,
	}
}

type Segment struct {
	Type     SegmentType
	LeadIn   *LeadIn
	MetaData *MetaData
}

func (reader *Reader) NextSegment() (*Segment, error) {
	var segment Segment
	var err error

	segmentOffset, err := reader.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	fmt.Printf("segment offset: 0x%x (%d)\n", segmentOffset, segmentOffset)

	segment.Type, segment.LeadIn, err = readLeadIn(reader.r)
	if err != nil {
		return nil, err
	}
	if segment.LeadIn.ToC.MetaData() {
		segment.MetaData, err = reader.ReadMetaData(segment.LeadIn.ToC)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO
	}

	pos, err := reader.r.Seek(segmentOffset+int64(segment.LeadIn.RawDataOffset)+leadInByteLength, io.SeekStart)
	if err != nil {
		return nil, err
	}
	fmt.Printf("pos: 0x%x (%d)\n", pos, pos)

	_, err = reader.r.Seek(segmentOffset+int64(segment.LeadIn.NextSegmentOffset)+leadInByteLength, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &segment, nil
}

const (
	NoRawData                     = 0xffffffff
	DAQmxFormatChangingScalerType = 0x00001269
	DAQmxDigitalLineScalerType    = 0x0000126a
)

type MetaData struct {
	Objects []*Object
}

type Object struct {
	Path         string
	RawDataIndex uint32
	V            any
	// TODO
	Properties map[string]any
}

type DAQmxRawDataIndex struct {
	DataType       DataType
	ArrayDimension uint32
	ChunkSize      uint64
	Scalers        []*DAQmxFormatChangingScaler
	RawDataWidths  []uint32
}

func (reader *Reader) ReadMetaData(toc TableOfContents) (*MetaData, error) {
	valueReader := toc.ValueReader()
	// TODO handle NewObjList
	var metadata MetaData
	numberOfObjects, err := valueReader.ReadU32(reader.r)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(numberOfObjects); i++ {
		var object Object
		object.Path, err = valueReader.ReadString(reader.r)
		if err != nil {
			return nil, err
		}
		object.RawDataIndex, err = valueReader.ReadU32(reader.r)
		if err != nil {
			return nil, err
		}
		if object.RawDataIndex == NoRawData {
			// no raw data assigned
		} else if toc.DAQmxRawData() {
			if object.RawDataIndex == DAQmxFormatChangingScalerType {
				var daQmxRawDataIndex DAQmxRawDataIndex
				object.V = &daQmxRawDataIndex

				daQmxRawDataIndex.DataType, err = valueReader.ReadDataType(reader.r)
				if err != nil {
					return nil, err
				}
				daQmxRawDataIndex.ArrayDimension, err = valueReader.ReadU32(reader.r)
				if err != nil {
					return nil, err
				}
				daQmxRawDataIndex.ChunkSize, err = valueReader.ReadU64(reader.r)
				if err != nil {
					return nil, err
				}
				scalerVectorSize, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return nil, err
				}
				for i := 0; i < int(scalerVectorSize); i++ {
					scaler, err := readDAQmxFormatChangingScaler(reader.r, valueReader)
					if err != nil {
						return nil, err
					}
					daQmxRawDataIndex.Scalers = append(daQmxRawDataIndex.Scalers, scaler)
				}
				rawDataWidthVectorSize, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return nil, err
				}
				for i := 0; i < int(rawDataWidthVectorSize); i++ {
					rawDataWidth, err := valueReader.ReadU32(reader.r)
					if err != nil {
						return nil, err
					}
					daQmxRawDataIndex.RawDataWidths = append(daQmxRawDataIndex.RawDataWidths, rawDataWidth)
				}
			} else if object.RawDataIndex == DAQmxDigitalLineScalerType {
				// TODO
				return nil, fmt.Errorf("unsupported rawDataIndex: 0x%08x", object.RawDataIndex)
			} else {
				return nil, fmt.Errorf("unsupported rawDataIndex: 0x%08x", object.RawDataIndex)
			}
		} else {
			// TODO
		}

		object.Properties = make(map[string]any)
		numberOfProperties, err := valueReader.ReadU32(reader.r)
		if err != nil {
			return nil, err
		}
		for j := 0; j < int(numberOfProperties); j++ {
			propertyName, err := valueReader.ReadString(reader.r)
			if err != nil {
				return nil, err
			}
			propertyValue, err := valueReader.ReadValue(reader.r)
			if err != nil {
				return nil, err
			}
			object.Properties[propertyName] = propertyValue
		}

		metadata.Objects = append(metadata.Objects, &object)
	}

	return &metadata, nil
}
