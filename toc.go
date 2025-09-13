package tdms

type TableOfContents uint32

func (toc TableOfContents) MetaData() bool {
	return toc&(1<<1) != 0
}

func (toc TableOfContents) RawData() bool {
	return toc&(1<<3) != 0
}

func (toc TableOfContents) DAQmxRawData() bool {
	return toc&(1<<7) != 0
}

func (toc TableOfContents) InterleavedData() bool {
	return toc&(1<<5) != 0
}

func (toc TableOfContents) BigEndian() bool {
	return toc&(1<<6) != 0
}

func (toc TableOfContents) NewObjList() bool {
	return toc&(1<<2) != 0
}

func (toc TableOfContents) ValueReader() *ValueReader {
	if toc.BigEndian() {
		return BigEndianValueReader
	} else {
		return LittleEndianValueReader
	}
}
