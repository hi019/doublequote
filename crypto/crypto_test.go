package crypto

import "testing"

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		NewService("test")
	})
}
