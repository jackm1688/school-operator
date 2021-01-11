package operator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildImageAndPush(t *testing.T) {

	repo := "gdsz.harbor.com/library"
	appName := "school-operator"
	version := "v2"
	app := App{
		Repo:    repo,
		AppName: appName,
		Version: version,
	}
	err, _ := app.BuildImageAndPush()
	if err != nil {
		t.Error(err)
	}

	assert.CallerInfo()
}
