package class

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/labels"
	"sort"

	"context"
	"fmt"
	schoolv1 "github.com/school/school-operator/pkg/apis/school/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/school/school-operator/pkg/apis/school/v1"
	appv1 "github.com/school/school-operator/pkg/apis/school/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	k8sv1 "github.com/school/school-operator/pkg/apis/school/v1"

	"encoding/json"
)

var log = logf.Log.WithName("controller_class")

const CLASS_SIZE int32 = 10 //默认班级人数容量
// GetStatus 获取更新状态
func GetStatus(size int) string  {
	switch  {
	case size == 10:
		return "Ready"
	case size <10:
		return "NotReady"
	case size > 10:
		return "Error"
	}
	return "Unknown"
}
/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Class Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileClass{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("class-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Class
	err = c.Watch(&source.Kind{Type: &schoolv1.Class{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Class
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &schoolv1.Class{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileClass implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileClass{}

// ReconcileClass reconciles a Class object
type ReconcileClass struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Class object and makes changes based on the state read
// and what is in the Class.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileClass) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Class")

	// Fetch the Class instance
	instance := &schoolv1.Class{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	label := labels.Set{"app": instance.Name}

	existingPods := &corev1.PodList{}
	//获取name对应的所以的pod列表
	err = r.client.List(context.TODO(), existingPods, &client.ListOptions{
		Namespace:     request.Namespace,
		LabelSelector: labels.SelectorFromSet(label),
	})

	if err != nil {
		reqLogger.Error(err, "取已经存在的pod失败")
		return reconcile.Result{}, err
	}

	//取到pod列表中的pod name
	var existingPodNames []string
	for _, pod := range existingPods.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodRunning {
			existingPodNames = append(existingPodNames, pod.GetObjectMeta().GetName())
		}
	}

	reqLogger.Info("existingPodNames","size",len(existingPodNames))

	//pod数量为0，进行初始化
	if len(existingPodNames) == 0 {
		reqLogger.Info("初始化创建Pod:", "existingPodNames",
			existingPodNames,"replicas", instance.Spec.Replicas)
		//create deploy
		return r.createDeploy(instance)
	}

	////期望的status
	status := k8sv1.ClassStatus{
		PodNames: existingPodNames,
		Replicas: len(existingPodNames),
		Status: GetStatus(len(existingPodNames)),
	}
	sort.Strings(instance.Status.PodNames)
	sort.Strings(status.PodNames)
	reqLogger.Info("instance podStatus","status",instance.Status)
	//对比两个状态是否一致，如果不一致则进行更新
	if !reflect.DeepEqual(instance.Status, status) {
		reqLogger.Info("更新pod的状态","status",status)
		instance.Status = status
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "更新pod的状态失败")
			return reconcile.Result{}, err
		}
	}
	//期望值小，需要scale down
	reqLogger.Info("instance.Spec.Replicas","replicas",
		instance.Spec.Replicas,"existingPodNames replicas",len(existingPodNames))
	if len(existingPodNames) > int(*instance.Spec.Replicas) {
		return r.scaleDown(len(existingPodNames),int(*instance.Spec.Replicas),existingPods,reqLogger,instance,&request)
	}

	//len(pod)<运行中的len(pod.replace)，期望值大，需要scale up create
	if len(existingPodNames) < int(*instance.Spec.Replicas){
		return r.scaleUp(existingPodNames,instance,&request,reqLogger)
	}
	return reconcile.Result{}, nil
}

// scaleUp
func (r *ReconcileClass) scaleUp(existingPodNames []string,instance *schoolv1.Class,
	request *reconcile.Request,reqLogger logr.Logger) (reconcile.Result, error) {
	if len(existingPodNames) < int(*instance.Spec.Replicas) {
		//create
		reqLogger.Info("正在创建Pod,当前的podnames和期望的Replicas:", "existingPodNames",existingPodNames,
			"Spec.Replicas",instance.Spec.Replicas)
		deploy := &appsv1.Deployment{}
        //err != nil && errors.IsNotFound(err)
		if err := r.client.Get(context.TODO(), request.NamespacedName,deploy); err == nil  {
			//创建 Deploy
			return r.updateDeploy(instance)
		}/*else{
			data,_ := json.Marshal(deploy)
			fmt.Printf("%s\n",string(data))
		}*/
	}
	return reconcile.Result{},nil
}

//scale len(pod)>运行中的len(pod.replace)，期望值小，需要scale down
func (r *ReconcileClass)scaleDown(exitsSize int,sepcReplics int,existingPods *corev1.PodList,
	reqLogger logr.Logger,instance *schoolv1.Class,request *reconcile.Request) (reconcile.Result, error) {
	if exitsSize > sepcReplics {
		//delete
		pod := existingPods.Items[0]
		reqLogger.Info("正在删除Pod,当前的podnames和期望的Replicas:",
			"exitsPod",pod.Name, "Spec.Replicas",sepcReplics)
		err := r.client.Delete(context.TODO(), &pod)
		if err != nil {
			reqLogger.Error(err, "删除pod失败")
			return reconcile.Result{}, err
		}
		deploy := &appsv1.Deployment{}
		//err != nil && errors.IsNotFound(err)
		if err := r.client.Get(context.TODO(), request.NamespacedName,deploy); err == nil  {
			//创建 Deploy
			//*instance.Spec.Replicas = *instance.Spec.Replicas -1
			return r.updateDeploy(instance)
		}
		//r.updateDeploy()
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileClass) createDeploy(instance *schoolv1.Class) (reconcile.Result, error)  {
	deploy := NewDeploy(instance)
	if err := r.client.Create(context.TODO(), deploy); err != nil {
		return reconcile.Result{}, err
	}
	//关联 Annotations
	data, _ := json.Marshal(instance.Spec)
	if instance.Annotations != nil {
		instance.Annotations["spec"] = string(data)
	} else {
		instance.Annotations = map[string]string{"spec": string(data)}
	}
	if err := r.client.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, nil
	}
	return reconcile.Result{}, nil
}



func (r *ReconcileClass) updateDeploy(instance *schoolv1.Class) (reconcile.Result, error)  {
	deploy := NewDeploy(instance)
	if err := r.client.Update(context.TODO(), deploy); err != nil {
		return reconcile.Result{}, err
	}
	//关联 Annotations(重复代码，后面优化掉)
	data, _ := json.Marshal(instance.Spec)
	if instance.Annotations != nil {
		instance.Annotations["spec"] = string(data)
	} else {
		instance.Annotations = map[string]string{"spec": string(data)}
	}
	if err := r.client.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, nil
	}
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *schoolv1.Class,num int) *corev1.Pod {
	labels := map[string]string{
		"app": fmt.Sprintf("%s-%d",cr.Name,num),
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func NewDeploy(app *appv1.Class) *appsv1.Deployment {

	labels := map[string]string{"app": app.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group: v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind: "Class",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: newContainers(app),
				},
			},
			Selector: selector,
		},
	}
}

func newContainers(app *v1.Class) []corev1.Container {
	/*containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range app.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}*/
	return []corev1.Container{
		{
			Name: app.Name,
			Image: app.Spec.Image,
			//Resources: app.Spec.Resources,
			//Ports: containerPorts,
			ImagePullPolicy: corev1.PullIfNotPresent,
			//Env: app.Spec.Envs,
		},
	}
}