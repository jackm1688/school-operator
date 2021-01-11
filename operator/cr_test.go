package test

import (
	"testing"
	"time"

	"github.com/school/school-operator/models"
	"github.com/school/school-operator/operator"
)

func GetCRFile() *operator.CR {
	st := models.Student{
		ApiVersion: "school.crd.io/v1",
		Kind:       "Class",
		Metadata: models.Metadata{
			Name: "student",
		},
		Spec: models.Spec{
			ClassName:   "计算机应用二班",
			TeacherName: "杰克",
			Replicas:    11,
			Image:       "gdsz.harbor.com/library/myhttp:v1",
		},
	}

	cr := operator.CR{
		YamlFileName: "student_cr.yaml",
		Stu:          &st,
	}
	return &cr
}

func TestGenYAMLFile(t *testing.T) {

	cr := GetCRFile()

	err := cr.GenYAMLFile()
	if err != nil {
		t.Error(err)
	}
}

// create cr resource
func TestCreateCR(t *testing.T) {
	cr := GetCRFile()
	err := cr.CreateCR()
	if err != nil {
		t.Error(err)
	}
}

// update cr resource
func TestUpdateCR(t *testing.T) {
	cr := GetCRFile()
	err := cr.UpdateCR()
	if err != nil {
		t.Error(err)
	}
}

// delete resource
func TestDeleteCR(t *testing.T) {
	cr := GetCRFile()
	err := cr.DeleteCR()
	if err != nil {
		t.Error(err)
	}
}

// get cr resource
func TestGetCR(t *testing.T) {
	cr := GetCRFile()
	err := cr.GetCR()
	if err != nil {
		t.Error(err)
	}
}

// sclae cr resource
func TestScaleCR(t *testing.T) {
	cr := GetCRFile()
	err := cr.ScaleCR(10)
	if err != nil {
		t.Error(err)
	}
}

func TestDescribe(t *testing.T) {
	cr := GetCRFile()
	err := cr.Describe()
	if err != nil {
		t.Error(err)
	}
}

func TestCustomeScaleAndCheckResult(t *testing.T) {
	cr := GetCRFile()
	var replics = 10
	err := cr.CustomeScaleAndCheckResult(replics)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(10 * time.Second)
	_, status := cr.GetStatus()
	t.Logf("status:%#v\n", status)
	if status.Replicas == replics {
		t.Logf("scale update sucess,status:%s,replics:%d, scaleReplics:%d\n", status.Status, replics, status.Replicas)
	} else {
		t.Logf("scale update failed,status:%s,replics:%d, scaleReplics:%d\n", status.Status, replics, status.Replicas)
	}

}
