package model

import (
	"fmt"
	"strings"
)

type AsPath struct {
	indexMap map[int]int
	path     []int
}

func NewAsPath() *AsPath {
	return &AsPath{
		indexMap: map[int]int{},
		path:     []int{},
	}
}

// Add a new AS number to the front of the path
func (ap *AsPath) prepend(as int) error {
	if _, ok := ap.indexMap[as]; ok {
		return fmt.Errorf("duplicate AS number in path: %d", as)
	}
	ap.indexMap[as] = len(ap.path)
	ap.path = append(ap.path, as)
	return nil
}

func (ap *AsPath) HasAs(as int) bool {
	_, ok := ap.indexMap[as]
	return ok
}

func (ap *AsPath) Len() int {
	return len(ap.path)
}

func (ap *AsPath) String() string {
	var sb strings.Builder
	for i := len(ap.path) - 1; i >= 1; i-- {
		sb.WriteString(fmt.Sprintf("%d ", ap.path[i]))
	}
	sb.WriteString(fmt.Sprintf("%d", ap.path[0]))
	return sb.String()
}
