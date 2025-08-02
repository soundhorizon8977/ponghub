package endpoint_type

type EndpointType string

const (
	// HEALTH represents a health check port
	HEALTH EndpointType = "health"

	// API represents an API port
	API EndpointType = "api"

	// UNKNOWN represents an unknown port type
	UNKNOWN EndpointType = "unknown"
)

// String returns the string representation of the EndpointType
func (pt EndpointType) String() string {
	switch pt {
	case HEALTH:
		return "health"
	case API:
		return "api"
	default:
		return "unknown"
	}
}

// IsValid checks if the EndpointType is valid
func (pt EndpointType) IsValid() bool {
	return pt == HEALTH || pt == API
}

// ParseEndpointType parses a string into a EndpointType
func ParseEndpointType(s string) EndpointType {
	switch s {
	case "health":
		return HEALTH
	case "api":
		return API
	default:
		return UNKNOWN
	}
}

// CheckEndpointType checks if the given port type is valid
func (pt EndpointType) CheckEndpointType() bool {
	return pt.IsValid()
}
