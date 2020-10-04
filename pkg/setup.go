package pkg

import (
    "encoding/json"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    "io/ioutil"
    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/clientcmd"
    "log"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

// SetupClient creates and returns the k8s through which we interact with the cluster.
// Both the seldon and k8s.io/api apis are added to the scheme so we can manipulate resources on the cluster.
func SetupClient(k string) client.Client {
    config, err := clientcmd.BuildConfigFromFlags("", k)
    if err != nil {
        log.Fatalf("Error - failed to build config: %v", err)
    }

    sch := runtime.NewScheme()
    if err = seldonv2.AddToScheme(sch); err != nil {
        log.Fatalf("Error - couldn't add SeldonDeployment to scheme: %v", err)
    }
    if err = apiv1.AddToScheme(sch); err != nil {
        log.Fatalf("Error - couldn't add apiv1 to scheme: %v", err)
    }

    k8sClient, err := client.New(config, client.Options{Scheme: sch})
    if err != nil {
        log.Fatalf("Error - failed to create client: %v", err)
    }
    return k8sClient
}

// CreateDeployment reads the CRD json and populates a SeldonDeployment struct with the data.
func CreateDeployment(f string) *seldonv2.SeldonDeployment {
    sD := &seldonv2.SeldonDeployment{}

    if f[len(f)-5:] != ".json" {
        log.Fatalln("Error - must pass a json filename to CreateDeployment. Include the .json suffix.")
    }

    file, err := ioutil.ReadFile(f)
    if err != nil {
        log.Fatal("Error - failed to read json: %v", err)
    }

    if err = json.Unmarshal(file, &sD); err != nil {
        log.Fatal("Error - failed to unmarshall json into SeldonDeployment: %v", err)
    }
    return sD
}

// AddNamespace creates a namespace variable and adds the namespace name to the SeldonDeployment CRD.
func AddNamespace(sd *seldonv2.SeldonDeployment, n string) *apiv1.Namespace {
    ns := &apiv1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: n,
        },
    }
    sd.ObjectMeta.Namespace = ns.ObjectMeta.Name
    return ns
}