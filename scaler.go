package tdms

type Scaler interface {
	Scale(v any) float64
}
