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

import "fmt"

var valueFromFuncs = map[string]func(interface{}) (interface{}, error){
	"count": vfCount,
	"sum":   vfSum,
	"min":   vfMin,
	"max":   vfMax,
}

func vfCount(i interface{}) (interface{}, error) {
	switch val := i.(type) {
	case []interface{}:
		return float64(len(val)), nil
	case map[string]interface{}:
		return float64(len(val)), nil
	}

	return float64(0), fmt.Errorf("cannot count non-collection type")
}

func vfSum(i interface{}) (interface{}, error) {
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
	default:
		return float64(0), fmt.Errorf("cannot sum non-collection type")
	}
	return sum, nil
}

func vfMin(i interface{}) (interface{}, error) {
	var min *float64
	switch val := i.(type) {
	case []interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if min == nil || num < *min {
					min = &num
				}
			}
		}
	case map[string]interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if min == nil || num < *min {
					min = &num
				}
			}
		}
	}
	if min == nil {
		return nil, nil
	}
	return *min, nil
}

func vfMax(i interface{}) (interface{}, error) {
	var max *float64
	switch val := i.(type) {
	case []interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if max == nil || num > *max {
					max = &num
				}
			}
		}
	case map[string]interface{}:
		for _, item := range val {
			if num, err := toFloat64(item, false); err == nil {
				if max == nil || num > *max {
					max = &num
				}
			}
		}
	}
	if max == nil {
		return nil, nil
	}
	return *max, nil
}
