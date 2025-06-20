// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lucasbonna/contafacil_api/ent/accesslog"
	"github.com/lucasbonna/contafacil_api/ent/predicate"
)

// AccessLogUpdate is the builder for updating AccessLog entities.
type AccessLogUpdate struct {
	config
	hooks    []Hook
	mutation *AccessLogMutation
}

// Where appends a list predicates to the AccessLogUpdate builder.
func (alu *AccessLogUpdate) Where(ps ...predicate.AccessLog) *AccessLogUpdate {
	alu.mutation.Where(ps...)
	return alu
}

// SetIP sets the "ip" field.
func (alu *AccessLogUpdate) SetIP(s string) *AccessLogUpdate {
	alu.mutation.SetIP(s)
	return alu
}

// SetNillableIP sets the "ip" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableIP(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetIP(*s)
	}
	return alu
}

// SetMethod sets the "method" field.
func (alu *AccessLogUpdate) SetMethod(s string) *AccessLogUpdate {
	alu.mutation.SetMethod(s)
	return alu
}

// SetNillableMethod sets the "method" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableMethod(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetMethod(*s)
	}
	return alu
}

// SetEndpoint sets the "endpoint" field.
func (alu *AccessLogUpdate) SetEndpoint(s string) *AccessLogUpdate {
	alu.mutation.SetEndpoint(s)
	return alu
}

// SetNillableEndpoint sets the "endpoint" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableEndpoint(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetEndpoint(*s)
	}
	return alu
}

// SetRequestBody sets the "request_body" field.
func (alu *AccessLogUpdate) SetRequestBody(s string) *AccessLogUpdate {
	alu.mutation.SetRequestBody(s)
	return alu
}

// SetNillableRequestBody sets the "request_body" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableRequestBody(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetRequestBody(*s)
	}
	return alu
}

// ClearRequestBody clears the value of the "request_body" field.
func (alu *AccessLogUpdate) ClearRequestBody() *AccessLogUpdate {
	alu.mutation.ClearRequestBody()
	return alu
}

// SetRequestHeaders sets the "request_headers" field.
func (alu *AccessLogUpdate) SetRequestHeaders(s string) *AccessLogUpdate {
	alu.mutation.SetRequestHeaders(s)
	return alu
}

// SetNillableRequestHeaders sets the "request_headers" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableRequestHeaders(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetRequestHeaders(*s)
	}
	return alu
}

// ClearRequestHeaders clears the value of the "request_headers" field.
func (alu *AccessLogUpdate) ClearRequestHeaders() *AccessLogUpdate {
	alu.mutation.ClearRequestHeaders()
	return alu
}

// SetRequestParams sets the "request_params" field.
func (alu *AccessLogUpdate) SetRequestParams(s string) *AccessLogUpdate {
	alu.mutation.SetRequestParams(s)
	return alu
}

// SetNillableRequestParams sets the "request_params" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableRequestParams(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetRequestParams(*s)
	}
	return alu
}

// ClearRequestParams clears the value of the "request_params" field.
func (alu *AccessLogUpdate) ClearRequestParams() *AccessLogUpdate {
	alu.mutation.ClearRequestParams()
	return alu
}

// SetRequestQuery sets the "request_query" field.
func (alu *AccessLogUpdate) SetRequestQuery(s string) *AccessLogUpdate {
	alu.mutation.SetRequestQuery(s)
	return alu
}

// SetNillableRequestQuery sets the "request_query" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableRequestQuery(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetRequestQuery(*s)
	}
	return alu
}

// ClearRequestQuery clears the value of the "request_query" field.
func (alu *AccessLogUpdate) ClearRequestQuery() *AccessLogUpdate {
	alu.mutation.ClearRequestQuery()
	return alu
}

// SetResponseBody sets the "response_body" field.
func (alu *AccessLogUpdate) SetResponseBody(s string) *AccessLogUpdate {
	alu.mutation.SetResponseBody(s)
	return alu
}

// SetNillableResponseBody sets the "response_body" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableResponseBody(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetResponseBody(*s)
	}
	return alu
}

// ClearResponseBody clears the value of the "response_body" field.
func (alu *AccessLogUpdate) ClearResponseBody() *AccessLogUpdate {
	alu.mutation.ClearResponseBody()
	return alu
}

// SetResponseHeaders sets the "response_headers" field.
func (alu *AccessLogUpdate) SetResponseHeaders(s string) *AccessLogUpdate {
	alu.mutation.SetResponseHeaders(s)
	return alu
}

// SetNillableResponseHeaders sets the "response_headers" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableResponseHeaders(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetResponseHeaders(*s)
	}
	return alu
}

// ClearResponseHeaders clears the value of the "response_headers" field.
func (alu *AccessLogUpdate) ClearResponseHeaders() *AccessLogUpdate {
	alu.mutation.ClearResponseHeaders()
	return alu
}

// SetResponseTime sets the "response_time" field.
func (alu *AccessLogUpdate) SetResponseTime(s string) *AccessLogUpdate {
	alu.mutation.SetResponseTime(s)
	return alu
}

// SetNillableResponseTime sets the "response_time" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableResponseTime(s *string) *AccessLogUpdate {
	if s != nil {
		alu.SetResponseTime(*s)
	}
	return alu
}

// ClearResponseTime clears the value of the "response_time" field.
func (alu *AccessLogUpdate) ClearResponseTime() *AccessLogUpdate {
	alu.mutation.ClearResponseTime()
	return alu
}

// SetStatusCode sets the "status_code" field.
func (alu *AccessLogUpdate) SetStatusCode(i int) *AccessLogUpdate {
	alu.mutation.ResetStatusCode()
	alu.mutation.SetStatusCode(i)
	return alu
}

// SetNillableStatusCode sets the "status_code" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableStatusCode(i *int) *AccessLogUpdate {
	if i != nil {
		alu.SetStatusCode(*i)
	}
	return alu
}

// AddStatusCode adds i to the "status_code" field.
func (alu *AccessLogUpdate) AddStatusCode(i int) *AccessLogUpdate {
	alu.mutation.AddStatusCode(i)
	return alu
}

// ClearStatusCode clears the value of the "status_code" field.
func (alu *AccessLogUpdate) ClearStatusCode() *AccessLogUpdate {
	alu.mutation.ClearStatusCode()
	return alu
}

// SetCreatedAt sets the "created_at" field.
func (alu *AccessLogUpdate) SetCreatedAt(t time.Time) *AccessLogUpdate {
	alu.mutation.SetCreatedAt(t)
	return alu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableCreatedAt(t *time.Time) *AccessLogUpdate {
	if t != nil {
		alu.SetCreatedAt(*t)
	}
	return alu
}

// SetUpdatedAt sets the "updated_at" field.
func (alu *AccessLogUpdate) SetUpdatedAt(t time.Time) *AccessLogUpdate {
	alu.mutation.SetUpdatedAt(t)
	return alu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (alu *AccessLogUpdate) SetNillableUpdatedAt(t *time.Time) *AccessLogUpdate {
	if t != nil {
		alu.SetUpdatedAt(*t)
	}
	return alu
}

// Mutation returns the AccessLogMutation object of the builder.
func (alu *AccessLogUpdate) Mutation() *AccessLogMutation {
	return alu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (alu *AccessLogUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, alu.sqlSave, alu.mutation, alu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (alu *AccessLogUpdate) SaveX(ctx context.Context) int {
	affected, err := alu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (alu *AccessLogUpdate) Exec(ctx context.Context) error {
	_, err := alu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (alu *AccessLogUpdate) ExecX(ctx context.Context) {
	if err := alu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (alu *AccessLogUpdate) check() error {
	if v, ok := alu.mutation.ResponseTime(); ok {
		if err := accesslog.ResponseTimeValidator(v); err != nil {
			return &ValidationError{Name: "response_time", err: fmt.Errorf(`ent: validator failed for field "AccessLog.response_time": %w`, err)}
		}
	}
	return nil
}

func (alu *AccessLogUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := alu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(accesslog.Table, accesslog.Columns, sqlgraph.NewFieldSpec(accesslog.FieldID, field.TypeUUID))
	if ps := alu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := alu.mutation.IP(); ok {
		_spec.SetField(accesslog.FieldIP, field.TypeString, value)
	}
	if value, ok := alu.mutation.Method(); ok {
		_spec.SetField(accesslog.FieldMethod, field.TypeString, value)
	}
	if value, ok := alu.mutation.Endpoint(); ok {
		_spec.SetField(accesslog.FieldEndpoint, field.TypeString, value)
	}
	if value, ok := alu.mutation.RequestBody(); ok {
		_spec.SetField(accesslog.FieldRequestBody, field.TypeString, value)
	}
	if alu.mutation.RequestBodyCleared() {
		_spec.ClearField(accesslog.FieldRequestBody, field.TypeString)
	}
	if value, ok := alu.mutation.RequestHeaders(); ok {
		_spec.SetField(accesslog.FieldRequestHeaders, field.TypeString, value)
	}
	if alu.mutation.RequestHeadersCleared() {
		_spec.ClearField(accesslog.FieldRequestHeaders, field.TypeString)
	}
	if value, ok := alu.mutation.RequestParams(); ok {
		_spec.SetField(accesslog.FieldRequestParams, field.TypeString, value)
	}
	if alu.mutation.RequestParamsCleared() {
		_spec.ClearField(accesslog.FieldRequestParams, field.TypeString)
	}
	if value, ok := alu.mutation.RequestQuery(); ok {
		_spec.SetField(accesslog.FieldRequestQuery, field.TypeString, value)
	}
	if alu.mutation.RequestQueryCleared() {
		_spec.ClearField(accesslog.FieldRequestQuery, field.TypeString)
	}
	if value, ok := alu.mutation.ResponseBody(); ok {
		_spec.SetField(accesslog.FieldResponseBody, field.TypeString, value)
	}
	if alu.mutation.ResponseBodyCleared() {
		_spec.ClearField(accesslog.FieldResponseBody, field.TypeString)
	}
	if value, ok := alu.mutation.ResponseHeaders(); ok {
		_spec.SetField(accesslog.FieldResponseHeaders, field.TypeString, value)
	}
	if alu.mutation.ResponseHeadersCleared() {
		_spec.ClearField(accesslog.FieldResponseHeaders, field.TypeString)
	}
	if value, ok := alu.mutation.ResponseTime(); ok {
		_spec.SetField(accesslog.FieldResponseTime, field.TypeString, value)
	}
	if alu.mutation.ResponseTimeCleared() {
		_spec.ClearField(accesslog.FieldResponseTime, field.TypeString)
	}
	if value, ok := alu.mutation.StatusCode(); ok {
		_spec.SetField(accesslog.FieldStatusCode, field.TypeInt, value)
	}
	if value, ok := alu.mutation.AddedStatusCode(); ok {
		_spec.AddField(accesslog.FieldStatusCode, field.TypeInt, value)
	}
	if alu.mutation.StatusCodeCleared() {
		_spec.ClearField(accesslog.FieldStatusCode, field.TypeInt)
	}
	if value, ok := alu.mutation.CreatedAt(); ok {
		_spec.SetField(accesslog.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := alu.mutation.UpdatedAt(); ok {
		_spec.SetField(accesslog.FieldUpdatedAt, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, alu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{accesslog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	alu.mutation.done = true
	return n, nil
}

// AccessLogUpdateOne is the builder for updating a single AccessLog entity.
type AccessLogUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AccessLogMutation
}

// SetIP sets the "ip" field.
func (aluo *AccessLogUpdateOne) SetIP(s string) *AccessLogUpdateOne {
	aluo.mutation.SetIP(s)
	return aluo
}

// SetNillableIP sets the "ip" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableIP(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetIP(*s)
	}
	return aluo
}

// SetMethod sets the "method" field.
func (aluo *AccessLogUpdateOne) SetMethod(s string) *AccessLogUpdateOne {
	aluo.mutation.SetMethod(s)
	return aluo
}

// SetNillableMethod sets the "method" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableMethod(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetMethod(*s)
	}
	return aluo
}

// SetEndpoint sets the "endpoint" field.
func (aluo *AccessLogUpdateOne) SetEndpoint(s string) *AccessLogUpdateOne {
	aluo.mutation.SetEndpoint(s)
	return aluo
}

// SetNillableEndpoint sets the "endpoint" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableEndpoint(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetEndpoint(*s)
	}
	return aluo
}

// SetRequestBody sets the "request_body" field.
func (aluo *AccessLogUpdateOne) SetRequestBody(s string) *AccessLogUpdateOne {
	aluo.mutation.SetRequestBody(s)
	return aluo
}

// SetNillableRequestBody sets the "request_body" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableRequestBody(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetRequestBody(*s)
	}
	return aluo
}

// ClearRequestBody clears the value of the "request_body" field.
func (aluo *AccessLogUpdateOne) ClearRequestBody() *AccessLogUpdateOne {
	aluo.mutation.ClearRequestBody()
	return aluo
}

// SetRequestHeaders sets the "request_headers" field.
func (aluo *AccessLogUpdateOne) SetRequestHeaders(s string) *AccessLogUpdateOne {
	aluo.mutation.SetRequestHeaders(s)
	return aluo
}

// SetNillableRequestHeaders sets the "request_headers" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableRequestHeaders(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetRequestHeaders(*s)
	}
	return aluo
}

// ClearRequestHeaders clears the value of the "request_headers" field.
func (aluo *AccessLogUpdateOne) ClearRequestHeaders() *AccessLogUpdateOne {
	aluo.mutation.ClearRequestHeaders()
	return aluo
}

// SetRequestParams sets the "request_params" field.
func (aluo *AccessLogUpdateOne) SetRequestParams(s string) *AccessLogUpdateOne {
	aluo.mutation.SetRequestParams(s)
	return aluo
}

// SetNillableRequestParams sets the "request_params" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableRequestParams(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetRequestParams(*s)
	}
	return aluo
}

// ClearRequestParams clears the value of the "request_params" field.
func (aluo *AccessLogUpdateOne) ClearRequestParams() *AccessLogUpdateOne {
	aluo.mutation.ClearRequestParams()
	return aluo
}

// SetRequestQuery sets the "request_query" field.
func (aluo *AccessLogUpdateOne) SetRequestQuery(s string) *AccessLogUpdateOne {
	aluo.mutation.SetRequestQuery(s)
	return aluo
}

// SetNillableRequestQuery sets the "request_query" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableRequestQuery(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetRequestQuery(*s)
	}
	return aluo
}

// ClearRequestQuery clears the value of the "request_query" field.
func (aluo *AccessLogUpdateOne) ClearRequestQuery() *AccessLogUpdateOne {
	aluo.mutation.ClearRequestQuery()
	return aluo
}

// SetResponseBody sets the "response_body" field.
func (aluo *AccessLogUpdateOne) SetResponseBody(s string) *AccessLogUpdateOne {
	aluo.mutation.SetResponseBody(s)
	return aluo
}

// SetNillableResponseBody sets the "response_body" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableResponseBody(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetResponseBody(*s)
	}
	return aluo
}

// ClearResponseBody clears the value of the "response_body" field.
func (aluo *AccessLogUpdateOne) ClearResponseBody() *AccessLogUpdateOne {
	aluo.mutation.ClearResponseBody()
	return aluo
}

// SetResponseHeaders sets the "response_headers" field.
func (aluo *AccessLogUpdateOne) SetResponseHeaders(s string) *AccessLogUpdateOne {
	aluo.mutation.SetResponseHeaders(s)
	return aluo
}

// SetNillableResponseHeaders sets the "response_headers" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableResponseHeaders(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetResponseHeaders(*s)
	}
	return aluo
}

// ClearResponseHeaders clears the value of the "response_headers" field.
func (aluo *AccessLogUpdateOne) ClearResponseHeaders() *AccessLogUpdateOne {
	aluo.mutation.ClearResponseHeaders()
	return aluo
}

// SetResponseTime sets the "response_time" field.
func (aluo *AccessLogUpdateOne) SetResponseTime(s string) *AccessLogUpdateOne {
	aluo.mutation.SetResponseTime(s)
	return aluo
}

// SetNillableResponseTime sets the "response_time" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableResponseTime(s *string) *AccessLogUpdateOne {
	if s != nil {
		aluo.SetResponseTime(*s)
	}
	return aluo
}

// ClearResponseTime clears the value of the "response_time" field.
func (aluo *AccessLogUpdateOne) ClearResponseTime() *AccessLogUpdateOne {
	aluo.mutation.ClearResponseTime()
	return aluo
}

// SetStatusCode sets the "status_code" field.
func (aluo *AccessLogUpdateOne) SetStatusCode(i int) *AccessLogUpdateOne {
	aluo.mutation.ResetStatusCode()
	aluo.mutation.SetStatusCode(i)
	return aluo
}

// SetNillableStatusCode sets the "status_code" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableStatusCode(i *int) *AccessLogUpdateOne {
	if i != nil {
		aluo.SetStatusCode(*i)
	}
	return aluo
}

// AddStatusCode adds i to the "status_code" field.
func (aluo *AccessLogUpdateOne) AddStatusCode(i int) *AccessLogUpdateOne {
	aluo.mutation.AddStatusCode(i)
	return aluo
}

// ClearStatusCode clears the value of the "status_code" field.
func (aluo *AccessLogUpdateOne) ClearStatusCode() *AccessLogUpdateOne {
	aluo.mutation.ClearStatusCode()
	return aluo
}

// SetCreatedAt sets the "created_at" field.
func (aluo *AccessLogUpdateOne) SetCreatedAt(t time.Time) *AccessLogUpdateOne {
	aluo.mutation.SetCreatedAt(t)
	return aluo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableCreatedAt(t *time.Time) *AccessLogUpdateOne {
	if t != nil {
		aluo.SetCreatedAt(*t)
	}
	return aluo
}

// SetUpdatedAt sets the "updated_at" field.
func (aluo *AccessLogUpdateOne) SetUpdatedAt(t time.Time) *AccessLogUpdateOne {
	aluo.mutation.SetUpdatedAt(t)
	return aluo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (aluo *AccessLogUpdateOne) SetNillableUpdatedAt(t *time.Time) *AccessLogUpdateOne {
	if t != nil {
		aluo.SetUpdatedAt(*t)
	}
	return aluo
}

// Mutation returns the AccessLogMutation object of the builder.
func (aluo *AccessLogUpdateOne) Mutation() *AccessLogMutation {
	return aluo.mutation
}

// Where appends a list predicates to the AccessLogUpdate builder.
func (aluo *AccessLogUpdateOne) Where(ps ...predicate.AccessLog) *AccessLogUpdateOne {
	aluo.mutation.Where(ps...)
	return aluo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (aluo *AccessLogUpdateOne) Select(field string, fields ...string) *AccessLogUpdateOne {
	aluo.fields = append([]string{field}, fields...)
	return aluo
}

// Save executes the query and returns the updated AccessLog entity.
func (aluo *AccessLogUpdateOne) Save(ctx context.Context) (*AccessLog, error) {
	return withHooks(ctx, aluo.sqlSave, aluo.mutation, aluo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aluo *AccessLogUpdateOne) SaveX(ctx context.Context) *AccessLog {
	node, err := aluo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (aluo *AccessLogUpdateOne) Exec(ctx context.Context) error {
	_, err := aluo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aluo *AccessLogUpdateOne) ExecX(ctx context.Context) {
	if err := aluo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (aluo *AccessLogUpdateOne) check() error {
	if v, ok := aluo.mutation.ResponseTime(); ok {
		if err := accesslog.ResponseTimeValidator(v); err != nil {
			return &ValidationError{Name: "response_time", err: fmt.Errorf(`ent: validator failed for field "AccessLog.response_time": %w`, err)}
		}
	}
	return nil
}

func (aluo *AccessLogUpdateOne) sqlSave(ctx context.Context) (_node *AccessLog, err error) {
	if err := aluo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(accesslog.Table, accesslog.Columns, sqlgraph.NewFieldSpec(accesslog.FieldID, field.TypeUUID))
	id, ok := aluo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "AccessLog.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := aluo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, accesslog.FieldID)
		for _, f := range fields {
			if !accesslog.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != accesslog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := aluo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aluo.mutation.IP(); ok {
		_spec.SetField(accesslog.FieldIP, field.TypeString, value)
	}
	if value, ok := aluo.mutation.Method(); ok {
		_spec.SetField(accesslog.FieldMethod, field.TypeString, value)
	}
	if value, ok := aluo.mutation.Endpoint(); ok {
		_spec.SetField(accesslog.FieldEndpoint, field.TypeString, value)
	}
	if value, ok := aluo.mutation.RequestBody(); ok {
		_spec.SetField(accesslog.FieldRequestBody, field.TypeString, value)
	}
	if aluo.mutation.RequestBodyCleared() {
		_spec.ClearField(accesslog.FieldRequestBody, field.TypeString)
	}
	if value, ok := aluo.mutation.RequestHeaders(); ok {
		_spec.SetField(accesslog.FieldRequestHeaders, field.TypeString, value)
	}
	if aluo.mutation.RequestHeadersCleared() {
		_spec.ClearField(accesslog.FieldRequestHeaders, field.TypeString)
	}
	if value, ok := aluo.mutation.RequestParams(); ok {
		_spec.SetField(accesslog.FieldRequestParams, field.TypeString, value)
	}
	if aluo.mutation.RequestParamsCleared() {
		_spec.ClearField(accesslog.FieldRequestParams, field.TypeString)
	}
	if value, ok := aluo.mutation.RequestQuery(); ok {
		_spec.SetField(accesslog.FieldRequestQuery, field.TypeString, value)
	}
	if aluo.mutation.RequestQueryCleared() {
		_spec.ClearField(accesslog.FieldRequestQuery, field.TypeString)
	}
	if value, ok := aluo.mutation.ResponseBody(); ok {
		_spec.SetField(accesslog.FieldResponseBody, field.TypeString, value)
	}
	if aluo.mutation.ResponseBodyCleared() {
		_spec.ClearField(accesslog.FieldResponseBody, field.TypeString)
	}
	if value, ok := aluo.mutation.ResponseHeaders(); ok {
		_spec.SetField(accesslog.FieldResponseHeaders, field.TypeString, value)
	}
	if aluo.mutation.ResponseHeadersCleared() {
		_spec.ClearField(accesslog.FieldResponseHeaders, field.TypeString)
	}
	if value, ok := aluo.mutation.ResponseTime(); ok {
		_spec.SetField(accesslog.FieldResponseTime, field.TypeString, value)
	}
	if aluo.mutation.ResponseTimeCleared() {
		_spec.ClearField(accesslog.FieldResponseTime, field.TypeString)
	}
	if value, ok := aluo.mutation.StatusCode(); ok {
		_spec.SetField(accesslog.FieldStatusCode, field.TypeInt, value)
	}
	if value, ok := aluo.mutation.AddedStatusCode(); ok {
		_spec.AddField(accesslog.FieldStatusCode, field.TypeInt, value)
	}
	if aluo.mutation.StatusCodeCleared() {
		_spec.ClearField(accesslog.FieldStatusCode, field.TypeInt)
	}
	if value, ok := aluo.mutation.CreatedAt(); ok {
		_spec.SetField(accesslog.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := aluo.mutation.UpdatedAt(); ok {
		_spec.SetField(accesslog.FieldUpdatedAt, field.TypeTime, value)
	}
	_node = &AccessLog{config: aluo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, aluo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{accesslog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	aluo.mutation.done = true
	return _node, nil
}
