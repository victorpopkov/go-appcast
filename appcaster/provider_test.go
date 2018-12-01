package appcaster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_String(t *testing.T) {
	p := Provider(0)
	assert.Panics(t, func() {
		_ = p.String()
	})
}
