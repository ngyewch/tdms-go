package tdms

import (
	"fmt"
	"io"
)

type Reader struct {
	r io.ReadSeeker
}

func NewReader(r io.ReadSeeker) *Reader {
	return &Reader{
		r: r,
	}
}

type Segment struct{}

func (reader *Reader) NextSegment() (*Segment, error) {
	segmentType, leadIn, err := readLeadIn(reader.r)
	if err != nil {
		return nil, err
	}
	fmt.Printf("segmentType: %v, leadIn: %v\n", segmentType, leadIn)
	if leadIn.ToC.MetaData() {
		err = reader.ReadMetaData(leadIn.ToC)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO
	}
	return &Segment{}, nil
}

const (
	NoRawData                 = 0xffffffff
	DAQmxFormatChangingScaler = 0x00001269
	DAQmxDigitalLineScaler    = 0x0000126a
)

func (reader *Reader) ReadMetaData(toc TableOfContents) error {
	valueReader := toc.ValueReader()

	// TODO handle NewObjList

	numberOfObjects, err := valueReader.ReadU32(reader.r)
	if err != nil {
		return err
	}
	fmt.Println("----")
	for i := 0; i < int(numberOfObjects); i++ {
		objectPath, err := valueReader.ReadString(reader.r)
		if err != nil {
			return err
		}
		fmt.Printf("objectPath: %s\n", objectPath)
		rawDataIndex, err := valueReader.ReadU32(reader.r)
		if err != nil {
			return err
		}
		fmt.Printf("- rawDataIndex: %d\n", rawDataIndex)
		if rawDataIndex == NoRawData {
			// no raw data assigned
		} else if toc.DAQmxRawData() {
			if rawDataIndex == DAQmxFormatChangingScaler {
				// TODO
			} else if rawDataIndex == DAQmxDigitalLineScaler {
				// TODO
			} else {
				return fmt.Errorf("invalid rawDataIndex: 0x%08x", rawDataIndex)
			}

			dataType, err := valueReader.ReadDataType(reader.r)
			if err != nil {
				return err
			}
			fmt.Printf("dataType: %v\n", dataType)
			arrayDimension, err := valueReader.ReadU32(reader.r)
			if err != nil {
				return err
			}
			fmt.Printf("arrayDimension: %v\n", arrayDimension)
			chunkSize, err := valueReader.ReadU64(reader.r)
			if err != nil {
				return err
			}
			fmt.Printf("chunkSize: %v\n", chunkSize)

			scalerVectorSize, err := valueReader.ReadU32(reader.r)
			if err != nil {
				return err
			}

			for i := 0; i < int(scalerVectorSize); i++ {
				daqMxDataType, err := valueReader.ReadDataType(reader.r)
				if err != nil {
					return err
				}
				fmt.Printf("daqMxDataType: %v\n", daqMxDataType)
				rawBufferIndex, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return err
				}
				fmt.Printf("rawBufferIndex: %v\n", rawBufferIndex)
				rawByteOffsetWithinTheStride, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return err
				}
				fmt.Printf("rawByteOffsetWithinTheStride: %v\n", rawByteOffsetWithinTheStride)
				sampleFormatBitmap, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return err
				}
				fmt.Printf("sampleFormatBitmap: %v\n", sampleFormatBitmap)
				scaleId, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return err
				}
				fmt.Printf("scaleId: %v\n", scaleId)

				rawDataVectorSize, err := valueReader.ReadU32(reader.r)
				if err != nil {
					return err
				}
				var rawData []uint32
				for i := 0; i < int(rawDataVectorSize); i++ {
					v, err := valueReader.ReadU32(reader.r)
					if err != nil {
						return err
					}
					rawData = append(rawData, v)
				}
			}
		} else {
			// TODO
		}

		numberOfProperties, err := valueReader.ReadU32(reader.r)
		if numberOfProperties > 0 {
			fmt.Println("- properties:")
			for j := 0; j < int(numberOfProperties); j++ {
				propertyName, err := valueReader.ReadString(reader.r)
				if err != nil {
					return err
				}
				propertyValue, err := valueReader.ReadValue(reader.r)
				if err != nil {
					return err
				}
				fmt.Printf("  - %s: %v\n", propertyName, propertyValue)
			}
		}
	}

	return nil
}
