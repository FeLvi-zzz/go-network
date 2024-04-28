package types

type Flags uint8

func (f Flags) ToString() string {
	switch {
	case f.IsDF():
		return "Don't Fragment"
	case f.IsMF():
		return "More Fragments"
	}

	return "None"
}

func (f Flags) IsDF() bool {
	return f&0b010 != 0
}

func (f Flags) IsMF() bool {
	return f&0b001 != 0
}
