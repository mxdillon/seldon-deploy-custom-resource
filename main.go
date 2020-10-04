package main

import (
    "context"
    "flag"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
    apiv1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/clientcmd"
    "log"
    "seldon-deploy-custom-resource/pkg"
    ctrl "sigs.k8s.io/controller-runtime"
    "time"
)

var k, f, nsName string

func init() {
    flag.StringVar(&k, "kconfig", "./config", "Location of the k8s config file")
    flag.StringVar(&f, "filename", "./seldon-crd.json", "Location of the Seldon CRD json file")
    flag.StringVar(&nsName, "ns", "seldon-crd", "Name of the namespace to deploy into")
    flag.Parse()
}

func main() {

    //clientset := pkg.SetupClient(k)
    //api := clientset.CoreV1()

    //// access the API to list pods
    //pods, _ := api.Pods("").List(context.Background(), metav1.ListOptions{})
    //fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

    seldonDeployment := pkg.CreateDeployment(f)
    ns := pkg.AddNamespace(seldonDeployment, nsName)

    sch := runtime.NewScheme()

    config, err := clientcmd.BuildConfigFromFlags("", k)

    // add Seldon api to scheme to enable creation of SeldonDeployment resource
    if err = seldonv2.AddToScheme(sch); err != nil {
        log.Fatalf("Error - couldn't add SeldonDeployment to scheme: %v", err)
    }
    // add apiv1 to scheme to enable creation of namespace
    if err = apiv1.AddToScheme(sch); err != nil {
        log.Fatalf("Error - couldn't add apiv1 to scheme: %v", err)
    }

    // initialise manager for creating controllers
    k8sManager, err := ctrl.NewManager(config, ctrl.Options{
        Scheme: sch,
        LeaderElection: false,
    })
    // create client
    k8sClient := k8sManager.GetClient()


    // if specified namespace doesn't exist, create it
    if err = k8sClient.Create(context.Background(), ns); err != nil {
        log.Printf("Namespace %v might already exist: %v", ns.ObjectMeta.Name, err)
    }

    // apply deployment to cluster
    if err = k8sClient.Create(context.Background(), seldonDeployment); err != nil {
        log.Fatalf("Error - couldn't create CRD: %v", err)
    }

    time.Sleep(time.Second * 20)

    // scale the resource to 2 replicas
    seldonDeployment.Spec.Predictors[0].Replicas = int32Ptr(2)

    // apply update to cluster
    if err = k8sClient.Update(context.Background(), seldonDeployment); err != nil {
       log.Fatalf("Error - couldn't update CRD: %v", err)
    }

    time.Sleep(time.Second * 20)

    // delete SeldonDeployment from cluster
    if err = k8sClient.Delete(context.Background(), seldonDeployment); err != nil {
        log.Fatalf("Error - couldn't delete CRD: %v", err)
    }

}

func int32Ptr(i int32) *int32 { return &i }
