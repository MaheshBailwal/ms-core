package persistence

type IQueryHandler interface {
	HandlesQuery() string
	Handle(IDomainQuery, *DBContext) any
}
