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
	"fmt"
)

// FuncType indicates how a function processes values.
type FuncType int

const (
	// FuncIterative means the function is applied to each value individually.
	// Results in N metrics for N input values.
	FuncIterative FuncType = iota
	// FuncAggregate means the function receives all values and returns a single result.
	// Results in 1 metric regardless of input count.
	FuncAggregate
)

// ValueFunc is a function that computes or transforms a value.
// It receives the input value and optional string arguments.
type ValueFunc func(value interface{}, args []string) (interface{}, error)

// ValueFuncDef defines a value function with its type metadata.
type ValueFuncDef struct {
	Func ValueFunc
	Type FuncType
}

// valueFuncs maps function names to their definitions.
var valueFuncs = map[string]ValueFuncDef{
	// Iterative functions - applied to each value, returns N metrics
	"filter":   {Func: vfFilter, Type: FuncIterative},
	"multiply": {Func: vfMultiply, Type: FuncIterative},
	"add":      {Func: vfAdd, Type: FuncIterative},
	// Aggregate functions - receives all values, returns 1 metric
	"count":  {Func: vfCount, Type: FuncAggregate},
	"sum":    {Func: vfSum, Type: FuncAggregate},
	"min":    {Func: vfMin, Type: FuncAggregate},
	"max":    {Func: vfMax, Type: FuncAggregate},
	"scalar": {Func: vfScalar, Type: FuncAggregate},
}

// vfFilter passes the value through unchanged (identity function).
func vfFilter(i interface{}, _ []string) (interface{}, error) {
	return i, nil
}

// vfMultiply multiplies the value by args[0].
func vfMultiply(i interface{}, args []string) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("multiply function requires a multiplier argument")
	}
	multiplier, err := toFloat64(args[0], false)
	if err != nil {
		return nil, fmt.Errorf("invalid multiplier: %w", err)
	}
	val, err := toFloat64(i, false)
	if err != nil {
		return nil, err
	}
	return val * multiplier, nil
}

// vfAdd adds args[0] to the value.
func vfAdd(i interface{}, args []string) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("add function requires an addend argument")
	}
	addend, err := toFloat64(args[0], false)
	if err != nil {
		return nil, fmt.Errorf("invalid addend: %w", err)
	}
	val, err := toFloat64(i, false)
	if err != nil {
		return nil, err
	}
	return val + addend, nil
}

// vfScalar returns a static value specified in args[0], ignoring the input.
// Usage: func: {name: scalar, args: ["42"]}
func vfScalar(_ interface{}, args []string) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("scalar function requires at least one argument")
	}
	return toFloat64(args[0], false)
}

func vfCount(i interface{}, _ []string) (interface{}, error) {
	switch val := i.(type) {
	case []interface{}:
		return float64(len(val)), nil
	case map[string]interface{}:
		return float64(len(val)), nil
	case nil:
		return nil, nil
	}

	return nil, fmt.Errorf("cannot count non-collection type")
}

func vfSum(i interface{}, _ []string) (interface{}, error) {
	sum := float64(0)
	switch val := i.(type) {
	case []interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				sum += num
			}
		}
	case map[string]interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				sum += num
			}
		}
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot sum non-collection type")
	}

	return sum, nil
}

func vfMin(i interface{}, _ []string) (interface{}, error) {
	var minVal *float64
	switch val := i.(type) {
	case []interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if minVal == nil || num < *minVal {
					minVal = &num
				}
			}
		}
	case map[string]interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if minVal == nil || num < *minVal {
					minVal = &num
				}
			}
		}
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot find min of non-collection type")
	}

	if minVal == nil {
		return nil, nil
	}

	return *minVal, nil
}

func vfMax(i interface{}, _ []string) (interface{}, error) {
	var maxVal *float64
	switch val := i.(type) {
	case []interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if maxVal == nil || num > *maxVal {
					maxVal = &num
				}
			}
		}
	case map[string]interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if maxVal == nil || num > *maxVal {
					maxVal = &num
				}
			}
		}
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot find max of non-collection type")
	}

	if maxVal == nil {
		return nil, nil
	}

	return *maxVal, nil
}
