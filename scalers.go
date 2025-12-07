package tdms

import (
	"fmt"
	"math"
	"strings"
)

type ScalingStatus string

const (
	ScalingStatusUnscaled ScalingStatus = "unscaled"
	ScalingStatusScaled   ScalingStatus = "scaled"
)

type Scaler interface {
	Scale(v any) float64
}

type LinearScaler struct {
	LinearInputSource uint
	LinearSlope       float64
	LinearYIntercept  float64
}

func NewLinearScaler(props map[string]any) (*LinearScaler, error) {
	linearInputSource, hasLinearInputSource, err := getUint(props, "Linear_Input_Source")
	if err != nil {
		return nil, err
	}
	if !hasLinearInputSource {
		return nil, fmt.Errorf("Line_Input_Source not specified")
	}
	linearSlope, hasLinearSlope, err := getFloat64(props, "Linear_Slope")
	if !hasLinearSlope {
		return nil, fmt.Errorf("Linear_Slope not specified")
	}
	linearYIntercept, hasLinearYIntercept, err := getFloat64(props, "Linear_Y_Intercept")
	if !hasLinearYIntercept {
		return nil, fmt.Errorf("Linear_Y_Intercept not specified")
	}
	return &LinearScaler{
		LinearInputSource: linearInputSource,
		LinearSlope:       linearSlope,
		LinearYIntercept:  linearYIntercept,
	}, nil
}

func (scaler *LinearScaler) Scale(v any) float64 {
	n, err := asFloat64(v)
	if err != nil {
		return math.NaN()
	}
	return (n * scaler.LinearSlope) + scaler.LinearYIntercept
}

func GetScalers(props map[string]any) ([]Scaler, error) {
	numberOfScales, hasNumberOfScales, err := getInt(props, "NI_Number_Of_Scales")
	if err != nil {
		return nil, err
	}
	if !hasNumberOfScales || (numberOfScales <= 0) {
		return nil, nil
	}

	sScalingStatus, hasScalingStatus, err := getString(props, "NI_Scaling_Status")
	if err != nil {
		return nil, err
	}
	if !hasScalingStatus {
		return nil, fmt.Errorf("scaling status not specified")
	}
	scalingStatus := ScalingStatus(sScalingStatus)
	switch scalingStatus {
	case ScalingStatusUnscaled:
		// continue
	case ScalingStatusScaled:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown scaling status: %v", scalingStatus)
	}

	scalers := make([]Scaler, numberOfScales)
	for i := 0; i < numberOfScales; i++ {
		scalerProps := make(map[string]any)
		prefix := fmt.Sprintf("NI_Scale[%d]_", i)
		for name, value := range props {
			if strings.HasPrefix(name, prefix) {
				scalerProps[name[len(prefix):]] = value
			}
		}
		scaleType, hasScaleType, err := getString(scalerProps, "Scale_Type")
		if err != nil {
			return nil, err
		}
		if hasScaleType {
			switch scaleType {
			case "Linear":
				scaler, err := NewLinearScaler(scalerProps)
				if err != nil {
					return nil, err
				}
				scalers[i] = scaler
			default:
				return nil, fmt.Errorf("unknown scale type '%v'", scaleType)
			}
		}
	}

	return scalers, nil
}
