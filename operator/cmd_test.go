package test

import (
	"testing"

	"github.com/school/school-operator/operator"
)

func TestReadImage(t *testing.T) {
	err, r := operator.ReadImage()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("=====> result:%s\n", r)
	}
}
