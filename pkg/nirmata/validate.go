package nirmata

import (
	"fmt"
	"regexp"
)

func validateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters", k))
	}

	if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must match [\\w+=,.@-]", k))
	}

	return
}

func validateAksFields(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must match [\\w+=,.@-]", k))
	}
	return
}

func validateNodeCount(v interface{}, k string) (ws []string, errors []error) {
	count := v.(int)
	if count < 0 || count > 1000 {
		errors = append(errors, fmt.Errorf("node count (%d) must be between 0 and 1000", count))
	}

	return
}

func validateGKEMachineType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must match [\\w+=,.@-]", k))
	}
	return
}

func validateGKEDiskSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 9 {
		errors = append(errors, fmt.Errorf("disk_size (%s) must be greater than 9", k))
	}
	return
}

func validateGKELocationType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Zonal" && value != "Regional" {
		errors = append(errors, fmt.Errorf("location_type (%s) must Zonal or Regional", v))
	}

	return
}

func validateEKSDiskSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 9 {
		errors = append(errors, fmt.Errorf(
			"%q The disk size must be grater than 9", k))
	}

	return
}

func validateDeleteAction(v interface{}, k string) (ws []string, errors []error) {
	val := v.(string)
	if val != "delete" && val != "remove" {
		errors = append(errors, fmt.Errorf(
			"%s is not a valid delete_action; expected 'delete' or 'remove'", val))
	}

	return
}

func validatePermission(v interface{}, k string) (ws []string, errors []error) {
	val := v.(string)
	if val != "admin" && val != "edit" && val != "view" {
		errors = append(errors, fmt.Errorf(
			"%s is not a valid permission; expected 'admin' or 'edit' or 'view'", val))
	}

	return
}

func validateEntity(v interface{}, k string) (ws []string, errors []error) {
	val := v.(string)
	if val != "user" && val != "team" {
		errors = append(errors, fmt.Errorf(
			"%s is not a valid entity; expected 'user' or 'team'", val))
	}

	return
}
