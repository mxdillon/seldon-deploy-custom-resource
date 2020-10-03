package main

import (
    "context"
    "flag"
    "fmt"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "log"
    "seldon-deploy-custom-resource/pkg"
    ctrl "sigs.k8s.io/controller-runtime"
)

var k string

// init allows user to specify location of config file. Defaults to working directory.
func init() {

    flag.StringVar(&k, "kconfig", "./config", "Location of the k8s config file")
    flag.Parse()
}

func main() {

    clientset := pkg.SetupClient(k)
    api := clientset.CoreV1()

    // access the API to list pods
    pods, _ := api.Pods("").List(context.TODO(), metav1.ListOptions{})
    fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

    seldonDeployment := pkg.CreateDeployment("seldon-crd.json")


    // specify namespace to deploy into
    ns := &apiv1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: "seldon-crd",
        },
    }

    // set namespace in SeldonDeployment
    seldonDeployment.ObjectMeta.Namespace = ns.ObjectMeta.Name

    sch := runtime.NewScheme()

    // create manager for creating controllers
    cfg, err := ctrl.GetConfig()
    k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
        Scheme:         sch,
        LeaderElection: false,
        Namespace: ns.ObjectMeta.Name,
    })
    // create client
    k8sClient := k8sManager.GetClient()

    // add kind SeldonDeployment to scheme to enable creation
    if err = seldonv2.AddToScheme(sch); err != nil {
        log.Fatalf("Error - couldn't add SeldonDeployment to scheme: %v", err)
    }

    // add kind SeldonDeployment to scheme to enable creation
    if err = apiv1.AddToScheme(sch); err != nil {
        log.Fatalf("Error - couldn't add apiv1 to scheme: %v", err)
    }

    // if specified ns doesn't exist, create it
    //if err = k8sClient.Create(context.Background(), ns); err != nil {
    //    log.Fatalf("Error - couldn't create namespace: %v", err)
    //}

    // apply deployment to cluster
    if err = k8sClient.Create(context.Background(), seldonDeployment); err != nil {
        log.Fatalf("Error - couldn't create CRD: %v", err)
    }



}
