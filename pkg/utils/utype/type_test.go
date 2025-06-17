package utype_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/wheatley/pkg/utils/utype"
)

func TestValueToPtr(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		expected := "hello"

		assert.Equal(t, &expected, utype.ValueToPtr("hello"))
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		expected := 10

		assert.Equal(t, &expected, utype.ValueToPtr(10))
	})

	t.Run("float", func(t *testing.T) {
		t.Parallel()

		expected := 3.14

		assert.Equal(t, &expected, utype.ValueToPtr(3.14))
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		expected := true

		assert.Equal(t, &expected, utype.ValueToPtr(true))
	})

	t.Run("nil_struct", func(t *testing.T) {
		t.Parallel()

		var expected *struct{ Name string }

		assert.Equal(t, &expected, utype.ValueToPtr(expected))
	})
}

func TestPtrToValue(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		value := "hello"

		assert.Equal(t, "hello", utype.PtrToValue(&value))
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		value := 10

		assert.Equal(t, 10, utype.PtrToValue(&value))
	})

	t.Run("float", func(t *testing.T) {
		t.Parallel()

		value := 3.14

		assert.InEpsilon(t, 3.14, utype.PtrToValue(&value), 0.001)
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		value := true

		assert.True(t, utype.PtrToValue(&value))
	})

	t.Run("nil_string", func(t *testing.T) {
		t.Parallel()

		var value *string

		assert.Empty(t, utype.PtrToValue(value))
	})

	t.Run("nil_int", func(t *testing.T) {
		t.Parallel()

		var value *int

		assert.Empty(t, utype.PtrToValue(value))
	})

	t.Run("nil_float", func(t *testing.T) {
		t.Parallel()

		var value *float32

		assert.Empty(t, utype.PtrToValue(value))
	})

	t.Run("nil_bool", func(t *testing.T) {
		t.Parallel()

		var value *bool

		assert.False(t, utype.PtrToValue(value))
	})

	type myTestingStruct struct {
		Name string
	}

	t.Run("struct_instance", func(t *testing.T) {
		t.Parallel()

		value := myTestingStruct{Name: "struct_instance_test"}

		assert.Equal(t, myTestingStruct{Name: "struct_instance_test"}, utype.PtrToValue(&value))
	})

	t.Run("struct_nil", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, myTestingStruct{}, utype.PtrToValue[myTestingStruct](nil))
	})
}
