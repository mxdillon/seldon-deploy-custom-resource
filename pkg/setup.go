package pkg

import (
    "encoding/json"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    "io/ioutil"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "log"
)

// SetupClient creates the clientset and returns the API to interact with the k8s cluster.
func SetupClient(k string) *kubernetes.Clientset {
    config, err := clientcmd.BuildConfigFromFlags("", k)
    if err != nil {
        log.Fatalf("Error - failed to build config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error - failed to create clientset: %v", err)
    }

    return clientset
}

// CreateDeployment reads the CRD json and populates a SeldonDeployment struct with the data.
func CreateDeployment(f string) *seldonv2.SeldonDeployment {
    sD := &seldonv2.SeldonDeployment{}

    if f[len(f)-5:] != ".json" {
        log.Fatalln("Error - must pass a json filename to createDeployment. Include the .json suffix.")
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
