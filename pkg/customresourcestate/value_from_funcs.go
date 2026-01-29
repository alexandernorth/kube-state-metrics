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
	"strconv"
)

const LegacyFunctionName = "legacy"

type ValueFromFunc struct {
	FuncName    string
	Func        func(v interface{}, args ...string) (interface{}, error)
	IsAggregate bool
}

var valueFromFuncs = map[string]ValueFromFunc{
	LegacyFunctionName: {
		FuncName:    LegacyFunctionName,
		Func:        vfLegacy,
		IsAggregate: false,
	},
	"count": {
		FuncName:    "count",
		Func:        vfCount,
		IsAggregate: true,
	},
	"literal": {
		FuncName:    "literal",
		Func:        vfLiteral,
		IsAggregate: false,
	},
}

func vfLegacy(v interface{}, args ...string) (interface{}, error) {
	if len(args) == 0 {
		return v, nil
	}
	// TODO: Should we return an error if the path does not exist? Probably not.
	switch vv := v.(type) {
	case []interface{}:
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		if i < 0 {
			// negative index
			i += len(vv)
		}
		if i < 0 || i >= len(vv) {
			return nil, fmt.Errorf("list index out of range: %s", args[0])
		}
		return vfLegacy(vv[i], args[1:]...)
	case map[string]interface{}:
		found, ok := vv[args[0]]
		if !ok {
			return nil, nil
		}
		return vfLegacy(found, args[1:]...)
	default:
		return nil, nil
	}
}

func vfCount(v interface{}, args ...string) (interface{}, error) {
	switch vv := v.(type) {
	case []interface{}:
		return len(vv), nil
	case map[string]interface{}:
		return len(vv), nil
	default:
		return 0, nil
	}
}

func vfLiteral(v interface{}, args ...string) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("literal function requires exactly one argument")
	}

	f, err := toFloat64(args[0], false)
	if err != nil {
		return nil, fmt.Errorf("could not parse literal value %s: %w", args[0], err)
	}
	return f, nil
}
