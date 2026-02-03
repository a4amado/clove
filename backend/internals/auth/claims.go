package auth

import (
	repository "clove/internals/services/generatedRepo"
	"encoding/json"
	"slices"
)

type ResourcePermessionList map[repository.Resource][]repository.Operation
type ClaimsBuilder struct {
	Permessions ResourcePermessionList `json:"Permessions"`
}

func ParseClaims(token string) (*ClaimsBuilder, error) {
	claims := ClaimsBuilder{}
	err := json.Unmarshal([]byte(token), &claims)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
func (c *ClaimsBuilder) Can(r repository.Resource, o repository.Operation) bool {
	if c.Permessions[r] == nil {
		return false
	}

	return slices.Contains(c.Permessions[r], o)
}

func (c *ClaimsBuilder) Allow(r repository.Resource, o repository.Operation) *ClaimsBuilder {
	if c.Permessions[r] == nil {
		c.Permessions[r] = []repository.Operation{}
	}
	if slices.Index(c.Permessions[r], o) == -1 {
		c.Permessions[r] = append(c.Permessions[r], o)
	}
	return c
}
func (c *ClaimsBuilder) Deny(r repository.Resource, o repository.Operation) *ClaimsBuilder {
	if c.Permessions[r] == nil {
		c.Permessions[r] = []repository.Operation{}
	}
	idx := slices.Index(c.Permessions[r], o)
	if idx == -1 {
		return c
	}
	c.Permessions[r] = slices.Delete(c.Permessions[r], idx, idx+1)
	return c
}

func (c *ClaimsBuilder) String() (string, error) {
	byts, err := json.Marshal(*c)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

func NewClaimsBuilder() ClaimsBuilder {
	return ClaimsBuilder{
		Permessions: make(ResourcePermessionList),
	}
}
