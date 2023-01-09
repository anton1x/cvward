package internal

type FrameInfo struct {
	info      map[int]int
	threshold int
}

func NewFrameInfo() *FrameInfo {
	return &FrameInfo{
		info: make(map[int]int),
	}
}

func (fi *FrameInfo) Put(frameID, numObj int) {
	fi.info[frameID] = numObj
}
