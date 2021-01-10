operator-sdk命令介绍:<br>
operator-sdk add [api|controller]命令将控制器或资源添加到项目中。该命令必须从Operator项目的根目录运行。 <br>
operator-sdk build 命令编译代码并生成可执行文件<br>
operator-sdk completion 生成bash补全<br>
operator-sdk print-deps 命令显示操作员所需的最新Golang软件包和版本。默认情况下，它以列格式打印<br>
operator-sdk generate 命令将调用特定的生成器以根据需要生成代码(k8s,crds)<br>
operator-sdk olm-catalog gen-csv 子命令将集群服务版本（CSV）清单和可选的自定义资源定义（CRD）文件写入 deploy/olm-catalog/<operator_name>/<csv_version> <br>


=================================================================================================<br>
<br>
工作流程
SDK提供以下工作流程来开发新的Operator：<br>
  1.使用SDK命令行界面（CLI）创建新的Operator项目<br>
  2.通过添加自定义资源定义（CRD）定义新资源API<br>
  3.使用SDK API监控指定的资源<br>
  4.在指定的处理程序中定义Operator协调逻辑(对比期望状态与实际状态)，并使用SDK API与资源进行交互<br>
  5.使用SDK CLI构建并生成Operator部署manifests<br>
Operator使用SDK在用户自定义的处理程序中以高级API处理监视资源的事件，并采取措施来reconcile（对比期望状态与实际状态）应用程序的状态。<br><br>

<br>
=================================================================================================<br>
1.创建一个学校仓库 <br>
operator-sdk new school-operator --repo=github.com/school/school-operator<br>
cd school-operator<br>

2.Add a new Custom Resource Definition<br>
operator-sdk  add api --api-version=school.crd.io/v1 --kind=Class<br>

3.Define the spec and status<br>

// ClassSpec defines the desired state of Class<br>
type ClassSpec struct {<br>
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster<br>
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file<br>
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html<br>
	ClassName   string `json:"className"`   //班级名称<br>
	TeacherName string `json:"teacherName"` //班主任姓名<br>
	ClassSize   *int32 `json:"classSize"`   //学生人数<br>
	Image       string `json:"image"`<br>
}

// ClassStatus defines the observed state of Class<br>
type ClassStatus struct {<br>
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster<br>
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file<br>
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html<br>
	appsv1.DeploymentStatus `json:",inline"`<br>
}

operator-sdk generate k8s<br>
<br>
4.Updating CRD manifests<br>
operator-sdk generate crds<br>

5.Add a new Controller<br>
operator-sdk add controller --api-version=school.crd.io/v1 --kind=Class<br>


//controller<br>

func NewDeploy(app *appv1.Class) *appsv1.Deployment {<br>
	labels := map[string]string{"app": app.Name}<br>
	selector := &metav1.LabelSelector{MatchLabels: labels}<br>
	return &appsv1.Deployment{<br>
		TypeMeta: metav1.TypeMeta{<br>
			APIVersion: "apps/v1",<br>
			Kind:       "Deployment",<br>
		},
		ObjectMeta: metav1.ObjectMeta{<br>
			Name:      app.Name,<br>
			Namespace: app.Namespace,<br>

			OwnerReferences: []metav1.OwnerReference{<br>
				*metav1.NewControllerRef(app, schema.GroupVersionKind{<br>
					Group: v1.SchemeGroupVersion.Group,<br>
					Version: v1.SchemeGroupVersion.Version,<br>
					Kind: "AppService",<br>
				}),<br>
			},<br>
		},<br>
		Spec: appsv1.DeploymentSpec{<br>
			Replicas: app.Spec.ClassSize,<br>
			Template: corev1.PodTemplateSpec{<br>
				ObjectMeta: metav1.ObjectMeta{<br>
					Labels: labels,<br>
				},
				Spec: corev1.PodSpec{<br>
					Containers: newContainers(app),<br>
				},<br>
			},<br>

			Selector: selector,<br>
		},<br>
	}<br>
}<br>

func newContainers(app *v1.Class) []corev1.Container {<br>
	/*containerPorts := []corev1.ContainerPort{}<br>
	for _, svcPort := range app.Spec.Ports {<br>
		cport := corev1.ContainerPort{}<br>
		cport.ContainerPort = svcPort.TargetPort.IntVal<br>
		containerPorts = append(containerPorts, cport)<br>
	}*/<br>
	return []corev1.Container{<br>
		{<br>
			Name: app.Name,<br>
			Image: app.Spec.Image,<br>
			//Resources: app.Spec.Resources,<br>
			//Ports: containerPorts,<br>
			ImagePullPolicy: corev1.PullIfNotPresent,<br>
			//Env: app.Spec.Envs,<br>
		},<br>
	}<br>
}<br>
<br>

6.Build and run the operator<br>
operator-sdk build gdsz.harbor.com/library/school-operator:v1<br>
docker push gdsz.harbor.com/library/school-operator:v1<br>


kubectl create -f deploy/crds/classes.crd.io_classes_crd.yaml<br>
customresourcedefinition.apiextensions.k8s.io/classes.classes.crd.io created<br>

//本地(监听)测试<br>
operator-sdk run local --watch-namespace=default<br><br>
<br>

7.Run as a Deployment inside the cluster<br>
kubectl create -f deploy/service_account.yaml<br>
kubectl create -f deploy/role.yaml<br>
kubectl create -f deploy/role_binding.yaml<br>
kubectl create -f deploy/operator.yaml<br>

<br>
kubectl delete -f deploy/service_account.yaml<br>
kubectl delete -f deploy/role.yaml<br>
kubectl delete -f deploy/role_binding.yaml<br>
kubectl delete -f deploy/operator.yaml<br>



kubectl get deployment school-operator<br>
NAME              READY   UP-TO-DATE   AVAILABLE   AGE<br>
school-operator   0/1     1            0           2m18s<br>

kubectl get deployment school-operator<br>
NAME              READY   UP-TO-DATE   AVAILABLE   AGE<br>
school-operator   1/1     1            1           2m49s<br>


kubectl get  pod school-operator-7c7bb89f54-x48j8<br>
NAME                               READY   STATUS    RESTARTS   AGE<br>
school-operator-7c7bb89f54-x48j8   1/1     Running   1          3m7s<br>

<br>
8.Create a Memcached CR<br>

apiVersion: school.crd.io/v1<br>
kind: Class<br>
metadata:<br>
  name: student<br>
spec:<br>
  # Add fields here<br>
  className: "计算机应用二班"<br>
  teacherName: "杰克"<br>
  replicas: 10<br>
  image: gdsz.harbor.com/library/myhttp:v1<br>
<br>
kubectl apply -f  deploy/crds/classes.crd.io_v1_class_cr.yaml<br>
class.classes.crd.io/lemon-class created<br>

# 查看应用<br>
kubectl get deployment student<br>
NAME      READY   UP-TO-DATE   AVAILABLE   AGE<br>
student   10/10   10           10          15m<br>

<br>
kubectl get pod  -l app=student  -L app<br>
NAME                       READY   STATUS    RESTARTS   AGE   APP<br>
student-6b7597f46b-48xsv   1/1     Running   0          19m   student<br>
student-6b7597f46b-7qkbj   1/1     Running   0          19m   student<br>
student-6b7597f46b-fccbb   1/1     Running   0          19m   student<br>
student-6b7597f46b-hksgn   1/1     Running   0          19m   student<br>
student-6b7597f46b-v9pdc   1/1     Running   0          19m   student<br>
student-6b7597f46b-vjt89   1/1     Running   0          19m   student<br>
student-6b7597f46b-w2qkv   1/1     Running   0          19m   student<br>
student-6b7597f46b-w7h4h   1/1     Running   0          19m   student<br>
student-6b7597f46b-x5sgz   1/1     Running   0          19m   student<br>
student-6b7597f46b-xkwt2   1/1     Running   0          19m   student<br>


<br>

kubectl get crd classes.school.crd.io  <br>
NAME                    CREATED AT <br>
classes.school.crd.io   2021-01-10T17:25:52Z <br>
ip-192-168-43-34:~ JackMeng$ kubectl describe crd classes.school.crd.io <br>

<br>

kubectl get class <br>
NAME      AGE <br>
student   21m <br>

kubectl describe  class <br>
Name:         student <br>
Namespace:    default <br>
Labels:       <none> <br>
Annotations:  kubectl.kubernetes.io/last-applied-configuration: <br>
                {"apiVersion":"school.crd.io/v1","kind":"Class","metadata":{"annotations":{},"name":"student","namespace":"default"},"spec":{"className":"...<br>
              spec: {"className":"计算机应用二班","teacherName":"杰克","replicas":10,"image":"gdsz.harbor.com/library/myhttp:v1"} <br>
API Version:  school.crd.io/v1<br>
Kind:         Class<br>
Metadata:<br>
  Creation Timestamp:  2021-01-10T17:51:54Z<br>
  Generation:          1<br>
  Resource Version:    12002040<br>
  Self Link:           /apis/school.crd.io/v1/namespaces/default/classes/student<br>
  UID:                 7a31decb-0a96-4543-bd32-4b706667d9e3<br>
Spec:<br>
  Class Name:    计算机应用二班<br>
  Image:         gdsz.harbor.com/library/myhttp:v1<br>
  Replicas:      10<br>
  Teacher Name:  杰克<br>
Status:<br>
  Pod Names:<br>
    student-6b7597f46b-x5sgz<br>
    student-6b7597f46b-fccbb<br>
    student-6b7597f46b-48xsv<br>
    student-6b7597f46b-vjt89<br>
    student-6b7597f46b-xkwt2<br>
    student-6b7597f46b-7qkbj<br>
    student-6b7597f46b-w2qkv<br>
    student-6b7597f46b-v9pdc<br>
    student-6b7597f46b-w7h4h<br>
    student-6b7597f46b-hksgn<br>
  Replicas:  10<br>
  Status:    Ready<br>
Events:      <none><br>

<br>
# 未解决的问题:<br>
  状态管理(Ready OR Not Ready控制)，还需要多研究源码下及组件的工作原理。<br>
  显示班级、老师及学生数(v0.18.0没有查阅相关资料)<br>
