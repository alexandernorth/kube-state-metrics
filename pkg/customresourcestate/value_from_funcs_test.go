/*
Copyright 2021 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package customresourcestate

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_vfCount(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{name: "empty array", input: []interface{}{}, want: float64(0)},
		{name: "array with elements", input: []interface{}{1, 2, 3}, want: float64(3)},
		{name: "empty map", input: map[string]interface{}{}, want: float64(0)},
		{name: "map with elements", input: map[string]interface{}{"a": 1, "b": 2}, want: float64(2)},
		{name: "non-collection", input: "hello", wantErr: true},
		{name: "nil", input: nil, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vfCount(tt.input, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_vfSum(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{name: "empty array", input: []interface{}{}, want: float64(0)},
		{name: "empty map", input: map[string]interface{}{}, want: float64(0)},
		{name: "array with numbers", input: []interface{}{float64(1), float64(2), float64(3)}, want: float64(6)},
		{name: "map with numbers", input: map[string]interface{}{"a": float64(1), "b": float64(2)}, want: float64(3)},
		{name: "mixed types ignores non-numeric", input: []interface{}{float64(1), "skip", nil, float64(2)}, want: float64(3)},
		{name: "non-collection", input: "hello", wantErr: true},
		{name: "nil", input: nil, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vfSum(tt.input, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_vfMin(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{name: "empty array", input: []interface{}{}, want: nil},
		{name: "empty map", input: map[string]interface{}{}, want: nil},
		{name: "array finds minimum", input: []interface{}{float64(3), float64(1), float64(2)}, want: float64(1)},
		{name: "map finds minimum", input: map[string]interface{}{"a": float64(5), "b": float64(2)}, want: float64(2)},
		{name: "zero as minimum", input: []interface{}{float64(0), float64(1)}, want: float64(0)},
		{name: "negative infinity minimum", input: []interface{}{math.Inf(-1), float64(1)}, want: math.Inf(-1)},
		{name: "ignores non-numeric", input: []interface{}{"skip", float64(5), nil}, want: float64(5)},
		{name: "all non-numeric", input: []interface{}{"a", nil}, want: nil},
		{name: "non-collection", input: "hello", wantErr: true},
		{name: "nil", input: nil, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vfMin(tt.input, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_vfMax(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{name: "empty array", input: []interface{}{}, want: nil},
		{name: "empty map", input: map[string]interface{}{}, want: nil},
		{name: "array finds maximum", input: []interface{}{float64(3), float64(5), float64(2)}, want: float64(5)},
		{name: "map finds maximum", input: map[string]interface{}{"a": float64(1), "b": float64(8)}, want: float64(8)},
		{name: "zero as maximum", input: []interface{}{float64(-1), float64(0)}, want: float64(0)},
		{name: "infinity", input: []interface{}{float64(1), math.Inf(1)}, want: math.Inf(1)},
		{name: "ignores non-numeric", input: []interface{}{"skip", float64(5), nil}, want: float64(5)},
		{name: "all non-numeric", input: []interface{}{"a", nil}, want: nil},
		{name: "non-collection", input: "hello", wantErr: true},
		{name: "nil", input: nil, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vfMax(tt.input, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_valueFromFuncsRegistry(t *testing.T) {
	expected := []string{"filter", "count", "sum", "min", "max", "scalar"}
	assert.Len(t, valueFuncs, len(expected))
	for _, name := range expected {
		assert.Contains(t, valueFuncs, name)
	}
}
