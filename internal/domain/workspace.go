package domain

import (
	"strings"
	"time"
	"unicode"
)

// WorkspaceRole defines the possible roles within a workspace.
type WorkspaceRole string

const (
	// WorkspaceRoleAdmin can manage members and projects
	WorkspaceRoleAdmin WorkspaceRole = "admin"
	// WorkspaceRoleMember can view and create tasks
	WorkspaceRoleMember WorkspaceRole = "member"
	// WorkspaceRoleViewer can only view
	WorkspaceRoleViewer WorkspaceRole = "viewer"
)

// IsValid checks that the role is one of the allowed values
func (r WorkspaceRole) IsValid() bool {
	switch r {
	case WorkspaceRoleAdmin, WorkspaceRoleMember, WorkspaceRoleViewer:
		return true
	default:
		return false
	}
}

// Workspace represents a workspace that groups users and projects
type Workspace struct {
	ID          string
	Name        string
	Slug        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WorkspaceMember represents a user's membership in a workspace
type WorkspaceMember struct {
	ID          string
	WorkspaceID string
	UserID      string
	Role        WorkspaceRole
	JoinedAt    time.Time
}

// NewWorkspace builds a Workspace with an automatically generated slug
func NewWorkspace(name, description, createdBy string) *Workspace {
	now := time.Now()
	return &Workspace{
		Name:        name,
		Slug:        GenerateSlug(name),
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GenerateSlug converts a name into a URL-friendly slug. "My Dev Team" -> "my-dev-team"
func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	var result strings.Builder

	for _, c := range slug {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			result.WriteRune(c)
		} else {
			result.WriteRune('-')
		}
	}

	cleaned := strings.Join(strings.FieldsFunc(result.String(), func(r rune) bool {
		return r == '-'
	}), "-")

	return cleaned
}

// NewWorkspaceMember builds a membership with a default role
func NewWorkspaceMember(workspaceID, userID string, role WorkspaceRole) *WorkspaceMember {
	return &WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now(),
	}
}
