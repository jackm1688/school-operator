package controller

import (
	"github.com/school/school-operator/pkg/controller/class"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, class.Add)
}
