package client

import (
	"fmt"
	"strings"
)

// Service is an enumeration type for available services
type Service int

const (
	// ServiceCatalogs ...
	ServiceCatalogs Service = iota + 1

	// ServiceEnvironments ...
	ServiceEnvironments

	// ServiceClusters ...
	ServiceClusters

	// ServiceUsers ...
	ServiceUsers

	// ServiceSecurity ...
	ServiceSecurity

	// ServiceConfig ...
	ServiceConfig
)

// Name returns the service name
func (s Service) Name() string {
	switch s {
	case ServiceCatalogs:
		return "catalog"

	case ServiceEnvironments:
		return "environments"

	case ServiceClusters:
		return "cluster"

	case ServiceUsers:
		return "users"

	case ServiceSecurity:
		return "security"

	case ServiceConfig:
		return "config"

	default:
		return ""
	}
}

// ParseService converts a string to a service
func ParseService(s string) (Service, error) {
	switch strings.ToLower(s) {
	case "catalog":
		return ServiceCatalogs, nil

	case "environments":
		return ServiceEnvironments, nil

	case "cluster":
		return ServiceClusters, nil

	case "users":
		return ServiceUsers, nil

	case "security":
		return ServiceSecurity, nil

	case "config":
		return ServiceConfig, nil

	}

	return -1, fmt.Errorf("Invalid service %s", s)
}
