package main

import (
	"fmt"
	"log"
)

type valset struct {
	val []interface{}
}

// VarTracker tracks variable values
type VarTracker struct {
	Names  []string
	Values []valset
}

// NewVarTracker returns a new VarTracker object
func NewVarTracker() *VarTracker {
	return &VarTracker{}
}

// SetNames sets name of variables to be tracked
func (v *VarTracker) SetNames(names ...string) {
	for i := range names {
		v.Names = append(v.Names, names[i])
		v.Values = append(v.Values, valset{})
	}
}

// Store a set of values
func (v *VarTracker) Store(newvals ...interface{}) {
	if len(newvals) != len(v.Values) {
		log.Println("Wrong number of variables", len(newvals), len(v.Values))
		return
	}
	for i := range newvals {
		v.Values[i].val = append(v.Values[i].val, newvals[i])
	}
}

// Reset the stored values
func (v *VarTracker) Reset() {
	for i := range v.Values {
		v.Values[i] = valset{}
	}
}

// Dump out the values being tracked
func (v *VarTracker) Dump() {
	for i := range v.Names {
		fmt.Print(v.Names[i], ",\t")
	}
	fmt.Println()

	for i := range v.Values[0].val {
		for j := range v.Values {
			fmt.Print(v.Values[j].val[i], ",\t")
		}
		fmt.Println()
	}
}
