package tdms

type Object struct {
	Path         string
	RawDataIndex RawDataIndex
	Properties   map[string]any
}
