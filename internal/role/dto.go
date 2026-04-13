package role

import (
	"slices"

	"github.com/cyradin/fixik/internal/db"
)

type Role struct {
	Code        Type
	Name        string
	Description string
	Permissions Permission
}

type Type = db.Role

const (
	User    = db.RoleUser
	Manager = db.RoleManager
	Admin   = db.RoleAdmin
)

var allRoles = []Role{
	{
		Name:        "Пользователь",
		Code:        User,
		Description: "Может работать с инцидентами (всё, кроме удаления)",
		Permissions: IncidentGet |
			IncidentCreate |
			IncidentUpdate |
			CommentGet |
			CommentCreate |
			CommentUpdate |
			UserGet |
			TeamGet |
			PriorityGet |
			StatusGet,
	},
	{
		Name:        "Менеджер",
		Code:        Manager,
		Description: "Может выполнять все операции, кроме изменения пользователей и команд",
		Permissions: IncidentGet |
			IncidentCreate |
			IncidentUpdate |
			IncidentDelete |
			CommentGet |
			CommentCreate |
			CommentUpdate |
			CommentDelete |
			UserGet |
			TeamGet |
			PriorityGet |
			PriorityCreate |
			PriorityUpdate |
			PriorityDelete |
			StatusGet |
			StatusCreate |
			StatusUpdate |
			StatusDelete,
	},
	{
		Name:        "Администратор",
		Code:        Admin,
		Description: "Может выполнять все операции",
		Permissions: IncidentGet |
			IncidentCreate |
			IncidentUpdate |
			IncidentDelete |
			CommentGet |
			CommentCreate |
			CommentUpdate |
			CommentDelete |
			UserGet |
			UserCreate |
			UserUpdate |
			UserDelete |
			TeamGet |
			TeamCreate |
			TeamUpdate |
			TeamDelete |
			PriorityGet |
			PriorityCreate |
			PriorityUpdate |
			PriorityDelete |
			StatusGet |
			StatusCreate |
			StatusUpdate |
			StatusDelete,
	},
}

func Can(code Type, permissions Permission) bool {
	for _, r := range allRoles {
		if r.Code == code {
			return r.Permissions&permissions > 0
		}
	}

	return false
}

func List() []Role {
	return slices.Clone(allRoles)
}

func Types() []Type {
	return []Type{
		User,
		Manager,
		Admin,
	}
}

type Permission uint64

const (
	// Incident
	IncidentGet Permission = 1 << iota
	IncidentCreate
	IncidentUpdate
	IncidentDelete

	// Comment
	CommentGet
	CommentCreate
	CommentUpdate
	CommentDelete

	// User
	UserGet
	UserCreate
	UserUpdate
	UserDelete

	// Team
	TeamGet
	TeamCreate
	TeamUpdate
	TeamDelete

	// Priority
	PriorityGet
	PriorityCreate
	PriorityUpdate
	PriorityDelete

	// Status
	StatusGet
	StatusCreate
	StatusUpdate
	StatusDelete
)

var permissionCodes = map[Permission]string{
	IncidentGet:    "INCIDENT_GET",
	IncidentCreate: "INCIDENT_CREATE",
	IncidentUpdate: "INCIDENT_UPDATE",
	IncidentDelete: "INCIDENT_DELETE",
	CommentGet:     "COMMENT_GET",
	CommentCreate:  "COMMENT_CREATE",
	CommentUpdate:  "COMMENT_UPDATE",
	CommentDelete:  "COMMENT_DELETE",
	UserGet:        "USER_GET",
	UserCreate:     "USER_CREATE",
	UserUpdate:     "USER_UPDATE",
	UserDelete:     "USER_DELETE",
	TeamGet:        "TEAM_GET",
	TeamCreate:     "TEAM_CREATE",
	TeamUpdate:     "TEAM_UPDATE",
	TeamDelete:     "TEAM_DELETE",
	PriorityGet:    "PRIORITY_GET",
	PriorityCreate: "PRIORITY_CREATE",
	PriorityUpdate: "PRIORITY_UPDATE",
	PriorityDelete: "PRIORITY_DELETE",
	StatusGet:      "STATUS_GET",
	StatusCreate:   "STATUS_CREATE",
	StatusUpdate:   "STATUS_UPDATE",
	StatusDelete:   "STATUS_DELETE",
}

func (p Permission) Codes() []string {
	var result []string

	for pp, code := range permissionCodes {
		if p&pp > 0 {
			result = append(result, code)
		}
	}

	return result
}
