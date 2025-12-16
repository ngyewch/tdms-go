package tdms

import (
	"fmt"

	"github.com/ngyewch/tdms-go/utils"
)

type LinearScaler struct {
	scaleId           uint32
	linearInputSource uint
	linearSlope       float64
	linearYIntercept  float64
}

func NewLinearScaler(scaleId uint32, props map[string]any) (*LinearScaler, error) {
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
		scaleId:           scaleId,
		linearInputSource: linearInputSource,
		linearSlope:       linearSlope,
		linearYIntercept:  linearYIntercept,
	}, nil
}

func (scaler *LinearScaler) ScaleId() uint32 {
	return scaler.scaleId
}

func (scaler *LinearScaler) Scale(v any) (float64, error) {
	n, err := utils.AsFloat64(v)
	if err != nil {
		return 0, err
	}
	return (n * scaler.linearSlope) + scaler.linearYIntercept, nil
}
