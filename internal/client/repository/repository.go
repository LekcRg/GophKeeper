package repository

type Repository interface {
	GetTable() [][]string
}
