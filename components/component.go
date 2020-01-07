package components

import (
	"database/sql"
	"errors"
	"path"
	"time"
)

// Service is a component type which represents a long-running service that must be available as
// part of a simplex data processing flow
var Service = "service"

// Task is a component type which represents a process that must be run to completion as part of a
// simplex data processing flow
var Task = "task"

// ComponentTypes is a set (of keys) enumerating the types of components that simplex respects
var ComponentTypes = map[string]bool{
	Service: true,
	Task:    true,
}

// ErrInvalidComponentType signifies that a caller attempted to create component metadata with
// a component type which wasn't included in the ComponentTypes map
var ErrInvalidComponentType = errors.New("Invalid ComponentType")

// ErrEmptyID signifies that a caller attempted to create component metadata in which the ID string
// was the empty string
var ErrEmptyID = errors.New("ID must be a non-empty string")

// ErrEmptyComponentPath signifies that a caller attempted to create component metadata in which the
// ComponentPath string was the empty string
var ErrEmptyComponentPath = errors.New("ComponentPath must be a non-empty string")

// ComponentMetadata - the metadata about a component that gets stored in the state database
type ComponentMetadata struct {
	ID                string    `json:"id"`
	ComponentType     string    `json:"component_type"`
	ComponentPath     string    `json:"component_path"`
	SpecificationPath string    `json:"specification_path"`
	CreatedAt         time.Time `json:"created_at"`
}

// DefaultSpecificationFileName - this is the name of the file inside the component directory
// representing the simplex specification of the component.
var DefaultSpecificationFileName = "component.json"

// GenerateComponentMetadata creates a ComponentMetadata instance from the specified parameters,
// applying defaults as required and reasonable. It also performs validation on its inputs and
// returns an error describing the reasons for rejection of invalid component metadata. Component
// metadata requires that:
// 1. id be non-null (ErrEmptyID returned otherwise)
// 2. componentType be one of the keys of the ComponentTypes map (ErrInvalidComponentType returned otherwise)
// 3. componentPath be non-empty (ErrEmptyComponentPath returned otherwise)
func GenerateComponentMetadata(id, componentType, componentPath, specificationPath string) (ComponentMetadata, error) {
	if id == "" {
		return ComponentMetadata{}, ErrEmptyID
	}

	if componentPath == "" {
		return ComponentMetadata{}, ErrEmptyComponentPath
	}

	if _, ok := ComponentTypes[componentType]; !ok {
		return ComponentMetadata{}, ErrInvalidComponentType
	}

	if specificationPath == "" {
		specificationPath = path.Join(componentPath, DefaultSpecificationFileName)
	}

	createdAt := time.Now()

	metadata := ComponentMetadata{
		ID:                id,
		ComponentType:     componentType,
		ComponentPath:     componentPath,
		SpecificationPath: specificationPath,
		CreatedAt:         createdAt,
	}
	return metadata, nil
}

// AddComponent registers a component (by metadata) against a simplex state database. It applies
// reasonable defaults where possible (e.g. on SpecificationPath).
// This is the handler for `simplex components add`
func AddComponent(db *sql.DB, id, componentType, componentPath, specificationPath string) (ComponentMetadata, error) {
	metadata, err := GenerateComponentMetadata(id, componentType, componentPath, specificationPath)
	if err != nil {
		return metadata, err
	}

	err = InsertComponent(db, metadata)

	return metadata, err
}