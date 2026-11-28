package domain

import "time"

// Project represent a project
type Project struct {
	ID          string
	WorkspaceID string
	Name        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProject constructor
func NewProject(workspaceID, name, description, createdBy string) *Project {
	now := time.Now()

	return &Project{
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
