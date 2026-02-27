package middleware

import (
	"encoding/json"
	"slices"
)

type Resource string

const (
	APP           Resource = "APP"
	KEY           Resource = "KEY"
	OneTimeToken  Resource = "ONE_TIME_TOKEN"
	DELIVERY      Resource = "DELIVERY"
	RESOURCES_ALL Resource = "RESOURCES_ALL"
)

type Operation string

const (
	CREATE         Operation = "CREATE"
	READ           Operation = "READ"
	DESTROY        Operation = "DESTROY"
	UPDATE         Operation = "UPDATE"
	OPERATIONS_ALL Operation = "OPERATIONS_ALL"
)

type ResourcePermissionList map[Resource][]Operation
type PermissionsBuilder struct {
	Permissions ResourcePermissionList `json:"permissions"`
}

func ParsePermissions(token []byte) (*PermissionsBuilder, error) {
	Permissions := PermissionsBuilder{}
	err := json.Unmarshal(token, &Permissions)
	if err != nil {
		return nil, err
	}
	return &Permissions, nil
}

func (c *PermissionsBuilder) Can(r Resource, o Operation) bool {
	// does this token has a wildcard permissions
	if c.Permissions[RESOURCES_ALL] != nil && slices.Contains(c.Permissions[RESOURCES_ALL], OPERATIONS_ALL) {
		return true
	}

	if c.Permissions[RESOURCES_ALL] != nil && slices.Contains(c.Permissions[RESOURCES_ALL], o) {
		return true
	}

	if c.Permissions[r] == nil {
		return false
	}

	return slices.Contains(c.Permissions[r], o)
}

func (c *PermissionsBuilder) Allow(r Resource, o Operation) *PermissionsBuilder {
	if c.Permissions[r] == nil {
		c.Permissions[r] = []Operation{}
	}
	if slices.Index(c.Permissions[r], o) == -1 {
		c.Permissions[r] = append(c.Permissions[r], o)
	}
	return c
}
func (c *PermissionsBuilder) Deny(r Resource, o Operation) *PermissionsBuilder {
	if c.Permissions[r] == nil {
		c.Permissions[r] = []Operation{}
	}
	idx := slices.Index(c.Permissions[r], o)
	if idx == -1 {
		return c
	}
	c.Permissions[r] = slices.Delete(c.Permissions[r], idx, idx+1)
	return c
}

func (c *PermissionsBuilder) String() (string, error) {
	byts, err := json.Marshal(*c)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

func NewPermissionsBuilder() PermissionsBuilder {
	return PermissionsBuilder{
		Permissions: make(ResourcePermissionList),
	}
}
