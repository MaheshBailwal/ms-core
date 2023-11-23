package persistence

import (
	"context"
	"log"

	"github.com/MaheshBailwal/mscore/core"
)

type GenericRepository[TE core.IEntity] interface {
	All(ctx context.Context, filters QueryFilter) ([]TE, error)
	ById(ctx context.Context, id int64) (*TE, error)
	Create(ctx context.Context, entity *TE) (int64, error)
	Update(ctx context.Context, entity *TE) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Query(ctx context.Context, query IDomainQuery) (any, error)
}

type GenericRepositorySql[TE interface{}] struct {
	DBContext     *DBContext
	queryhandlers map[string]IQueryHandler
}

func NewGenericRepositorySql[TE interface{}](db *DBContext, queryhandlers map[string]IQueryHandler) GenericRepositorySql[TE] {
	return GenericRepositorySql[TE]{
		DBContext:     db,
		queryhandlers: queryhandlers,
	}
}

func (r *GenericRepositorySql[TE]) Create(ctx context.Context, entity *TE) (int64, error) {

	var query = CreateInsertQuery(*entity, r.DBContext.Schema)
	row := r.DBContext.ExecuteCommand(query)
	var id int64
	row.Scan(&id)
	return id, nil
}

func (r *GenericRepositorySql[TE]) Update(ctx context.Context, entity *TE) (int64, error) {
	var query = CreateUpdateQuery(*entity, r.DBContext.Schema)
	r.DBContext.ExecuteCommand(query)
	return 1, nil
}

func (r *GenericRepositorySql[TE]) All(ctx context.Context, filters QueryFilter) ([]TE, error) {
	var query = CreateReadAllQuery[TE](filters, r.DBContext.Schema)
	rows, err := r.DBContext.ExecuteQuery(query)
	if err != nil {
		log.Println(err)
	}
	result := Map[TE](*rows)

	return result, nil
}

func (r *GenericRepositorySql[TE]) ById(ctx context.Context, id int64) (*TE, error) {
	var query = CreateReadByIdQuery[TE](id, r.DBContext.Schema)
	r.DBContext.ExecuteQuery(query)
	rows, err := r.DBContext.ExecuteQuery(query)
	if err != nil {
		log.Println(err)
	}
	result := Map[TE](*rows)

	if len(result) < 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (r *GenericRepositorySql[TE]) Delete(ctx context.Context, id int64) (int64, error) {
	var query = CreateDeleteQuery[TE](id, r.DBContext.Schema)
	return r.DBContext.ExecuteDelete(query)
}

func (r *GenericRepositorySql[TE]) Query(ctx context.Context, query IDomainQuery) (any, error) {

	//TODO: singleton
	// mapping := make(map[string]IQueryHandler)
	// mapping[queries.CouponInUse{}.GetQueryName()] = queryhandlers.CouponInUseHandler{}
	handler := r.queryhandlers[query.GetQueryName()]
	return handler.Handle(query, r.DBContext), nil
}
