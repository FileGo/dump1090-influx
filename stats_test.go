package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnixTimeUnmarshalJSON(t *testing.T) {
	assert := assert.New(t)

	t.Run("pass", func(t *testing.T) {
		jsonText := []byte(`946684800.5`)

		var out unixTime
		err := json.Unmarshal(jsonText, &out)
		if assert.Nil(err) {
			assert.Equal(2000, out.Year())
			assert.Equal(1, out.Day())
			assert.Equal(0, out.Hour())
			assert.Equal(0, out.Minute())
			assert.Equal(0, out.Second())
			assert.Equal(int(5*1e8), out.Nanosecond())
		}
	})

	t.Run("fail", func(t *testing.T) {
		jsonText := []byte(`"notanumb3r"`)
		var out unixTime
		err := json.Unmarshal(jsonText, &out)
		assert.NotNil(err)
	})
}
