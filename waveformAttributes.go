package tdms

import (
	"time"

	"github.com/ngyewch/tdms-go/utils"
)

type WaveformAttributes struct {
	StartTime       time.Time
	StartOffset     float64
	Increment       float64
	Samples         int
	Unit            string
	UnitDescription string
}

func GetWaveformAttributes(props map[string]any) (*WaveformAttributes, error) {
	startTime, _, err := utils.GetTime(props, "wf_start_time")
	if err != nil {
		return nil, err
	}
	startOffset, _, err := utils.GetFloat64(props, "wf_start_offset")
	if err != nil {
		return nil, err
	}
	increment, _, err := utils.GetFloat64(props, "wf_increment")
	if err != nil {
		return nil, err
	}
	samples, _, err := utils.GetInt(props, "wf_samples")
	if err != nil {
		return nil, err
	}
	unit, _, err := utils.GetString(props, "unit_string")
	if err != nil {
		return nil, err
	}
	unitDescription, _, err := utils.GetString(props, "NI_UnitDescription")
	if err != nil {
		return nil, err
	}
	return &WaveformAttributes{
		StartTime:       startTime,
		StartOffset:     startOffset,
		Increment:       increment,
		Samples:         samples,
		Unit:            unit,
		UnitDescription: unitDescription,
	}, nil
}
