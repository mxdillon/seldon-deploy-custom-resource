module seldon-deploy-custom-resource

go 1.15

require (
	github.com/seldonio/seldon-core/operator v0.0.0-20200930164230-41a462013a9c
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	sigs.k8s.io/controller-runtime v0.6.3
)
