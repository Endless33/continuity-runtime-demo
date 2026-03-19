package runtime

type Dedup struct {
	seen map[int]bool
}

func NewDedup() *Dedup {
	return &Dedup{
		seen: make(map[int]bool),
	}
}

func (d *Dedup) Seen(id int) bool {
	if d.seen[id] {
		return true
	}

	d.seen[id] = true
	return false
}