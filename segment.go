package tdms

import "io"

type Segment struct {
	Type     SegmentType
	LeadIn   *LeadIn
	MetaData *MetaData
	Offset   int64
}

func ReadSegment(r io.Reader, previousSegment *Segment) (*Segment, error) {
	var segment Segment
	var err error

	segment.Type, segment.LeadIn, err = readLeadIn(r)
	if err != nil {
		return nil, err
	}
	if segment.LeadIn.ToC.MetaData() {
		segment.MetaData, err = ReadMetaData(r, segment.LeadIn.ToC, previousSegment)
		if err != nil {
			return nil, err
		}
	}

	return &segment, nil
}
