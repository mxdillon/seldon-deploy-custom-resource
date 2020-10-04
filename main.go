package main

import (
    "context"
    "flag"
    "github.com/mxdillon/seldon-deploy-custom-resource/pkg"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    "k8s.io/apimachinery/pkg/types"
    "log"
)

var (
    k, f, nsName string
    n            int
)

func init() {
    flag.StringVar(&k, "kconfig", "./config", "Location of the k8s config file")
    flag.StringVar(&f, "filename", "./seldon-crd.json", "Location of the Seldon CRD json file")
    flag.StringVar(&nsName, "ns", "seldon-crd", "Name of the namespace to deploy into")
    flag.IntVar(&n, "replicas", 2, "Number of replicas to scale up to")
    flag.Parse()
}

func main() {

    k8sClient, err := pkg.SetupClient(k)
    if err != nil {
        log.Fatalf("Error setting up client: %v", err)
    }

    seldonDeployment, err := pkg.CreateDeployment(f)
    if err != nil {
        log.Fatalf("Error creating Seldon Deployment CRD: %v", err)
    }

    ns := pkg.AddNamespace(seldonDeployment, nsName)

    // if specified namespace doesn't exist, create it
    if err := k8sClient.Create(context.Background(), ns); err != nil {
        log.Printf("Namespace %v might already exist: %v", ns.ObjectMeta.Name, err)
    }

    // apply deployment to cluster
    log.Printf("Applying CRD %v to cluster.", seldonDeployment.Name)
    if err := k8sClient.Create(context.Background(), seldonDeployment); err != nil {
        log.Fatalf("Error - couldn't create CRD: %v", err)
    }

    // wait for pod to become available
    log.Printf("Waiting for CRD pod to become available.")
    if err := pkg.WaitForStatus("Available", seldonDeployment, k8sClient); err != nil {
        log.Fatalf("Error waiting for CRD to be ready: %v", err)
    }

    // get latest state of seldonDeployment to update
    latest := &seldonv2.SeldonDeployment{}
    key := types.NamespacedName{
        Name:      seldonDeployment.Name,
        Namespace: seldonDeployment.Namespace,
    }
    if err := k8sClient.Get(context.Background(), key, latest); err != nil {
        log.Fatalf("Error - couldn't get latest info on CRD: %v", err)
    }

    // update number of replicas and apply to cluster
    latest.Spec.Predictors[0].Replicas = int32Ptr(int32(n))
    if err := k8sClient.Update(context.Background(), latest); err != nil {
        log.Fatalf("Error - couldn't update CRD: %v", err)
    }

    // wait for replicas to become available
    log.Printf("Waiting for all %v replicas of CRD to become available.", n)
    if err := pkg.WaitForReplicas(n, latest, k8sClient); err != nil {
        log.Fatalf("Error waiting for replicas to be available: %v", err)
    }

    // delete SeldonDeployment CRD from cluster
    if err := k8sClient.Delete(context.Background(), latest); err != nil {
        log.Fatalf("Error - couldn't delete CRD: %v", err)
    }
    log.Printf("Deleting CRD %v from namespace %v.", latest.Name, latest.Namespace)

}

func int32Ptr(i int32) *int32 { return &i }
