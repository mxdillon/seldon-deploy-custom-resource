package pkg

import (
    "encoding/json"
    "errors"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    "io/ioutil"
    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/clientcmd"
    "log"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

// SetupClient creates and returns the k8s client through which we interact with the cluster.
// Both the seldon and k8s.io/api apis are added to the scheme so we can manipulate resources on the cluster.
func SetupClient(k string) (client.Client, error) {
    config, err := clientcmd.BuildConfigFromFlags("", k)
    if err != nil {
        return nil, err
    }

    sch := runtime.NewScheme()
    if err = seldonv2.AddToScheme(sch); err != nil {
        return nil, err
    }
    if err = apiv1.AddToScheme(sch); err != nil {
        return nil, err
    }

    k8sClient, err := client.New(config, client.Options{Scheme: sch})
    if err != nil {
        return nil, err
    }
    log.Printf("Configured the k8s client.")
    return k8sClient, nil
}

// CreateDeployment reads the CRD json and populates a SeldonDeployment struct with the data.
// Includes error handling for when a non-json file is passed to the function.
func CreateDeployment(f string) (*seldonv2.SeldonDeployment, error) {
    sD := &seldonv2.SeldonDeployment{}

    if f[len(f)-5:] != ".json" {
        err := errors.New("Error - must pass a json filename to CreateDeployment. Include the .json suffix.")
        return nil, err
    }

    file, err := ioutil.ReadFile(f)
    if err != nil {
        return nil, err
    }

    if err = json.Unmarshal(file, &sD); err != nil {
        return nil, err
    }
    log.Printf("Created Seldon CRD from json.")
    return sD, nil
}

// AddNamespace creates a namespace variable and adds the namespace name to the SeldonDeployment CRD.
func AddNamespace(sd *seldonv2.SeldonDeployment, n string) *apiv1.Namespace {
    ns := &apiv1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: n,
        },
    }
    sd.ObjectMeta.Namespace = ns.ObjectMeta.Name
    log.Printf("Added namespace %v to seldon CRD definition.", n)
    return ns
}
