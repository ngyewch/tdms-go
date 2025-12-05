package tdms

import (
	"fmt"
	"io"
	"os"
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

	var root *Node
	err := file.iterateSegments(func(segment *Segment) error {
		for _, object := range segment.MetaData.objects {
			objectPath, err := ObjectPathFromString(object.Path)
			if err != nil {
				return err
			}
			if objectPath.IsRoot() {
				if root == nil {
					root = NewNode("", object.Path)
				}
				for name, value := range object.Properties {
					root.Properties().Insert(name, value)
				}
				continue
			}
			if objectPath.Group == "" {
				return fmt.Errorf("group name is empty")
			}
			if root == nil {
				return fmt.Errorf("root not defined")
			}
			if objectPath.IsGroup() {
				group := root.GetChildByName(objectPath.Group)
				if group == nil {
					group = NewNode(objectPath.Group, object.Path)
					root.AddChild(group)
				}
				for name, value := range object.Properties {
					group.Properties().Insert(name, value)
				}
			} else if objectPath.IsChannel() {
				group := root.GetChildByName(objectPath.Group)
				if group == nil {
					return fmt.Errorf("group not defined")
				}
				channel := group.GetChildByName(objectPath.Channel)
				if channel == nil {
					channel = NewNode(objectPath.Channel, object.Path)
					group.AddChild(channel)
				}
				for name, value := range object.Properties {
					channel.Properties().Insert(name, value)
				}
			}
		}
		file.segments = append(file.segments, segment)
		return nil
	})
	if (err != nil) && (err != io.EOF) {
		return err
	}

	for _, group := range root.Children() {
		fmt.Printf("%s\n", group.Name())
		for _, channel := range group.Children() {
			fmt.Printf("- %s\n", channel.Name())
			seq := channel.Properties().All()
			seq(func(name string, value any) bool {
				fmt.Printf("  - %s: %v\n", name, value)
				return true
			})
		}
	}

	return nil
}
