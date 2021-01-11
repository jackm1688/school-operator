package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ClassPodSpec struct {
	Replicas int `json:"replicas"`
}

type ClassPodStatus struct {
	Replicas int `json:"replicas"`
	PodNames []string `json:"podNames"`
}

// ClassSpec defines the desired state of Class
type ClassSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClassName   string `json:"className"`
	TeacherName string `json:"teacherName"` //班主任姓名
	Replicas    *int32 `json:"replicas"`    //学生人数
	Image       string `json:"image"`
}

// ClassStatus defines the observed state of Class
type ClassStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	appsv1.DeploymentStatus `json:",inline"`
	Replicas int `json:"replicas"`
	PodNames []string `json:"podNames,omitempty"`
	Status string `json:"status"`
}

type ClassScale struct {
	Replicas   *int32 `json:"replicas"`   //学生人数
}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Class is the Schema for the classes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=classes,scope=Namespaced
type Class struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec   ClassSpec   `json:"spec,omitempty"`
	Status ClassStatus `json:"status,omitempty"`
}

type ClassAdditionalPrinterColumns struct {
	ClassName   string `json:"className"`
	TeacherName string `json:"teacherName"` //班主任姓名
	Replicas    *int32 `json:"replicas"`    //学生人数
}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClassList contains a list of Class
type ClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items    []Class `json:"items"`
}


func init() {
	SchemeBuilder.Register(&Class{}, &ClassList{})
}
