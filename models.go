package azdevops

import "github.com/google/uuid"

type ResourceArea struct {
	Id          *uuid.UUID `json:"id,omitempty"`
	LocationUrl *string    `json:"locationUrl,omitempty"`
	Name        *string    `json:"name,omitempty"`
}

type Api struct {
	// Area name for this resource
	Area *string `json:"area,omitempty"`
	// Unique Identifier for this location
	Id *uuid.UUID `json:"id,omitempty"`
	// Maximum api version that this resource supports (current server version for this resource)
	MaxVersion *string `json:"maxVersion,omitempty"`
	// Minimum api version that this resource supports
	MinVersion *string `json:"minVersion,omitempty"`
	// The latest version of this resource location that is in "Release" (non-preview) mode
	ReleasedVersion *string `json:"releasedVersion,omitempty"`
	// Resource name
	ResourceName *string `json:"resourceName,omitempty"`
	// The current resource version supported by this resource location
	ResourceVersion *int `json:"resourceVersion,omitempty"`
	// This location's route template (templated relative path)
	RouteTemplate *string `json:"routeTemplate,omitempty"`
}
