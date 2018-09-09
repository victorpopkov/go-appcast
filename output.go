package appcast

// Outputer is the interface that wraps the Output methods.
//
// This interface should be embedded by more use-case specific Outputer
// interfaces.
type Outputer interface {
	Save() error
	GenerateChecksum(algorithm ChecksumAlgorithm)
	Content() []byte
	SetContent(content []byte)
	Checksum() *Checksum
	Provider() Provider
	SetProvider(provider Provider)
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
	provider Provider
}

// Save should implement a way of saving an appcast content from the
// Output.content depending on the chosen Outputer.
func (o *Output) Save() error {
	panic("implement me")
}

// GenerateChecksum creates a new Checksum instance based on the provided
// algorithm.
func (o *Output) GenerateChecksum(algorithm ChecksumAlgorithm) {
	o.checksum = NewChecksum(algorithm, o.content)
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
func (o *Output) Provider() Provider {
	return o.provider
}

// SetProvider is an Output.provider setter.
func (o *Output) SetProvider(provider Provider) {
	o.provider = provider
}
