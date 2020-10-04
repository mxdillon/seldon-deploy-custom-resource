# .Background-deploy-custom-resource





-------
### Instructions for use

Prerequisites:
1. You will need a kubernetes cluster with Seldon Core deployed on it.
1. Ensure you are connected to the context of this cluster. Use `kubectl config get-contexts` to verify.
2. Copy your kubeconfig file into this module directory. The default location for
the file is `~/.kube/config`.
You can copy it anywhere into the directory, but you
must use the command line flag `-kconfig=<relativepath>` when running the module. Leaving
the argument blank defaults to the top level, ie `-kconfig=./config`.
3. The Seldon CRD used for this module can be found [here](https://raw.githubusercontent.com/SeldonIO/seldon-core/master/notebooks/resources/model.json).
This can be pasted over with another CRD if desired.
