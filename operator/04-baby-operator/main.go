package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var gvr = schema.GroupVersionResource{
	Group:    "example.com",
	Version:  "v1alpha1",
	Resource: "foos",
}

// カスタムリソースFooの定義
type Foo struct {
	// inlineやomitemptyはJSONのマーシャリングに関する設定: https://qiita.com/nakampany/items/94c58340f81234970250
	metav1.TypeMeta   `json:",inline"`            // Kubernetes API のリソースタイプを保持
	metav1.ObjectMeta `json:"metadata,omitempty"` // Kubernetes リソースのメタデータを保持

	TestString string `json:"testString"` // カスタムフィールド
	TestNum    int    `json:"testNum"`
}

// 複数の Foo リソースをリスト形式で表現する構造体
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}

func listFoos(client dynamic.Interface, namespace string) (*FooList, error) {
	// gvrとnamespaceを指定してリソースを取得
	list, err := client.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// リソースをJSONに変換
	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var foolist FooList // FooList型
	// 取得したJSONをFooList型に変換
	if err := json.Unmarshal(data, &foolist); err != nil {
		return nil, err
	}
	return &foolist, nil
}

func createPod(clientset *kubernetes.Clientset, namespace, name string) error {
	// Podという構造体の初期化
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx:latest",
				},
			},
			RestartPolicy: v1.RestartPolicyAlways,
		},
	}
	// namespaceを指定、contextが必要なのでcontext.TODO()を指定
	_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("failed to create pod %v\n", err)
		return err
	}
	fmt.Println("successfully created pod")
	return nil
}

func main() {
	var defaultKubeconfigPath string
	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfigPath = filepath.Join(home, ".kube", "config")
	}
	kubeconfig := flag.String("kubeconfig", defaultKubeconfigPath, "kubeconfig config file")
	flag.Parse()
	// https://pkg.go.dev/k8s.io/client-go/tools/clientcmd#BuildConfigFromFlags
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// podではなくカスタムリソースを扱うためdynamic clientを使う
	// https://pkg.go.dev/k8s.io/client-go/dynamic#NewForConfig
	client, _ := dynamic.NewForConfig(config)

	clientset, _ := kubernetes.NewForConfig(config) // 既存のpodを確認するため

	// 無限ループ
	for {
		foos, _ := listFoos(client, "")
		fmt.Println("INDEX\tNAMESPACE\tNAME")
		for i, foo := range foos.Items {
			namespace := foo.GetNamespace()
			name := foo.GetName()
			fmt.Printf("%d\t%s\t%s\n", i, namespace, name)

			// 指定されたpodが存在するか確認
			_, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					fmt.Println("Pod doesn't exist. Creating new Pod")
					createPod(clientset, namespace, name)
				} else {
					fmt.Printf("failed to get pod %v\n", err)
				}
			} else {
				fmt.Printf("successfully got pod %s\n", name)
			}
		}
		// 1秒ごとにリソースを取得し、podの状態を更新
		time.Sleep(1 * time.Second)
	}

}
