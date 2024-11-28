package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var defaultKubeconfigPath string
	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfigPath = filepath.Join(home, ".kube", "config")
	}
	// コマンドラインで--kubeconfigを使えるようにする(デフォルトはdefaultKubeconfigPath)
	kubeconfig := flag.String("kubeconfig", defaultKubeconfigPath, "kubeconfig config file")
	flag.Parse()

	// _はエラー無視 指定したkubeconfigを用いてクライアントを作成
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// クラスタにアクセスするためのクライアントセットを作成
	// https://pkg.go.dev/k8s.io/client-go/kubernetes#NewForConfig
	clientset, _ := kubernetes.NewForConfig(config)

	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	// print
	fmt.Println("NAMESPACE\tNAME")
	for _, pod := range pods.Items {
		fmt.Printf("%s\t%s\n", pod.GetNamespace(), pod.GetName())
	}

}
