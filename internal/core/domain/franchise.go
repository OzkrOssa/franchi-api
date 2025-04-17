package domain

type Franchise struct {
	ID     uint64
	Name   string
	Branch []Branch
}
