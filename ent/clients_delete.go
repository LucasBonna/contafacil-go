// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lucasbonna/contafacil_api/ent/clients"
	"github.com/lucasbonna/contafacil_api/ent/predicate"
)

// ClientsDelete is the builder for deleting a Clients entity.
type ClientsDelete struct {
	config
	hooks    []Hook
	mutation *ClientsMutation
}

// Where appends a list predicates to the ClientsDelete builder.
func (cd *ClientsDelete) Where(ps ...predicate.Clients) *ClientsDelete {
	cd.mutation.Where(ps...)
	return cd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cd *ClientsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cd.sqlExec, cd.mutation, cd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cd *ClientsDelete) ExecX(ctx context.Context) int {
	n, err := cd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cd *ClientsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(clients.Table, sqlgraph.NewFieldSpec(clients.FieldID, field.TypeUUID))
	if ps := cd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cd.mutation.done = true
	return affected, err
}

// ClientsDeleteOne is the builder for deleting a single Clients entity.
type ClientsDeleteOne struct {
	cd *ClientsDelete
}

// Where appends a list predicates to the ClientsDelete builder.
func (cdo *ClientsDeleteOne) Where(ps ...predicate.Clients) *ClientsDeleteOne {
	cdo.cd.mutation.Where(ps...)
	return cdo
}

// Exec executes the deletion query.
func (cdo *ClientsDeleteOne) Exec(ctx context.Context) error {
	n, err := cdo.cd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{clients.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cdo *ClientsDeleteOne) ExecX(ctx context.Context) {
	if err := cdo.Exec(ctx); err != nil {
		panic(err)
	}
}
