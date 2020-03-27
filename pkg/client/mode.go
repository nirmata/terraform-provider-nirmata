package client

// OutputMode sets the output mode for a REST API call
type OutputMode int

const (
	// OutputModeNone ...
	OutputModeNone OutputMode = iota + 1

	// OutputModeDefault ...
	OutputModeDefault

	// OutputModeExport ...
	OutputModeExport

	// OutputModeExportDetails ...
	OutputModeExportDetails
)

// Name returns the service name
func (m OutputMode) String() string {
	switch m {
	case OutputModeNone:
		return ""

	case OutputModeDefault:
		return "default"

	case OutputModeExport:
		return "export"

	case OutputModeExportDetails:
		return "exportDetails"

	default:
		return ""
	}
}
