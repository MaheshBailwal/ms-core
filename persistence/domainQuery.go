package persistence

type IDomainQuery interface {
	GetQueryName() string
}
