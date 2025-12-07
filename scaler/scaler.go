package scaler

import (
	"fmt"
	"strings"

	"github.com/ngyewch/tdms-go/utils"
)

type ScalingStatus string

const (
	ScalingStatusUnscaled ScalingStatus = "unscaled"
	ScalingStatusScaled   ScalingStatus = "scaled"
)

type Scaler interface {
	Scale(v any) float64
}

func GetScalers(props map[string]any) ([]Scaler, error) {
	numberOfScales, hasNumberOfScales, err := utils.GetInt(props, "NI_Number_Of_Scales")
	if err != nil {
		return nil, err
	}
	if !hasNumberOfScales || (numberOfScales <= 0) {
		return nil, nil
	}

	sScalingStatus, hasScalingStatus, err := utils.GetString(props, "NI_Scaling_Status")
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
		scaleType, hasScaleType, err := utils.GetString(scalerProps, "Scale_Type")
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
