package tdms

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/egregors/sortedmap"
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

	root := NewRoot()
	for path := range objectMap {
		objectPath, err := ObjectPathFromString(path)
		if err != nil {
			return err
		}
		if objectPath.IsRoot() || (objectPath.Group == "") {
			continue
		}
		group := root.Group(objectPath.Group)
		if group == nil {
			group = NewGroup(objectPath.Group)
			root.AddGroup(group)
		}
		if objectPath.Channel != "" {
			channel := group.Channel(objectPath.Channel)
			if channel == nil {
				channel = NewChannel(objectPath.Channel)
				group.AddChannel(channel)
			}
		}
	}
	for _, group := range root.Groups() {
		fmt.Printf("%s\n", group.Name())
		for _, channel := range group.Channels() {
			fmt.Printf("  %s\n", channel.Name())
		}
	}

	if false {
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
	}

	return nil
}

type Node[T any] struct {
	name     string
	path     string
	childMap *sortedmap.SortedMap[map[string]T, string, T]
}

func NewNode[T any](name string, path string) *Node[T] {
	return &Node[T]{
		name: name,
		path: path,
		childMap: sortedmap.New[map[string]T, string, T](func(i, j sortedmap.KV[string, T]) bool {
			return i.Key < j.Key
		}),
	}
}

func (node *Node[T]) Name() string {
	return node.name
}

func (node *Node[T]) Path(name string) string {
	return node.path
}

func (node *Node[T]) Children() []T {
	return node.childMap.CollectValues()
}

func (node *Node[T]) GetChildByName(name string) T {
	child, _ := node.childMap.Get(name)
	return child
}

func (node *Node[T]) AddChild(name string, child T) {
	node.childMap.Insert(name, child)
}

type Root struct {
	groupMap *sortedmap.SortedMap[map[string]*Group, string, *Group]
}

func NewRoot() *Root {
	return &Root{
		groupMap: sortedmap.New[map[string]*Group, string, *Group](func(i, j sortedmap.KV[string, *Group]) bool {
			return i.Key < j.Key
		}),
	}
}

func (root *Root) Groups() []*Group {
	return root.groupMap.CollectValues()
}

func (root *Root) Group(name string) *Group {
	group, _ := root.groupMap.Get(name)
	return group
}

func (root *Root) AddGroup(group *Group) {
	root.groupMap.Insert(group.Name(), group)
}

type Group struct {
	object     *Object
	name       string
	channelMap *sortedmap.SortedMap[map[string]*Channel, string, *Channel]
}

func NewGroup(name string) *Group {
	return &Group{
		name: name,
		channelMap: sortedmap.New[map[string]*Channel, string, *Channel](func(i, j sortedmap.KV[string, *Channel]) bool {
			return i.Key < j.Key
		}),
	}
}

func (group *Group) Name() string {
	return group.name
}

func (group *Group) Channels() []*Channel {
	return group.channelMap.CollectValues()
}

func (group *Group) Channel(name string) *Channel {
	channel, _ := group.channelMap.Get(name)
	return channel
}

func (group *Group) AddChannel(channel *Channel) {
	group.channelMap.Insert(channel.Name(), channel)
}

type Channel struct {
	name       string
	properties map[string]any
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:       name,
		properties: make(map[string]any),
	}
}

func (channel *Channel) Name() string {
	return channel.name
}

func (channel *Channel) Properties() map[string]any {
	return channel.properties
}
