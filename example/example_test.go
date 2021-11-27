//go:generate go test . -ldflags "-X 'github.com/bennydictor/canonical.Canonize=true'"
package example

import (
	"testing"

	"github.com/bennydictor/canonical"
)

func TestFoo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		result, err := Foo("bar")
		canonical.Assert(t, result, canonical.Error(err))
	})
	t.Run("err", func(t *testing.T) {
		result, err := Foo("foo")
		canonical.Assert(t, result, canonical.Error(err))
	})
}
