package universe

import (
// "fmt"
)

type Region struct {
	Name          string
	Constelations []*Constelation
	Systems       []*System
}

type Constelation struct {
	Name    string
	Systems []*System
	Region  *Region
}

type System struct {
	Name         string
	Region       *Region
	Constelation *Constelation
}

var (
	Systems       map[string]*System       = make(map[string]*System)
	Constelations map[string]*Constelation = make(map[string]*Constelation)
	Regions       map[string]*Region       = make(map[string]*Region)
)
