package operator

import (
	"context"
	//"fmt"
	"net/http"

	//v1 "k8s.io/api/apps/v1"

	"github.com/school/school-operator/models"

	"encoding/json"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	//"k8s.io/client-go/tools/clientcmd"
)

var ClassUrl = "http://192.168.0.80:8080/apis/school.crd.io/v1/classes"

var gvr = schema.GroupVersionResource{
	Group:    "school.crd.io",
	Version:  "v1",
	Resource: "classes",
}

func GetClassStatus() (error, *models.Classes) {

	request, err := http.NewRequest(http.MethodGet, ClassUrl, nil)
	if err != nil {
		return err, nil
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	res := models.Classes{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err, nil
	}
	return nil, &res
}

func CreateStudentWithYaml(client dynamic.Interface, namespace string, yamlData string) (*models.Student, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode([]byte(yamlData), nil, obj); err != nil {
		return nil, err
	}

	utd, err := client.Resource(gvr).Namespace(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var st models.Student
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func UpdateStudentWithYaml(client dynamic.Interface, namespace string, yamlData string) (*models.Student, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode([]byte(yamlData), nil, obj); err != nil {
		return nil, err
	}

	utd, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	obj.SetResourceVersion(utd.GetResourceVersion())
	utd, err = client.Resource(gvr).Namespace(namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var st models.Student
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func DeleteStudent(client dynamic.Interface, namespace string, name string) error {
	return client.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

}
