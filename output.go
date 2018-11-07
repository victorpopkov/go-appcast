package appcast

// Outputer is the interface that wraps the Output methods.
//
// This interface should be embedded by more use-case specific Outputer
// interfaces.
type Outputer interface {
	Save() error
	GenerateChecksum(algorithm ChecksumAlgorithm) *Checksum
	Content() []byte
	SetContent(content []byte)
	Checksum() *Checksum
	Provider() Providerer
	SetProvider(provider Providerer)
	Appcast() Appcaster
	SetAppcast(appcast Appcaster)
}

// Output represents an appcast output data.
//
// This struct should be embedded by more use-case specific Output structs.
type Output struct {
	// content specifies an appcast data which will be outputted.
	content []byte

	// checksum specifies the Checksum pointer that hold a checksum data about
	// the Output.content.
	//
	// This value is set by calling Output.GenerateChecksum.
	checksum *Checksum

	// provider specifies the one of the supported appcast providers.
	provider Providerer

	// appcast specifies the provider-specific appcast guessed for the current
	// Output.
	appcast Appcaster
}

// Save should implement a way of saving an appcast content from the
// Output.content depending on the chosen Outputer.
func (o *Output) Save() error {
	panic("implement me")
}

// GenerateChecksum creates a new Checksum instance pointer based on the
// provided algorithm and sets it as an Output.checksum.
func (o *Output) GenerateChecksum(algorithm ChecksumAlgorithm) *Checksum {
	c := NewChecksum(algorithm, o.content)
	o.checksum = c
	return c
}

// Content is an Output.content getter.
func (o *Output) Content() []byte {
	return o.content
}

// SetContent is an Output.content setter.
func (o *Output) SetContent(content []byte) {
	o.content = content
}

// Checksum is an Output.checksum getter.
func (o *Output) Checksum() *Checksum {
	return o.checksum
}

// Provider is an Output.provider getter.
func (o *Output) Provider() Providerer {
	return o.provider
}

// SetProvider is an Output.provider setter.
func (o *Output) SetProvider(provider Providerer) {
	o.provider = provider
}

// Appcast is a Output.appcast getter.
func (o *Output) Appcast() Appcaster {
	return o.appcast
}

// SetAppcast is a Output.appcast setter.
func (o *Output) SetAppcast(appcast Appcaster) {
	o.appcast = appcast
}
