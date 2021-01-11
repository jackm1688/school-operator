package operator

import (
	"testing"
)

func TestReadImage(t *testing.T) {
	err, r := ReadImage()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("=====> result:%s\n", r)
	}
}
