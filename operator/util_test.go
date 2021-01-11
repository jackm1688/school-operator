package test

import (
	"testing"

	"github.com/school/school-operator/operator"
)

func TestGetClassStaus(t *testing.T) {
	err, res := operator.GetClassStatus()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", res.Items[0].Status)
}
