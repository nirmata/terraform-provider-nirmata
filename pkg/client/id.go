package client

// ID identifies a model object
type ID interface {
	Service() Service
	ModelIndex() string
	UUID() string
	Map() map[string]interface{}
}

type id struct {
	service    Service
	modelIndex string
	uuid       string
}

// NewID creates a new ID
func NewID(service Service, modelIndex string, uuid string) ID {
	return &id{service, modelIndex, uuid}
}

func (i *id) Service() Service {
	return i.service
}

func (i *id) ModelIndex() string {
	return i.modelIndex
}

func (i *id) UUID() string {
	return i.uuid
}

func (i *id) Map() map[string]interface{} {
	return map[string]interface{}{
		"service":    i.Service().Name(),
		"modelIndex": i.ModelIndex(),
		"id":         i.UUID(),
	}
}
