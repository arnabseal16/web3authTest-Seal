package main

import (
	"fmt"
	"context"
	"flag"
	"path/filepath"
	"encoding/json"

	types "k8s.io/apimachinery/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

var (

	deploymentName = "test-hpa-dep"
	memscaleFactor = uint32(20)
	cpuscaleFactor = uint32(2)		
)

type patchMapValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value interface{} `json:"value"`
}

//  patchUint32Value specifies a patch operation for a uint32.
type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

func scaleReplicationController(clientset *kubernetes.Clientset, deploymentName string) error {


	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err))
		}

		payload := []patchMapValue{{
			Op:    "replace",
			Path:  "/spec/containers[0]/resources/cpu",
			Value: "20m",
		}}
		payloadBytes, _ := json.Marshal(payload)

		for _, v := range pods.Items {
			fmt.Println(v.Name)
			newV1Pod, err :=clientset.CoreV1().Pods("default").Patch(context.TODO(), v.Name, types.StrategicMergePatchType, payloadBytes, metav1.PatchOptions{}, "")
			fmt.Println(newV1Pod)
			return err
		}
		return err
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	//fmt.Println("Updated deployment...")

	return nil
}

func int32Ptr(i int32) *int32 { return &i }

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	err = scaleReplicationController(clientset, deploymentName)
	if err != nil {
		panic(err.Error())
	}
}
