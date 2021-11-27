package canonical

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Canonize - set to any non-zero value to canonize tests.
var Canonize = ""

//goland:noinspection GoBoolExpressions
func canonizeEnabled() bool {
	return Canonize != ""
}

var canonicalValues = map[string]interface{}{}

const canonicalJson = "canonical.json"

var readOnce sync.Once

func read(t *testing.T) {
	if canonizeEnabled() {
		return
	}

	readOnce.Do(func() {
		file, err := os.Open(canonicalJson)
		if err != nil {
			t.Fatal("open canonical.json: ", err)
		}

		if err = json.NewDecoder(file).Decode(&canonicalValues); err != nil {
			t.Fatal("decode canonical.json: ", err)
		}
	})
}

var writeMu sync.Mutex

func write(t *testing.T) {
	buf, err := json.Marshal(canonicalValues)
	if err != nil {
		t.Fatal("marshal canonical.json: ", err)
	}

	var indented bytes.Buffer
	if err = json.Indent(&indented, buf, "", "\t"); err != nil {
		t.Fatal("indent canonical.json: ", err)
	}

	if err = ioutil.WriteFile(canonicalJson, indented.Bytes(), 0644); err != nil {
		t.Fatal("write canonical.json: ", err)
	}
}

func mustMarshal(t *testing.T, v interface{}) string {
	result, err := json.Marshal(v)
	if err != nil {
		t.Fatal("canonical: marshal failed: ", err)
	}
	return string(result)
}

// Assert asserts that the provided value matches the canonical value.
func Assert(t *testing.T, values ...interface{}) bool {
	if len(values) == 0 {
		t.Log("canonical.Assert of no values")
		return false
	}

	var value interface{}
	if len(values) == 1 {
		value = values[0]
	} else {
		value = values
	}

	if canonizeEnabled() {
		writeMu.Lock()
		defer writeMu.Unlock()

		canonicalValues[t.Name()] = value
		write(t)

		return true
	} else {
		read(t)

		canonicalValue, ok := canonicalValues[t.Name()]
		if !ok {
			t.Log("did not find canonical value")
			return false
		}

		return assert.JSONEq(t, mustMarshal(t, canonicalValue), mustMarshal(t, value))
	}
}

// Require requires that the provided value matches the canonical value.
func Require(t *testing.T, values ...interface{}) {
	if len(values) == 0 {
		t.Fatal("canonical.Require of no values")
	}

	var value interface{}
	if len(values) == 1 {
		value = values[0]
	} else {
		value = values
	}

	if canonizeEnabled() {
		writeMu.Lock()
		defer writeMu.Unlock()

		canonicalValues[t.Name()] = value
		write(t)
	} else {
		read(t)

		canonicalValue, ok := canonicalValues[t.Name()]
		if !ok {
			t.Fatal("did not find canonical value")
		}

		require.JSONEq(t, mustMarshal(t, canonicalValue), mustMarshal(t, value))
	}
}

// Error - canonical form of an error
func Error(err error) interface{} {
	if err == nil {
		return nil
	}

	return err.Error()
}
