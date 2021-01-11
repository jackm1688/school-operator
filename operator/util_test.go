package operator

import (
	"testing"
)

func TestGetClassStaus(t *testing.T) {
	err, res := GetClassStatus()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", res.Items[0].Status)
}
