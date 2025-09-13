package tdms

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type SegmentType string

const (
	SegmentTypeTDSm SegmentType = "TDSm"
	SegmentTypeTDSh SegmentType = "TDSh"
)

type LeadIn struct {
	ToC               TableOfContents
	VersionNumber     uint32
	NextSegmentOffset uint64
	RawDataOffset     uint64
}

func readLeadIn(r io.Reader) (SegmentType, *LeadIn, error) {
	segmentTag := make([]byte, 4)
	_, err := io.ReadFull(r, segmentTag)
	if err != nil {
		return "", nil, err
	}
	var segmentType SegmentType
	if bytes.Compare(segmentTag, []byte(SegmentTypeTDSm)) == 0 {
		segmentType = SegmentTypeTDSm
	} else if bytes.Compare(segmentTag, []byte(SegmentTypeTDSh)) == 0 {
		segmentType = SegmentTypeTDSh
	} else {
		return "", nil, fmt.Errorf("unknown segment tag: %v", segmentTag)
	}
	var leadIn LeadIn
	err = binary.Read(r, binary.LittleEndian, &leadIn.ToC)
	if err != nil {
		return "", nil, err
	}
	valueReader := leadIn.ToC.ValueReader()
	leadIn.VersionNumber, err = valueReader.ReadU32(r)
	if err != nil {
		return "", nil, err
	}
	leadIn.NextSegmentOffset, err = valueReader.ReadU64(r)
	if err != nil {
		return "", nil, err
	}
	leadIn.RawDataOffset, err = valueReader.ReadU64(r)
	if err != nil {
		return "", nil, err
	}
	return segmentType, &leadIn, nil
}
