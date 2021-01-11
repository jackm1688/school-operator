package models

type A struct {
}

type Student struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
}

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type Spec struct {
	ClassName   string `json:"className" yaml:"className"`
	TeacherName string `json:"teacherName" yaml:"teacherName"`
	Replicas    int    `json:"replicas" yaml:"replicas"`
	Image       string `json:"image" yaml:"image"`
}
