package pkg

import (
    "context"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    "k8s.io/apimachinery/pkg/types"
    "log"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "time"
)

// WaitForStatus halts execution until the specified status is returned by the resource.
func WaitForStatus(s string, sd *seldonv2.SeldonDeployment, c client.Client) {
    status := ""
    for status != s {
        latest := &seldonv2.SeldonDeployment{}
        key := types.NamespacedName{
            Name:      sd.Name,
            Namespace: sd.Namespace,
        }
        if err := c.Get(context.Background(), key, latest); err != nil {
            log.Fatalf("Error - couldn't get latest info on CRD: %v", err)
        }
        status = string(latest.Status.State)
        time.Sleep(time.Second * 2)
    }
    log.Printf("Received status %v from resource. Proceeding to next action.", status)
}

// WaitForReplicas halts execution until the specified number of replicas are available.
func WaitForReplicas(n int, sd *seldonv2.SeldonDeployment, c client.Client) {
    replicas := 0
    rName := sd.Name + "-" + sd.Spec.Predictors[0].Name + "-0-classifier"
    for replicas != n {
        latest := &seldonv2.SeldonDeployment{}
        key := types.NamespacedName{
            Name:      sd.Name,
            Namespace: sd.Namespace,
        }
        if err := c.Get(context.Background(), key, latest); err != nil {
            log.Fatalf("Error - couldn't get latest info on CRD: %v", err)
        }
        replicas = int(latest.Status.DeploymentStatus[rName].AvailableReplicas)
        time.Sleep(time.Second * 2)
    }
    log.Printf("%v replicas of CRD available. Proceeding to next action.", replicas)
}
