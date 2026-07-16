package env

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func uniqueKey(key string) string {
	prefix := strconv.Itoa(time.Now().Nanosecond())
	return fmt.Sprintf("%v_%v", prefix, key)
}

func Test_getStringValue(t *testing.T) {
	t.Run("env var exists", func(t *testing.T) {
		key := uniqueKey("existing_env")
		val := strconv.Itoa(rand.IntN(1000))

		os.Setenv(key, val)
		defer os.Unsetenv(key)

		assert.EqualValues(t, val, getStringValue(key, "not val"))
	})

	t.Run("env var does not exist", func(t *testing.T) {
		key := uniqueKey("non_existent_env")
		val := strconv.Itoa(rand.IntN(1000))

		assert.EqualValues(t, val, getStringValue(key, val))
	})
}

func Test_getIntValue(t *testing.T) {
	t.Run("existing int env var", func(t *testing.T) {
		key := uniqueKey("existing_int_env")
		val := rand.IntN(1000)

		os.Setenv(key, strconv.Itoa(val))
		defer os.Unsetenv(key)

		assert.EqualValues(t, val, getIntValue(key, 420))
	})

	t.Run("int env var does not exist", func(t *testing.T) {
		key := uniqueKey("existing_int_env")
		val := rand.IntN(1000)

		assert.EqualValues(t, val, getIntValue(key, val))
	})

	t.Run("non parseable int env var", func(t *testing.T) {
		key := uniqueKey("existing_int_env")
		val := rand.IntN(1000)

		os.Setenv(key, "not a parsebale int")
		defer os.Unsetenv(key)

		assert.EqualValues(t, val, getIntValue(key, val))
	})
}

func Test_getBoolValue(t *testing.T) {
	t.Run("existing bool env var", func(t *testing.T) {
		key := uniqueKey("bool_existing_env")
		val := true

		os.Setenv(key, strconv.FormatBool(val))
		defer os.Unsetenv(key)

		assert.EqualValues(t, val, getBoolValue(key, false))
	})

	t.Run("non-existent bool env var", func(t *testing.T) {
		key := uniqueKey("bool_existing_env")
		val := true

		assert.EqualValues(t, val, getBoolValue(key, val))
	})

	t.Run("non-parsebale bool env var", func(t *testing.T) {
		key := uniqueKey("bool_existing_env")
		val := true

		os.Setenv(key, "unparseble bool value")
		defer os.Unsetenv(key)

		assert.EqualValues(t, val, getBoolValue(key, val))
	})
}

func Test_getDurationValue(t *testing.T) {
	t.Run("existing duration env var", func(t *testing.T) {
		key := uniqueKey("duration_existing_env")
		randInt := rand.IntN(1000)
		val := time.Duration(randInt) * time.Second
		def := time.Duration(randInt+1) * time.Second

		os.Setenv(key, val.String())
		defer os.Unsetenv(key)

		assert.EqualValues(t, val, getDurationValue(key, def))
	})

	t.Run("non-existent duration env var", func(t *testing.T) {
		key := uniqueKey("duration_existing_env")
		randInt := rand.IntN(1000)

		def := time.Duration(randInt) * time.Second

		assert.EqualValues(t, def, getDurationValue(key, def))
	})

	t.Run("non-parseble duration env var", func(t *testing.T) {
		key := uniqueKey("duration_existing_env")
		randInt := rand.IntN(1000)

		def := time.Duration(randInt) * time.Second

		os.Setenv(key, "un parseble duration")

		assert.EqualValues(t, def, getDurationValue(key, def))
	})
}
