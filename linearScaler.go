package tdms

import (
	"fmt"
	"math"

	"github.com/ngyewch/tdms-go/utils"
)

type LinearScaler struct {
	LinearInputSource uint
	LinearSlope       float64
	LinearYIntercept  float64
}

func NewLinearScaler(props map[string]any) (*LinearScaler, error) {
	linearInputSource, hasLinearInputSource, err := utils.GetUint(props, "Linear_Input_Source")
	if err != nil {
		return nil, err
	}
	if !hasLinearInputSource {
		return nil, fmt.Errorf("Line_Input_Source not specified")
	}
	linearSlope, hasLinearSlope, err := utils.GetFloat64(props, "Linear_Slope")
	if !hasLinearSlope {
		return nil, fmt.Errorf("Linear_Slope not specified")
	}
	linearYIntercept, hasLinearYIntercept, err := utils.GetFloat64(props, "Linear_Y_Intercept")
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
	n, err := utils.AsFloat64(v)
	if err != nil {
		return math.NaN()
	}
	return (n * scaler.LinearSlope) + scaler.LinearYIntercept
}
