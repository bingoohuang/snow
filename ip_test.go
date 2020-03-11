package snow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferHostIPv4(t *testing.T) {
	that := assert.New(t)

	that.NotEmpty(InferHostIPv4(""))
}
