// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"doublequote/ent/collection"
	"doublequote/ent/collectionentry"
	"doublequote/ent/entry"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CollectionEntryCreate is the builder for creating a CollectionEntry entity.
type CollectionEntryCreate struct {
	config
	mutation *CollectionEntryMutation
	hooks    []Hook
}

// SetIsRead sets the "is_read" field.
func (cec *CollectionEntryCreate) SetIsRead(b bool) *CollectionEntryCreate {
	cec.mutation.SetIsRead(b)
	return cec
}

// SetNillableIsRead sets the "is_read" field if the given value is not nil.
func (cec *CollectionEntryCreate) SetNillableIsRead(b *bool) *CollectionEntryCreate {
	if b != nil {
		cec.SetIsRead(*b)
	}
	return cec
}

// SetCollectionID sets the "collection_id" field.
func (cec *CollectionEntryCreate) SetCollectionID(i int) *CollectionEntryCreate {
	cec.mutation.SetCollectionID(i)
	return cec
}

// SetEntryID sets the "entry_id" field.
func (cec *CollectionEntryCreate) SetEntryID(i int) *CollectionEntryCreate {
	cec.mutation.SetEntryID(i)
	return cec
}

// SetCreatedAt sets the "created_at" field.
func (cec *CollectionEntryCreate) SetCreatedAt(t time.Time) *CollectionEntryCreate {
	cec.mutation.SetCreatedAt(t)
	return cec
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cec *CollectionEntryCreate) SetNillableCreatedAt(t *time.Time) *CollectionEntryCreate {
	if t != nil {
		cec.SetCreatedAt(*t)
	}
	return cec
}

// SetUpdatedAt sets the "updated_at" field.
func (cec *CollectionEntryCreate) SetUpdatedAt(t time.Time) *CollectionEntryCreate {
	cec.mutation.SetUpdatedAt(t)
	return cec
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (cec *CollectionEntryCreate) SetNillableUpdatedAt(t *time.Time) *CollectionEntryCreate {
	if t != nil {
		cec.SetUpdatedAt(*t)
	}
	return cec
}

// SetCollection sets the "collection" edge to the Collection entity.
func (cec *CollectionEntryCreate) SetCollection(c *Collection) *CollectionEntryCreate {
	return cec.SetCollectionID(c.ID)
}

// SetEntry sets the "entry" edge to the Entry entity.
func (cec *CollectionEntryCreate) SetEntry(e *Entry) *CollectionEntryCreate {
	return cec.SetEntryID(e.ID)
}

// Mutation returns the CollectionEntryMutation object of the builder.
func (cec *CollectionEntryCreate) Mutation() *CollectionEntryMutation {
	return cec.mutation
}

// Save creates the CollectionEntry in the database.
func (cec *CollectionEntryCreate) Save(ctx context.Context) (*CollectionEntry, error) {
	var (
		err  error
		node *CollectionEntry
	)
	cec.defaults()
	if len(cec.hooks) == 0 {
		if err = cec.check(); err != nil {
			return nil, err
		}
		node, err = cec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CollectionEntryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cec.check(); err != nil {
				return nil, err
			}
			cec.mutation = mutation
			if node, err = cec.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(cec.hooks) - 1; i >= 0; i-- {
			if cec.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cec.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cec.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (cec *CollectionEntryCreate) SaveX(ctx context.Context) *CollectionEntry {
	v, err := cec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cec *CollectionEntryCreate) Exec(ctx context.Context) error {
	_, err := cec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cec *CollectionEntryCreate) ExecX(ctx context.Context) {
	if err := cec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cec *CollectionEntryCreate) defaults() {
	if _, ok := cec.mutation.IsRead(); !ok {
		v := collectionentry.DefaultIsRead
		cec.mutation.SetIsRead(v)
	}
	if _, ok := cec.mutation.CreatedAt(); !ok {
		v := collectionentry.DefaultCreatedAt()
		cec.mutation.SetCreatedAt(v)
	}
	if _, ok := cec.mutation.UpdatedAt(); !ok {
		v := collectionentry.DefaultUpdatedAt()
		cec.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cec *CollectionEntryCreate) check() error {
	if _, ok := cec.mutation.IsRead(); !ok {
		return &ValidationError{Name: "is_read", err: errors.New(`ent: missing required field "CollectionEntry.is_read"`)}
	}
	if _, ok := cec.mutation.CollectionID(); !ok {
		return &ValidationError{Name: "collection_id", err: errors.New(`ent: missing required field "CollectionEntry.collection_id"`)}
	}
	if _, ok := cec.mutation.EntryID(); !ok {
		return &ValidationError{Name: "entry_id", err: errors.New(`ent: missing required field "CollectionEntry.entry_id"`)}
	}
	if _, ok := cec.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "CollectionEntry.created_at"`)}
	}
	if _, ok := cec.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "CollectionEntry.updated_at"`)}
	}
	if _, ok := cec.mutation.CollectionID(); !ok {
		return &ValidationError{Name: "collection", err: errors.New(`ent: missing required edge "CollectionEntry.collection"`)}
	}
	if _, ok := cec.mutation.EntryID(); !ok {
		return &ValidationError{Name: "entry", err: errors.New(`ent: missing required edge "CollectionEntry.entry"`)}
	}
	return nil
}

func (cec *CollectionEntryCreate) sqlSave(ctx context.Context) (*CollectionEntry, error) {
	_node, _spec := cec.createSpec()
	if err := sqlgraph.CreateNode(ctx, cec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (cec *CollectionEntryCreate) createSpec() (*CollectionEntry, *sqlgraph.CreateSpec) {
	var (
		_node = &CollectionEntry{config: cec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: collectionentry.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: collectionentry.FieldID,
			},
		}
	)
	if value, ok := cec.mutation.IsRead(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: collectionentry.FieldIsRead,
		})
		_node.IsRead = value
	}
	if value, ok := cec.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: collectionentry.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := cec.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: collectionentry.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if nodes := cec.mutation.CollectionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   collectionentry.CollectionTable,
			Columns: []string{collectionentry.CollectionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: collection.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CollectionID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cec.mutation.EntryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   collectionentry.EntryTable,
			Columns: []string{collectionentry.EntryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: entry.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.EntryID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CollectionEntryCreateBulk is the builder for creating many CollectionEntry entities in bulk.
type CollectionEntryCreateBulk struct {
	config
	builders []*CollectionEntryCreate
}

// Save creates the CollectionEntry entities in the database.
func (cecb *CollectionEntryCreateBulk) Save(ctx context.Context) ([]*CollectionEntry, error) {
	specs := make([]*sqlgraph.CreateSpec, len(cecb.builders))
	nodes := make([]*CollectionEntry, len(cecb.builders))
	mutators := make([]Mutator, len(cecb.builders))
	for i := range cecb.builders {
		func(i int, root context.Context) {
			builder := cecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CollectionEntryMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, cecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, cecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cecb *CollectionEntryCreateBulk) SaveX(ctx context.Context) []*CollectionEntry {
	v, err := cecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cecb *CollectionEntryCreateBulk) Exec(ctx context.Context) error {
	_, err := cecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cecb *CollectionEntryCreateBulk) ExecX(ctx context.Context) {
	if err := cecb.Exec(ctx); err != nil {
		panic(err)
	}
}
