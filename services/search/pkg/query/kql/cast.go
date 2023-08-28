package kql

import (
	"fmt"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func toIfaceSlice(in interface{}) []interface{} {
	if in == nil {
		return nil
	}
	return in.([]interface{})
}

func toNode(in interface{}) (ast.Node, error) {
	out, ok := in.(ast.Node)
	if !ok {
		return nil, fmt.Errorf("can't convert '%T' to ast.Node", in)
	}

	return out, nil
}

func toNodes(in interface{}) ([]ast.Node, error) {

	switch v := in.(type) {
	case []interface{}:
		var nodes []ast.Node

		for _, el := range toIfaceSlice(v) {
			node, err := toNode(el)
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, node)
		}

		return nodes, nil
	case []ast.Node:
		return v, nil
	default:
		return nil, fmt.Errorf("can't convert '%T' to []ast.Node", in)
	}
}

func toString(in interface{}) (string, error) {
	switch v := in.(type) {
	case []byte:
		return string(v), nil
	case []interface{}:
		var str string

		for _, i := range v {
			j := i.([]uint8)
			str += string(j[0])
		}

		return str, nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("can't convert '%T' to string", v)
	}
}
