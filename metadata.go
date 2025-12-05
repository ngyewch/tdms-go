package tdms

import (
	"fmt"
	"io"
)

type MetaData struct {
	objects   []*Object
	objectMap map[string]*Object
}

func NewMetaData() *MetaData {
	return &MetaData{
		objectMap: make(map[string]*Object),
	}
}

func (m *MetaData) Objects() []*Object {
	return m.objects
}

func (m *MetaData) AddObject(object *Object) error {
	obj := m.objectMap[object.Path]
	if obj != nil {
		return fmt.Errorf("object %s already added", object.Path)
	}
	m.objects = append(m.objects, object)
	m.objectMap[object.Path] = object
	return nil
}

func (m *MetaData) GetObjectByPath(path string) *Object {
	return m.objectMap[path]
}

func ReadMetaData(r io.Reader, toc TableOfContents, previousSegment *Segment) (*MetaData, error) {
	valueReader := toc.ValueReader()
	// TODO handle NewObjList
	metadata := NewMetaData()
	numberOfObjects, err := valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(numberOfObjects); i++ {
		var object Object
		object.Path, err = valueReader.ReadString(r)
		if err != nil {
			return nil, err
		}
		rawDataIndexType, err := valueReader.ReadU32(r)
		if err != nil {
			return nil, err
		}
		if rawDataIndexType == RawDataIndexTypeNoRawData {
			// no raw data assigned
		} else if toc.DAQmxRawData() {
			if rawDataIndexType == RawDataIndexTypeDAQmxFormatChangingScalerType {
				object.RawDataIndex, err = ReadDAQmxRawDataIndex(r, valueReader)
				if err != nil {
					return nil, err
				}
			} else if rawDataIndexType == RawDataIndexTypeDAQmxDigitalLineScalerType {
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
		numberOfProperties, err := valueReader.ReadU32(r)
		if err != nil {
			return nil, err
		}
		for j := 0; j < int(numberOfProperties); j++ {
			propertyName, err := valueReader.ReadString(r)
			if err != nil {
				return nil, err
			}
			propertyValue, err := valueReader.ReadValue(r)
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
