package tdms

type Scaler interface {
	ScaleId() uint32
	Scale(v any) (float64, error)
}
