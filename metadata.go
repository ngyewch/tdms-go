package tdms

import "fmt"

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
