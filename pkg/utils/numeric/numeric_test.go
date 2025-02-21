package numeric_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/utils/numeric"
)

func TestFloat64Value(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name    string
		value   any
		want    float64
		wantErr bool
	}{
		{
			name:  "int",
			value: 1,
			want:  1,
		},
		{
			name:  "int8",
			value: int8(1),
			want:  1,
		},
		{
			name:  "int16",
			value: int16(1),
			want:  1,
		},
		{
			name:  "int32",
			value: int32(1),
			want:  1,
		},
		{
			name:  "int64",
			value: int64(1),
			want:  1,
		},
		{
			name:  "uint",
			value: uint(1),
			want:  1,
		},
		{
			name:  "uint8",
			value: uint8(1),
			want:  1,
		},
		{
			name:  "uint16",
			value: uint16(1),
			want:  1,
		},
		{
			name:  "uint32",
			value: uint32(1),
			want:  1,
		},
		{
			name:  "uint64",
			value: uint64(1),
			want:  1,
		},
		{
			name:  "float32",
			value: float32(1.0),
			want:  1.0,
		},
		{
			name:  "float64",
			value: 1.0,
			want:  1.0,
		},
		{
			name:  "string",
			value: "1.0",
			want:  1.0,
		},
		{
			name:    "unsupported",
			value:   struct{}{},
			wantErr: true,
		},
		{
			name:    "unsupported_bool",
			value:   true,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := numeric.Float64Value(test.value)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.InEpsilon(t, test.want, got, 0.0001)
		})
	}
}
