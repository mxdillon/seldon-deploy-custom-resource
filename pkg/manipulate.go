package pkg

import (
    "context"
    seldonv2 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha2"
    "k8s.io/apimachinery/pkg/types"
    "log"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "time"
)

// WaitForStatus halts execution until the specified status is returned by the resource. Checks every 2 seconds.
func WaitForStatus(s string, sd *seldonv2.SeldonDeployment, c client.Client) error {
    status := ""
    for status != s {
        latest := &seldonv2.SeldonDeployment{}
        key := types.NamespacedName{
            Name:      sd.Name,
            Namespace: sd.Namespace,
        }
        if err := c.Get(context.Background(), key, latest); err != nil {
            return err
        }
        status = string(latest.Status.State)
        time.Sleep(time.Second * 2)
    }
    log.Printf("Received status %v from resource. Proceeding to next action.", status)
    return nil
}

// WaitForReplicas halts execution until the specified number of replicas are available. Checks every 2 seconds.
// This function is brittle - it requires that the resource name end "-0-classifier" so may need adjusting for other
// Seldon CRDs.
func WaitForReplicas(n int, sd *seldonv2.SeldonDeployment, c client.Client) error {
    replicas := 0
    rName := sd.Name + "-" + sd.Spec.Predictors[0].Name + "-0-classifier"
    for replicas != n {
        latest := &seldonv2.SeldonDeployment{}
        key := types.NamespacedName{
            Name:      sd.Name,
            Namespace: sd.Namespace,
        }
        if err := c.Get(context.Background(), key, latest); err != nil {
            return err
        }
        replicas = int(latest.Status.DeploymentStatus[rName].AvailableReplicas)
        time.Sleep(time.Second * 2)
    }
    log.Printf("%v replicas of CRD available. Proceeding to next action.", replicas)
    return nil
}
