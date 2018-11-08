package appcaster

// Providerer is the Provider interface.
//
// This interface should be embedded by your own Providerer interface.
type Providerer interface {
	String() string
}

// Provider holds different supported providers.
//
// This should be aliased by your own Provider.
type Provider int

// String returns the string representation of the Provider.
func (p Provider) String() string {
	panic("implement me")
}
