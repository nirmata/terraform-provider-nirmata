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

func validateGKEMachineType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must match [\\w+=,.@-]", k))
	}
	return
}

func validateGKEDiskSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 9 {
		errors = append(errors, fmt.Errorf("%q The disk size must be grater than 9", k))
	}
	return
}

func validateGKELocationType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Zonal" && value != "Regional" {
		errors = append(errors, fmt.Errorf("location_type (%v) must Zonal or Regional", v))
	}

	return
}
