package nirmata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jmespath/go-jmespath"
	"log"
)

func updateData(fieldName string, d *schema.ResourceData, s *schema.Schema, jmesPath string, data interface{}) error {
	v, err := jmespath.Search(jmesPath, data)
	if err != nil {
		return err
	}

	cv := convertValue(v, s)
	err = d.Set(fieldName, cv)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] updated %s (%s) to %v", fieldName, jmesPath, cv)
	return nil
}

func convertValue(v interface{}, s *schema.Schema) interface{} {
	if v == nil {
		return nil
	}
	
	switch s.Type {
	case schema.TypeBool:
		return v.(bool)

	case schema.TypeFloat:
		return v.(float64)

	case schema.TypeInt:
		return v.(int)

	case schema.TypeString:
		return v.(string)

	case schema.TypeList:
		return makeList(v, s)

	case schema.TypeMap:
		return makeMap(v, s)

	case schema.TypeSet:
		return makeSet(v.([]interface{}))

	default:
		return v
	}
}

func makeList(v interface{}, s *schema.Schema) interface{} {
	if v == nil || s.Elem == nil {
		return nil
	}

	switch s.Elem.(type) {
	case *schema.Schema:
		es := s.Elem.(*schema.Schema)
		switch es.Type {
		case schema.TypeString:
			return makeListOfString(v.([]interface{}))

		case schema.TypeBool:
			return makeListOfBool(v.([]interface{}))

		case schema.TypeInt:
			return makeListOfInt(v.([]interface{}))

		case schema.TypeFloat:
			return makeListOfFloat64(v.([]interface{}))

		default:
			return v
		}

	default:
		return v
	}
}

func makeListOfString(v []interface{}) []string {
	if v == nil {
		return nil
	}

	results := make([]string, len(v))
	for i, e := range v {
		if e == nil {
			results[i] = ""
			continue
		}

		results[i] = e.(string)
	}

	return results
}

func makeListOfInt(v []interface{}) []int {
	if v == nil {
		return nil
	}

	results := make([]int, len(v))
	for i, e := range v {
		if e == nil {
			results[i] = 0
			continue
		}

		results[i] = e.(int)
	}

	return results
}

func makeListOfFloat64(v []interface{}) []float64 {
	if v == nil {
		return nil
	}

	results := make([]float64, len(v))
	for i, e := range v {
		if e == nil {
			results[i] = 0
			continue
		}

		results[i] = e.(float64)
	}

	return results
}

func makeListOfBool(v []interface{}) []bool {
	if v == nil {
		return nil
	}

	results := make([]bool, len(v))
	for i, e := range v {
		if e == nil {
			results[i] = false
			continue
		}

		results[i] = e.(bool)
	}

	return results
}

func makeMap(v interface{}, s *schema.Schema) interface{} {
	if v == nil || s.Elem == nil {
		return nil
	}

	switch s.Elem.(type) {
	case *schema.Schema:
		es := s.Elem.(*schema.Schema)
		switch es.Type {
		case schema.TypeString:
			return makeMapOfString(v.(map[string]interface{}))

		default:
			return v
		}

	default:
		return v
	}
}

func makeMapOfString(m map[string]interface{}) map[string]string {
	if m == nil {
		return nil
	}

	result := make(map[string]string)
	for k, v := range m {
		if v == nil {
			result[k] = ""
			continue
		}

		result[k] = v.(string)
	}

	return result
}

func makeSet(v []interface{}) *schema.Set {
	if v == nil {
		return nil
	}

	result := &schema.Set{}
	for _, e := range v {
		result.Add(e)
	}

	return result
}
