package operator_test

import (
	"fmt"
	"os/exec"
)

/**
operator-sdk build gdsz.harbor.com/library/school-operator:v1
docker  push  gdsz.harbor.com/library/school-operator:v1
 */
var (
	opSdk = "operator-sdk"
)

// BuildImageAndPush 镜像构建推送到镜像仓库
func BuildImageAndPush(appName string,version string)  string {
	imageName := fmt.Sprintf("%s:%s",appName,version)
	exec.Command(opSdk,"build",imageName)
	
}
