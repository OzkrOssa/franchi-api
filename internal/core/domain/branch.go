package domain

type Branch struct {
	ID          uint64
	Name        string
	FranchiseID uint64
	Product     []Product
}
