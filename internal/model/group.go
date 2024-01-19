package model

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Group struct {
	Id        Uuid      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewGroup(name string) (Group, error) {
	var g Group

	uuid, err := NewUuid()
	if err != nil {
		return g, fmt.Errorf("cannot create group: %v", err)
	}

	g.Id = uuid
	g.Name = name
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	return g, nil
}

func (g *Group) Validate() error {
	errs := make([]error, 0)
	if g.Id == "" {
		errs = append(errs, errors.New("missing ID"))
	}
	if g.Name == "" {
		errs = append(errs, errors.New("missing name"))
	}
	if g.CreatedAt.IsZero() {
		errs = append(errs, errors.New("create time is missing a value"))
	}
	if g.UpdatedAt.IsZero() {
		errs = append(errs, errors.New("update time is missing a value"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

type GroupService interface {
	CreateGroup(ctx context.Context, p CreateGroupParams) (Group, error)
	UpdateGroup(ctx context.Context, p UpdateGroupParams) error
	RemoveGroup(ctx context.Context, p RemoveGroupParams) error
	ListGroups(ctx context.Context, p ListGroupsParams) ([]Group, error)
	FindGroupById(ctx context.Context, id Uuid) (Group, error)
	FindGroupByName(ctx context.Context, name string) (Group, error)
}

type CreateGroupParams struct {
	Name          string `json:"name"`
	OwnerName     string `json:"ownerName"`
	OwnerUsername string `json:"ownerUsername"`
}

type UpdateGroupParams struct {
	Id   Uuid    `json:"id"`
	Name *string `json:"name"`
}

type RemoveGroupParams struct {
	Id Uuid `json:"id"`
}

type ListGroupsParams struct {
	Pattern *string `json:"pattern"`
}
