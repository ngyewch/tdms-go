package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gosimple/slug"
	"github.com/ngyewch/tdms-go"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.StringArg(inputFileArg.Name)

	tdmsFile, err := tdms.OpenFile(inputFile)
	if err != nil {
		return err
	}
	defer func(tdmsFile *tdms.File) {
		_ = tdmsFile.Close()
	}(tdmsFile)

	sampleCount, err := tdmsFile.GetSampleCount()
	if err != nil {
		return err
	}
	fmt.Printf("sampleCount: %d\n", sampleCount)

	if false {
		fmt.Println("/")
		tdmsFile.Root().Properties().All()(func(name string, value any) bool {
			fmt.Printf("  * %s: %v [%T]\n", name, value, value)
			return true
		})
		for _, group := range tdmsFile.Root().Children() {
			fmt.Printf("- %s\n", group.Name())
			group.Properties().All()(func(name string, value any) bool {
				fmt.Printf("    * %s: %v [%T]\n", name, value, value)
				return true
			})
			for _, channel := range group.Children() {
				fmt.Printf("  - %s\n", channel.Name())
				channel.Properties().All()(func(name string, value any) bool {
					fmt.Printf("      * %s: %v [%T]\n", name, value, value)
					return true
				})
			}
		}
	}

	err = writeToWav(tdmsFile)
	if err != nil {
		return err
	}

	return nil
}

func writeToWav(tdmsFile *tdms.File) error {
	initialized := false
	var files []*os.File
	var wavEncoders []*wav.Encoder

	defer func() {
		for _, f := range files {
			_ = f.Close()
		}
	}()
	defer func() {
		for _, wavEncoder := range wavEncoders {
			_ = wavEncoder.Close()
		}
	}()

	err := os.MkdirAll("output", 0755)
	if err != nil {
		return err
	}

	const wavBitDepth = 16
	const wavNumChannels = 1
	const wavAudioFormat = 1

	multiplier := math.Pow(2, float64(wavBitDepth-1)) - 1
	maxValue := float64(5)
	minValue := float64(-5)
	divisor := max(math.Abs(maxValue), math.Abs(minValue))

	err = tdmsFile.ReadData(func(chunk tdms.Chunk) error {
		if !initialized {
			files = make([]*os.File, len(chunk.Channels))
			wavEncoders = make([]*wav.Encoder, len(chunk.Channels))

			var err error
			for i, channel := range chunk.Channels {
				files[i], err = os.Create(filepath.Join("output", slug.Make(channel.Node.Path())+".wav"))
				if err != nil {
					return err
				}
				wavEncoders[i] = wav.NewEncoder(files[i], int(math.Round(channel.WaveformAttributes.SampleRate())), wavBitDepth, wavNumChannels, wavAudioFormat)
			}

			initialized = true
		}

		for channelNo, channel := range chunk.Channels {
			wavEncoder := wavEncoders[channelNo]
			buffer := &audio.IntBuffer{
				Format: &audio.Format{
					NumChannels: wavEncoder.NumChans,
					SampleRate:  wavEncoder.SampleRate,
				},
				Data:           make([]int, len(channel.Samples)),
				SourceBitDepth: wavEncoder.BitDepth,
			}

			for i, sample := range channel.Samples {
				if sample < minValue {
					sample = minValue
				} else if sample > maxValue {
					sample = maxValue
				}
				buffer.Data[i] = int(math.Round(sample * multiplier / divisor))
			}

			err = wavEncoder.Write(buffer)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
