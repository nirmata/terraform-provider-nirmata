package client

import "fmt"

// Object has an unique ID, attributes, and relations
type Object interface {
	ID() ID
	Data() map[string]interface{}
	GetString(name string) string
	GetRelation(name string) (ID, error)

	// TODO - to be implemented
	// GetRelations(name string) []ID
}

type object struct {
	id   ID
	data map[string]interface{}
}

// NewObject creates a new Object
func NewObject(data map[string]interface{}) (Object, error) {
	id, err := ParseIDFromMap(data)
	if err != nil {
		return nil, err
	}

	return &object{id, data}, nil
}

func (o *object) ID() ID {
	return o.id
}

func (o *object) Data() map[string]interface{} {
	return o.data
}

func (o *object) GetString(key string) string {
	val := o.data[key]
	if val == nil {
		return ""
	}

	return val.(string)
}

func (o *object) GetRelation(name string) (ID, error) {
	val := o.data[name]
	if val == nil {
		return nil, fmt.Errorf("Relation %s not found", name)
	}

	idData := val.(map[string]interface{})
	return ParseIDFromMap(idData)
}
