package pathfinding

import (
	"anvil/internal/grid"
)

type Result struct {
	path  []grid.Position
	cost  int
	valid bool
}

func NewResult(path []grid.Position, cost int, valid bool) *Result {
	return &Result{
		path:  path,
		cost:  cost,
		valid: valid,
	}
}

func (r *Result) Path() []grid.Position {
	return r.path
}

func (r *Result) Cost() int {
	return r.cost
}

func (r *Result) IsValid() bool {
	return r.valid
}
