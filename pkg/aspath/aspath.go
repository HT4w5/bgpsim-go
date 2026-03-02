package aspath

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strings"
)

type AsPath struct {
	idxMap map[uint32]int
	path   []uint32
}

func New() *AsPath {
	return &AsPath{
		idxMap: map[uint32]int{},
		path:   []uint32{},
	}
}

// Add a new AS number to the front of the path
func (ap *AsPath) Prepend(as uint32) error {
	if _, ok := ap.idxMap[as]; ok {
		return fmt.Errorf("duplicate AS number in path: %d", as)
	}
	ap.idxMap[as] = len(ap.path)
	ap.path = append(ap.path, as)
	return nil
}

func (ap *AsPath) HasAs(as uint32) bool {
	_, ok := ap.idxMap[as]
	return ok
}

func (ap *AsPath) Len() int {
	return len(ap.path)
}

func (ap *AsPath) GetPath() []uint32 {
	return ap.path
}

func (ap *AsPath) String() string {
	if len(ap.path) == 0 {
		return ""
	}
	var sb strings.Builder
	for i := len(ap.path) - 1; i >= 1; i-- {
		sb.WriteString(fmt.Sprintf("%d ", ap.path[i]))
	}
	sb.WriteString(fmt.Sprintf("%d", ap.path[0]))
	return sb.String()
}

func (ap *AsPath) Hash() uint32 {
	h := fnv.New32a()
	binary.Write(h, binary.BigEndian, ap.path)
	return h.Sum32()
}
