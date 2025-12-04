package tdms

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

const (
	leadInByteLength = 28
)

type File struct {
	r        io.ReadSeekCloser
	segments []*Segment
	mutex    sync.Mutex
}

func OpenFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tdmsFile := &File{
		r: f,
	}
	err = tdmsFile.readMetadata()
	if err != nil {
		return nil, err
	}
	return tdmsFile, nil
}

func (file *File) Close() error {
	return file.r.Close()
}

func (file *File) Segments() []*Segment {
	return file.segments
}

func (file *File) iterateSegments(handler func(segment *Segment) error) error {
	_, err := file.r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	var fileOffset int64
	var nextSegmentOffset int64
	for {
		var segment Segment
		segment.Offset = nextSegmentOffset

		fileOffset, err = file.r.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		segment.Type, segment.LeadIn, err = readLeadIn(file.r)
		if err != nil {
			return err
		}
		if segment.LeadIn.ToC.MetaData() {
			segment.MetaData, err = file.ReadMetaData(segment.LeadIn.ToC)
			if err != nil {
				return err
			}
		}

		_, err := file.r.Seek(fileOffset+int64(segment.LeadIn.RawDataOffset)+leadInByteLength, io.SeekStart)
		if err != nil {
			return err
		}

		err = handler(&segment)
		if err != nil {
			return err
		}

		nextSegmentOffset = segment.Offset + int64(segment.LeadIn.NextSegmentOffset) + leadInByteLength
		if segment.Type == SegmentTypeTDSm {
			_, err = file.r.Seek(nextSegmentOffset, io.SeekStart)
			if err != nil {
				return err
			}
		}
	}
}

func (file *File) readMetadata() error {
	file.mutex.Lock()
	defer file.mutex.Unlock()

	objectMap := make(map[string]*Object)
	err := file.iterateSegments(func(segment *Segment) error {
		fmt.Println()
		for _, object1 := range segment.MetaData.objects {
			object0, ok := objectMap[object1.Path]
			if ok {
				for name, value := range object1.Properties {
					object0.Properties[name] = value
				}
			} else {
				objectMap[object1.Path] = object1
			}
		}
		file.segments = append(file.segments, segment)
		return nil
	})
	if (err != nil) && (err != io.EOF) {
		return err
	}

	var paths []string
	for path := range objectMap {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		object := objectMap[path]
		fmt.Println()
		fmt.Println(object.Path)
		fmt.Println("- Properties:")
		var names []string
		for name := range object.Properties {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("  - %s: %v\n", name, object.Properties[name])
		}
	}

	return nil
}

const (
	NoRawData                     = 0xffffffff
	DAQmxFormatChangingScalerType = 0x00001269
	DAQmxDigitalLineScalerType    = 0x0000126a
)

type Group struct {
}

type Channel struct {
}

func (file *File) ReadMetaData(toc TableOfContents) (*MetaData, error) {
	valueReader := toc.ValueReader()
	// TODO handle NewObjList
	metadata := NewMetaData()
	numberOfObjects, err := valueReader.ReadU32(file.r)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(numberOfObjects); i++ {
		var object Object
		object.Path, err = valueReader.ReadString(file.r)
		if err != nil {
			return nil, err
		}
		rawDataIndexType, err := valueReader.ReadU32(file.r)
		if err != nil {
			return nil, err
		}
		if rawDataIndexType == NoRawData {
			// no raw data assigned
		} else if toc.DAQmxRawData() {
			if rawDataIndexType == DAQmxFormatChangingScalerType {
				var daQmxRawDataIndex DAQmxRawDataIndex
				object.RawDataIndex = &daQmxRawDataIndex
				daQmxRawDataIndex.DataType, err = valueReader.ReadDataType(file.r)
				if err != nil {
					return nil, err
				}
				daQmxRawDataIndex.ArrayDimension, err = valueReader.ReadU32(file.r)
				if err != nil {
					return nil, err
				}
				daQmxRawDataIndex.ChunkSize, err = valueReader.ReadU64(file.r)
				if err != nil {
					return nil, err
				}
				scalerVectorSize, err := valueReader.ReadU32(file.r)
				if err != nil {
					return nil, err
				}
				for i := 0; i < int(scalerVectorSize); i++ {
					scaler, err := readDAQmxFormatChangingScaler(file.r, valueReader)
					if err != nil {
						return nil, err
					}
					daQmxRawDataIndex.Scalers = append(daQmxRawDataIndex.Scalers, scaler)
				}
				rawDataWidthVectorSize, err := valueReader.ReadU32(file.r)
				if err != nil {
					return nil, err
				}
				for i := 0; i < int(rawDataWidthVectorSize); i++ {
					rawDataWidth, err := valueReader.ReadU32(file.r)
					if err != nil {
						return nil, err
					}
					daQmxRawDataIndex.RawDataWidths = append(daQmxRawDataIndex.RawDataWidths, rawDataWidth)
				}
			} else if rawDataIndexType == DAQmxDigitalLineScalerType {
				// TODO
				return nil, fmt.Errorf("unsupported rawDataIndexType: 0x%08x", object.RawDataIndex)
			} else {
				return nil, fmt.Errorf("unsupported rawDataIndexType: 0x%08x", object.RawDataIndex)
			}
		} else {
			// TODO
			return nil, fmt.Errorf("unsupported rawDataIndexType: 0x%08x", object.RawDataIndex)
		}

		object.Properties = make(map[string]any)
		numberOfProperties, err := valueReader.ReadU32(file.r)
		if err != nil {
			return nil, err
		}
		for j := 0; j < int(numberOfProperties); j++ {
			propertyName, err := valueReader.ReadString(file.r)
			if err != nil {
				return nil, err
			}
			propertyValue, err := valueReader.ReadValue(file.r)
			if err != nil {
				return nil, err
			}
			object.Properties[propertyName] = propertyValue
		}

		err = metadata.AddObject(&object)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}
