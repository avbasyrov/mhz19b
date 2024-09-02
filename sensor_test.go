package mhz19b

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	t.Parallel()

	sum := checksum([]byte{0xFF, 0x86, 0x02, 0x48, 0x3E, 0x00, 0x00, 0x00, 0xF2})
	assert.Equal(t, uint8(0xF2), sum)
}

func TestCheckMessage(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		err := checkMessage([]byte{0xFF, 0x86, 0x02, 0x48, 0x3E, 0x00, 0x00, 0x00, 0xF2})
		assert.NoError(t, err)
	})

	t.Run("too short", func(t *testing.T) {
		t.Parallel()

		err := checkMessage([]byte{0xFF, 0x86, 0x02, 0x48, 0x3E, 0x00, 0x00, 0x00})
		assert.ErrorContains(t, err, "unexpected reply length")
	})

	t.Run("wrong 1st byte", func(t *testing.T) {
		t.Parallel()

		err := checkMessage([]byte{0xAA, 0x86, 0x02, 0x48, 0x3E, 0x00, 0x00, 0x00, 0xF2})
		assert.ErrorContains(t, err, "unexpected 1st byte in reply")
	})

	t.Run("wrong checksum", func(t *testing.T) {
		t.Parallel()

		err := checkMessage([]byte{0xFF, 0x86, 0x02, 0x48, 0x3E, 0x00, 0x00, 0x00, 0xAA})
		assert.ErrorContains(t, err, "bad checksum")
	})
}
