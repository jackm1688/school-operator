package operator

import (
	"io/ioutil"

	"github.com/school/school-operator/models"

	"gopkg.in/yaml.v2"
)

//对cr资源的操作
type CR struct {
	YamlFileName string
	Stu          *models.Student
}

// 生成stu yaml文件
func (cr *CR) GenYAMLFile() error {
	out, err := yaml.Marshal(cr.Stu)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cr.YamlFileName, out, 0755)
}

// create cr resource
func (cr *CR) CreateCR() error {
	return execCmd(kubelet, "create", "-f", cr.YamlFileName)
}

// update cr resource
func (cr *CR) UpdateCR() error {
	_ = cr.GenYAMLFile()
	return execCmd(kubelet, "apply", "-f", cr.YamlFileName)
}

// delete resource
func (cr *CR) DeleteCR() error {
	return execCmd(kubelet, "delete", "-f", cr.YamlFileName)
}

// get cr resource
func (cr *CR) GetCR() error {
	return execCmd(kubelet, "get", "-f", cr.YamlFileName)
}

// sclae cr resource
func (cr *CR) ScaleCR(replicas int) error {
	cr.Stu.Spec.Replicas = replicas
	_ = cr.GenYAMLFile()
	return execCmd(kubelet, "apply", "-f", cr.YamlFileName)
}

//kubectl describe class student
func (cr *CR) Describe() error {
	return execCmd(kubelet, "describe", cr.Stu.Kind, cr.Stu.Metadata.Name)
}

//可以使用go-client客户来搞哈，目前先用命令+api接口来实现
func (cr *CR) CustomeScaleAndCheckResult(replicas int) error {
	cr.Stu.Spec.Replicas = replicas
	_ = cr.GenYAMLFile()
	err := execCmd(kubelet, "apply", "-f", cr.YamlFileName)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CR) GetStatus() (error, *models.ClassStatus) {
	err, res := GetClassStatus()
	if err != nil {
		return err, nil
	}
	return nil, &res.Items[0].Status
}
