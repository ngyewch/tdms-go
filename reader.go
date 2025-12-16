package tdms

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gosimple/slug"
	"github.com/ngyewch/tdms-go/utils"
	"github.com/samber/oops"
)

const (
	leadInByteLength = 28
)

type File struct {
	r       io.ReadSeekCloser
	root    *Node
	nodeMap map[string]*Node
	mutex   sync.Mutex
}

func OpenFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tdmsFile := &File{
		r:       f,
		nodeMap: make(map[string]*Node),
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

func (file *File) Root() *Node {
	return file.root
}

func (file *File) Node(path string) *Node {
	return file.nodeMap[path]
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
			if err == io.EOF {
				return err
			}
			return oops.
				With("segmentOffset", fileOffset).
				Wrapf(err, "invalid segment")
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
					file.nodeMap[object.Path] = root
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
					file.nodeMap[object.Path] = group
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
					file.nodeMap[object.Path] = channel
					group.AddChild(channel)
				}
				for name, value := range object.Properties {
					channel.Properties().Insert(name, value)
				}
			}
		}
		return nil
	})
	if (err != nil) && (err != io.EOF) {
		return err
	}

	file.root = root

	return nil
}

func (file *File) ReadData() error {
	err := file.iterateSegments(func(segment *Segment) error {
		if !segment.LeadIn.ToC.RawData() {
			return nil
		}
		if segment.MetaData == nil {
			return nil
		}

		if segment.LeadIn.ToC.DAQmxRawData() {
			type Channel struct {
				slug               string
				object             *Object
				rawDataIndex       *DAQmxRawDataIndex
				waveformAttributes *WaveformAttributes
				sampleRate         int
			}

			var channels []Channel
			var rawDataIndexes []*DAQmxRawDataIndex
			for _, object := range segment.MetaData.Objects() {
				if object.RawDataIndex != nil {
					daqmxRawDataIndex := object.RawDataIndex.(*DAQmxRawDataIndex)
					if daqmxRawDataIndex == nil {
						return fmt.Errorf("DAQmx raw data index expected")
					}
					if len(daqmxRawDataIndex.Scalers) <= 0 {
						return fmt.Errorf("no scalers defined")
					}
					daqmxFormatChangingScaler := daqmxRawDataIndex.Scalers[0].(*DAQmxFormatChangingScaler)
					if daqmxFormatChangingScaler == nil {
						return fmt.Errorf("DAQmx format changing scaler expected as first scaler")
					}
					if len(channels) > 0 {
						err := channels[0].rawDataIndex.CheckCompatibility(daqmxRawDataIndex)
						if err != nil {
							return err
						}
					}
					node := file.Node(object.Path)
					if node == nil {
						return fmt.Errorf("could not read waveform attributes")
					}
					waveformAttributes, err := GetWaveformAttributes(node.Properties().Collect())
					if err != nil {
						return err
					}
					channels = append(channels, Channel{
						slug:               slug.Make(object.Path),
						object:             object,
						rawDataIndex:       daqmxRawDataIndex,
						waveformAttributes: waveformAttributes,
						sampleRate:         int(math.Ceil(1 / waveformAttributes.Increment)),
					})
					rawDataIndexes = append(rawDataIndexes, daqmxRawDataIndex)
				}
			}
			if len(channels) > 0 {
				var wavEncoders []*wav.Encoder
				err := os.MkdirAll("output", 0755)
				if err != nil {
					return err
				}
				for _, channel := range channels {
					f, err := os.Create(filepath.Join("output", channel.slug+".wav"))
					if err != nil {
						return err
					}
					defer func(f *os.File) {
						_ = f.Close()
					}(f)
					wavEncoder := wav.NewEncoder(f, channel.sampleRate, 16, 1, 1)
					defer func(wavEncoder *wav.Encoder) {
						_ = wavEncoder.Close()
					}(wavEncoder)
					wavEncoders = append(wavEncoders, wavEncoder)
				}

				rawDataWidths := channels[0].rawDataIndex.RawDataWidths
				buffers := make([][]byte, len(rawDataWidths))
				var totalRawDataWidth uint32
				for i, rawDataWidth := range rawDataWidths {
					buffers[i] = make([]byte, rawDataWidth)
					totalRawDataWidth += rawDataWidth
				}

				valueReader := segment.LeadIn.ToC.ValueReader()
				rawDataSize := segment.LeadIn.NextSegmentOffset - segment.LeadIn.RawDataOffset
				sampleCount := int(rawDataSize / uint64(totalRawDataWidth))
				sampleList := make([][]float64, len(channels))
				for i := range sampleList {
					sampleList[i] = make([]float64, sampleCount)
				}

				for i := 0; i < sampleCount; i++ {
					for bufferNo := 0; bufferNo < len(buffers); bufferNo++ {
						_, err := file.r.Read(buffers[bufferNo])
						if err != nil {
							return err
						}
					}
					for channelNo, channel := range channels {
						firstScaler := channel.rawDataIndex.Scalers[0].(*DAQmxFormatChangingScaler)
						if firstScaler == nil {
							return fmt.Errorf("DAQmx format changing scaler expected as first scaler")
						}
						v0, err := firstScaler.ReadFromBuffer(valueReader, buffers)
						if err != nil {
							return err
						}
						v, err := utils.AsFloat64(v0)
						if err != nil {
							return err
						}
						for _, scaler := range channel.rawDataIndex.Scalers[1:] {
							v, err = scaler.Scale(v)
							if err != nil {
								return err
							}
						}
						sampleList[channelNo][i] = v
					}
				}

				for channelNo := range channels {
					samples := sampleList[channelNo]
					wavEncoder := wavEncoders[channelNo]
					buffer := &audio.IntBuffer{
						Format: &audio.Format{
							NumChannels: wavEncoder.NumChans,
							SampleRate:  wavEncoder.SampleRate,
						},
						Data:           make([]int, sampleCount),
						SourceBitDepth: wavEncoder.BitDepth,
					}
					multiplier := math.Pow(2, float64(wavEncoder.BitDepth-1)) - 1
					maxValue := float64(5)
					for i, sample := range samples {
						if sample < -maxValue {
							sample = -maxValue
						} else if sample > maxValue {
							sample = maxValue
						}
						buffer.Data[i] = int(math.Round(sample * multiplier / maxValue))
					}
					err = wavEncoder.Write(buffer)
					if err != nil {
						return err
					}
				}
			}
		} else {
			// TODO
			return fmt.Errorf("not supported yet")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
