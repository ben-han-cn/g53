package g53

import ()

type NameRef struct {
	inner       *Name
	parentLevel int
}

func fromName(name *Name) NameRef {
	return NameRef{
		inner:       name,
		parentLevel: 0,
	}
}

func (r *NameRef) Parent() {
	r.parentLevel += 1
}

func (r *NameRef) Raw() []byte {
	return r.inner.raw[r.inner.offsets[r.parentLevel]:]
}

func (r *NameRef) IsRoot() bool {
	return r.parentLevel+1 == int(r.inner.labelCount)
}

func (r *NameRef) Hash(caseSensitive bool) uint32 {
	raw := r.Raw()
	hashLen := len(raw)
	hash := uint32(0)
	if caseSensitive {
		for i := 0; i < hashLen; i++ {
			hash ^= uint32(raw[i]) + 0x9e3779b9 + (hash << 6) + (hash >> 2)
		}
	} else {
		for i := 0; i < hashLen; i++ {
			hash ^= uint32(maptolower[raw[i]]) + 0x9e3779b9 + (hash << 6) + (hash >> 2)
		}
	}
	return hash
}
