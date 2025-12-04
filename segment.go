package tdms

type Segment struct {
	Type     SegmentType
	LeadIn   *LeadIn
	MetaData *MetaData
	Offset   int64
}
