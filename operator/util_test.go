package operator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGetClassStaus(t *testing.T) {
	err, res := GetClassStatus()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", res.Items[0].Status)
}

func GetClient() dynamic.Interface {

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}

func TestCreateStudentWithYaml(t *testing.T) {
	stu := `
apiVersion: school.crd.io/v1
kind: Class
metadata:
  name: student
spec:
  className: "计算机应用二班"
  teacherName: "杰克"
  replicas: 8
  image: gdsz.harbor.com/library/myhttp:v1`

	st, err := CreateStudentWithYaml(GetClient(), "default", stu)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", st)
}

func TestUpdateStudentWithYaml(t *testing.T) {
	stu := `
apiVersion: school.crd.io/v1
kind: Class
metadata:
  name: student
spec:
  className: "计算机应用二班"
  teacherName: "杰克"
  replicas: 30
  image: gdsz.harbor.com/library/myhttp:v1`

	st, err := UpdateStudentWithYaml(GetClient(), "default", stu)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", st)
}

func TestDeleteStudent(t *testing.T) {
	err := DeleteStudent(GetClient(), "default", "student")
	if err != nil {
		t.Error(err)
	}
}

/**
状态更新可以定义三个状态：Ready/NotReady/Error，代码中可以定义一个默认班级的容量：
  比如10人，如果ClassSize小于10，状态为NotReady，
              ClassSize等于10，状态为Ready，
              ClassSize大于10状态为Error。
这里可以使用assert返回结果
*/
func TestReadyAndNotReadyAndErrStatus(t *testing.T) {

	_, err := UpdateStudentWithYaml(GetClient(), "default", genStuYAML(5))
	if err != nil {
		t.Fatal(err)
	} else {
		time.Sleep(30 * time.Second)
		err, status := GetClassStatus()
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("【ClassSize小于10，状态为NotReady】replicas:%d,status:%s\n", status.Items[0].Status.Replicas, status.Items[0].Status.Status)
		}
	}
	fmt.Printf("\n------------------------------------------------------------------------------------------\n")
	_, err = UpdateStudentWithYaml(GetClient(), "default", genStuYAML(10))
	if err != nil {
		t.Fatal(err)
	} else {
		time.Sleep(20 * time.Second)
		err, status := GetClassStatus()
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("【ClassSize等于10，状态为Ready】replicas:%d,status:%s\n", status.Items[0].Status.Replicas, status.Items[0].Status.Status)
		}
	}
	fmt.Printf("\n------------------------------------------------------------------------------------------\n")

	_, err = UpdateStudentWithYaml(GetClient(), "default", genStuYAML(15))
	if err != nil {
		t.Fatal(err)
	} else {
		time.Sleep(30 * time.Second)
		err, status := GetClassStatus()
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("【ClassSize大于10状态为Error】replicas:%d,status:%s\n", status.Items[0].Status.Replicas, status.Items[0].Status.Status)
		}
	}
	fmt.Printf("\n------------------------------------------------------------------------------------------\n")

	//t.Logf("%#v\n", st.Spec)
}

func genStuYAML(replicas int) string {

	return fmt.Sprintf(`
apiVersion: school.crd.io/v1
kind: Class
metadata:
  name: student
spec:
  className: "计算机应用二班"
  teacherName: "杰克"
  replicas: %d
  image: gdsz.harbor.com/library/myhttp:v1`, replicas)

}
