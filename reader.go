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
	var previousSegment *Segment
	for {
		fileOffset, err = file.r.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		segment, err := ReadSegment(file.r, previousSegment)
		if err != nil {
			return err
		}
		segment.Offset = nextSegmentOffset

		_, err = file.r.Seek(fileOffset+int64(segment.LeadIn.RawDataOffset)+leadInByteLength, io.SeekStart)
		if err != nil {
			return err
		}

		err = handler(segment)
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

		previousSegment = segment
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

type Group struct {
}

type Channel struct {
}
