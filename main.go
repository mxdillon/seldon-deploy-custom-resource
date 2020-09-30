package main

import (
    "context"
    "flag"
    "fmt"
    appsv1 "k8s.io/api/apps/v1"
    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "log"
)


var k string

// init allows user to specify location of config file. Defaults to working directory.
func init() {
    flag.StringVar(&k, "kconfig", "./config", "Location of the k8s config file")
    flag.Parse()
}

func main() {

    // uses the current context in kubeconfig
    // path-to-kubeconfig -- for example, /root/.kube/config
    config, err := clientcmd.BuildConfigFromFlags("", k)
    if err != nil{
        log.Panic(err)
    }

    // creates the clientset - HANDLE THE ERROR!
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Panic(err)
    }

    // return clientset.CoreV1()

    api := clientset.CoreV1()

    // access the API to list pods
    pods, _ := api.Pods("").List(context.TODO(), metav1.ListOptions{})
    fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

    // apply CRD to specified namespace

    deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: "demo-deployment",
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: int32Ptr(2),
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": "demo",
                },
            },
            Template: apiv1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": "demo",
                    },
                },
                Spec: apiv1.PodSpec{
                    Containers: []apiv1.Container{
                        {
                            Name:  "web",
                            Image: "nginx:1.12",
                            Ports: []apiv1.ContainerPort{
                                {
                                    Name:          "http",
                                    Protocol:      apiv1.ProtocolTCP,
                                    ContainerPort: 80,
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create Deployment
    fmt.Println("Creating deployment...")
    result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())



//    to monitor what is going on with the resource, v1.ListOptions{Watch: true}

}

func int32Ptr(i int32) *int32 { return &i }
