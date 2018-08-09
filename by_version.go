package appcast

// ByVersion implements sort.Interface for []Release based on the Version field.
type ByVersion []Release

func (a ByVersion) Len() int {
	return len(a)
}

func (a ByVersion) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByVersion) Less(i, j int) bool {
	return a[i].Version().LessThan(a[j].Version())
}
