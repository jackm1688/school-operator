package models

type Classes struct {
	ApiVersion string `json:"apiVersion"`
	Items      []Item `json:"items"`
}

type Item struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       ClassesSpec            `json:"spec"`
	Status     ClassStatus            `json:"status"`
}

type ClassesSpec struct {
	ClassName   string `json:"className"`
	Image       string `json:"image"`
	Replicas    int    `json:"replicas"`
	TeacherName string `json:"teacherName"`
}

type ClassStatus struct {
	PodNames []string `json:"podNames"`
	Replicas int      `json:"replicas"`
	Status   string   `json:"status"`
}

/**
"status": {
               "podNames": [
               ],
               "replicas": 10,
               "status": "Ready"
           }
*/

/*
{
    "apiVersion": "school.crd.io/v1",
    "items": [
        {
            "apiVersion": "school.crd.io/v1",
            "kind": "Class",
            "metadata": {
                "annotations": {
                    "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"school.crd.io/v1\",\"kind\":\"Class\",\"metadata\":{\"annotations\":{},\"name\":\"student\",\"namespace\":\"default\"},\"spec\":{\"className\":\"计算机应用二班\",\"image\":\"gdsz.harbor.com/library/myhttp:v1\",\"replicas\":10,\"teacherName\":\"杰克\"}}\n",
                    "spec": "{\"className\":\"计算机应用二班\",\"teacherName\":\"杰克\",\"replicas\":10,\"image\":\"gdsz.harbor.com/library/myhttp:v1\"}"
                },
                "creationTimestamp": "2021-01-11T11:40:54Z",
                "generation": 11,
                "name": "student",
                "namespace": "default",
                "resourceVersion": "12225068",
                "selfLink": "/apis/school.crd.io/v1/namespaces/default/classes/student",
                "uid": "efa79455-c8be-4894-a52e-034dac89c77c"
            },
            "spec": {
                "className": "计算机应用二班",
                "image": "gdsz.harbor.com/library/myhttp:v1",
                "replicas": 10,
                "teacherName": "杰克"
            },
            "status": {
                "podNames": [
                ],
                "replicas": 10,
                "status": "Ready"
            }
        }
    ],
    "kind": "ClassList",
    "metadata": {
        "continue": "",
        "resourceVersion": "12229422",
        "selfLink": "/apis/school.crd.io/v1/classes"
    }
}
*/
