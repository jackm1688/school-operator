package operator

import (
	"fmt"
)

//operator-sdk控制器的操作对象

const (
	F1 = "deploy/service_account.yaml"
	F2 = "deploy/role.yaml"
	F3 = "deploy/role_binding.yaml"
	F4 = "deploy/operator.yaml"
	F5 = "deploy/crds/school.crd.io_classes_crd.yaml"
)

type App struct {
	Repo    string
	AppName string
	Version string
}

// BuildImageAndPush 镜像构建推送到镜像仓库
func (app *App) BuildImageAndPush() (error, string) {
	imageName := fmt.Sprintf("%s/%s:%s", app.Repo, app.AppName, app.Version)
	err := execCmd(opSdk, "build", imageName)
	if err != nil {
		return err, imageName
	}
	err = execCmd(docker, "push", imageName)
	if err != nil {
		return err, imageName
	}
	return nil, imageName
}

func (app *App) UpdateImge() (error, string) {
	imageName := fmt.Sprintf("%s/%s:%s", app.Repo, app.AppName, app.Version)
	_, srcIm := ReadImage()
	if srcIm == "" {
		return fmt.Errorf("get image failed"), imageName
	}
	return replaceImage(srcIm, imageName), imageName
}

//部署operator controller
func (app *App) InstallApp() (err error) {

	if err = execCmd(kubelet, "apply", "-f", F1); err != nil {
		return err
	}

	if err = execCmd(kubelet, "apply", "-f", F2); err != nil {
		return err
	}

	if err = execCmd(kubelet, "apply", "-f", F3); err != nil {
		return err
	}

	if err = execCmd(kubelet, "apply", "-f", F5); err != nil {
		return err
	}

	if err = execCmd(kubelet, "apply", "-f", F4); err != nil {
		return err
	}
	return nil
}

//卸载operator controller
func (app *App) UninstallApp() (err error) {

	if err = execCmd(kubelet, "delete", "-f", F1); err != nil {
		return err
	}

	if err = execCmd(kubelet, "delete", "-f", F2); err != nil {
		return err
	}

	if err = execCmd(kubelet, "delete", "-f", F3); err != nil {
		return err
	}

	if err = execCmd(kubelet, "delete", "-f", F5); err != nil {
		return err
	}

	if err = execCmd(kubelet, "delete", "-f", F4); err != nil {
		return err
	}
	return nil
}

//获取operator controller 状态
func (app *App) CheckStatusForApp() (err error) {

	if err = execCmd(kubelet, "get", "-f", F1); err != nil {
		return err
	}

	if err = execCmd(kubelet, "get", "-f", F2); err != nil {
		return err
	}

	if err = execCmd(kubelet, "get", "-f", F3); err != nil {
		return err
	}

	if err = execCmd(kubelet, "get", "-f", F4); err != nil {
		return err
	}
	return nil
}

//更新operator controller
func (app *App) UpdateFormOperatorDeployApp(dstImage string) error {
	_, srcImage := ReadImage()
	var err error
	if srcImage != "" && dstImage != "" {
		err = replaceImage(srcImage, dstImage)
	}
	if err != nil {
		return err
	}

	if err = execCmd(kubelet, "delete", "-f", F4); err != nil {
		return err
	}

	if err = execCmd(kubelet, "apply", "-f", F4); err != nil {
		return err
	}
	return nil
}
