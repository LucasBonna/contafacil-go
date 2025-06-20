// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/lucasbonna/contafacil_api/ent/clients"
	"github.com/lucasbonna/contafacil_api/ent/user"
)

// ClientsCreate is the builder for creating a Clients entity.
type ClientsCreate struct {
	config
	mutation *ClientsMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (cc *ClientsCreate) SetName(s string) *ClientsCreate {
	cc.mutation.SetName(s)
	return cc
}

// SetCnpj sets the "cnpj" field.
func (cc *ClientsCreate) SetCnpj(s string) *ClientsCreate {
	cc.mutation.SetCnpj(s)
	return cc
}

// SetRole sets the "role" field.
func (cc *ClientsCreate) SetRole(c clients.Role) *ClientsCreate {
	cc.mutation.SetRole(c)
	return cc
}

// SetCreatedAt sets the "created_at" field.
func (cc *ClientsCreate) SetCreatedAt(t time.Time) *ClientsCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cc *ClientsCreate) SetNillableCreatedAt(t *time.Time) *ClientsCreate {
	if t != nil {
		cc.SetCreatedAt(*t)
	}
	return cc
}

// SetUpdatedAt sets the "updated_at" field.
func (cc *ClientsCreate) SetUpdatedAt(t time.Time) *ClientsCreate {
	cc.mutation.SetUpdatedAt(t)
	return cc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (cc *ClientsCreate) SetNillableUpdatedAt(t *time.Time) *ClientsCreate {
	if t != nil {
		cc.SetUpdatedAt(*t)
	}
	return cc
}

// SetDeletedAt sets the "deleted_at" field.
func (cc *ClientsCreate) SetDeletedAt(t time.Time) *ClientsCreate {
	cc.mutation.SetDeletedAt(t)
	return cc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (cc *ClientsCreate) SetNillableDeletedAt(t *time.Time) *ClientsCreate {
	if t != nil {
		cc.SetDeletedAt(*t)
	}
	return cc
}

// SetID sets the "id" field.
func (cc *ClientsCreate) SetID(u uuid.UUID) *ClientsCreate {
	cc.mutation.SetID(u)
	return cc
}

// AddUserIDs adds the "users" edge to the User entity by IDs.
func (cc *ClientsCreate) AddUserIDs(ids ...uuid.UUID) *ClientsCreate {
	cc.mutation.AddUserIDs(ids...)
	return cc
}

// AddUsers adds the "users" edges to the User entity.
func (cc *ClientsCreate) AddUsers(u ...*User) *ClientsCreate {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return cc.AddUserIDs(ids...)
}

// Mutation returns the ClientsMutation object of the builder.
func (cc *ClientsCreate) Mutation() *ClientsMutation {
	return cc.mutation
}

// Save creates the Clients in the database.
func (cc *ClientsCreate) Save(ctx context.Context) (*Clients, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ClientsCreate) SaveX(ctx context.Context) *Clients {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ClientsCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ClientsCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *ClientsCreate) defaults() {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		v := clients.DefaultCreatedAt()
		cc.mutation.SetCreatedAt(v)
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		v := clients.DefaultUpdatedAt()
		cc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ClientsCreate) check() error {
	if _, ok := cc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Clients.name"`)}
	}
	if v, ok := cc.mutation.Name(); ok {
		if err := clients.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Clients.name": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Cnpj(); !ok {
		return &ValidationError{Name: "cnpj", err: errors.New(`ent: missing required field "Clients.cnpj"`)}
	}
	if v, ok := cc.mutation.Cnpj(); ok {
		if err := clients.CnpjValidator(v); err != nil {
			return &ValidationError{Name: "cnpj", err: fmt.Errorf(`ent: validator failed for field "Clients.cnpj": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Role(); !ok {
		return &ValidationError{Name: "role", err: errors.New(`ent: missing required field "Clients.role"`)}
	}
	if v, ok := cc.mutation.Role(); ok {
		if err := clients.RoleValidator(v); err != nil {
			return &ValidationError{Name: "role", err: fmt.Errorf(`ent: validator failed for field "Clients.role": %w`, err)}
		}
	}
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Clients.created_at"`)}
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Clients.updated_at"`)}
	}
	return nil
}

func (cc *ClientsCreate) sqlSave(ctx context.Context) (*Clients, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *ClientsCreate) createSpec() (*Clients, *sqlgraph.CreateSpec) {
	var (
		_node = &Clients{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(clients.Table, sqlgraph.NewFieldSpec(clients.FieldID, field.TypeUUID))
	)
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cc.mutation.Name(); ok {
		_spec.SetField(clients.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := cc.mutation.Cnpj(); ok {
		_spec.SetField(clients.FieldCnpj, field.TypeString, value)
		_node.Cnpj = value
	}
	if value, ok := cc.mutation.Role(); ok {
		_spec.SetField(clients.FieldRole, field.TypeEnum, value)
		_node.Role = value
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.SetField(clients.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := cc.mutation.UpdatedAt(); ok {
		_spec.SetField(clients.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := cc.mutation.DeletedAt(); ok {
		_spec.SetField(clients.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = value
	}
	if nodes := cc.mutation.UsersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   clients.UsersTable,
			Columns: []string{clients.UsersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ClientsCreateBulk is the builder for creating many Clients entities in bulk.
type ClientsCreateBulk struct {
	config
	err      error
	builders []*ClientsCreate
}

// Save creates the Clients entities in the database.
func (ccb *ClientsCreateBulk) Save(ctx context.Context) ([]*Clients, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Clients, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ClientsMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ClientsCreateBulk) SaveX(ctx context.Context) []*Clients {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ClientsCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ClientsCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
