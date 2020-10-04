# seldon-deploy-custom-resource

This module:

* Creates a standalone program in Go which takes in a Seldon Core Custom Resource and creates it over the Kubernetes API;
* Watchs the created resource to wait for it to become available;
* Scales the resource to 2 replicas;
* When it is available deletes the Custom Resource;
* [_NOT IMPLEMENTED_] In parallel to the last 3 steps lists the Kubernetes Events emitted by the created resource until it is deleted.



-------
### Instructions for use

#### Prerequisites:
1. This module has been written in golang v1.15 and has not been tested on previous versions.
2. You will need a kubernetes cluster with Seldon Core deployed onto it. This module was developed using
[kind](https://kind.sigs.k8s.io/docs/user/quick-start/) and has not been tested on any other k8s infrastructure.
3. Ensure you are connected to the context of the k8s cluster you wish to manipulate. Use `kubectl config get-contexts` to verify.
4. Copy your kubeconfig file into this module directory. The default location for the file is `~/.kube/config`.
You can copy it anywhere into the directory, but you must use the command line flag `-kconfig=<relativepath>` when
running the module if so. Leaving the argument blank defaults to the top level of the project directory,
ie `-kconfig=./config`.
5. The Seldon CRD used for this module can be found
[here](https://raw.githubusercontent.com/SeldonIO/seldon-core/master/notebooks/resources/model.json).
This can be overwritten, or another json can be pointed to with the `filename` command line flag.


#### Running the module:
1. Clone this repo. Copy in your config file and adjust the json if desired (see above).
2. In the root of the project directory run the following commands:
    ```
    go get -d -v ./
    go install -v ./
    go test ./pkg
    go build -o ./bin/seldon-deploy
    ```
    These commands: download and install the required modules; run the suite of unit tests; and build the binary.
3. Run the module with the following command:
    ```
    ./bin/seldon-deploy <optional flags>
    ```
   Where the flags are as follows:

   | flag | description | example | default value |
   | :---: | :---: |:---: | :---: |
   | kconfig | Location of the k8s config file. | `-kconfig ./foo/config` | `-kconfig ./config`  |
   | filename | Location of the Seldon CRD json file. | `filename ./foo/my-crd.json` | `-filename ./seldon-crd.json`  |
   | ns | Name of the namespace to deploy into. If it doesn't exist already this module will create it. | `-ns hello` | `-ns seldon-crd`  |
   | replicas | Number of replicas to scale up to. | `-replicas 10` | `-replicas 2`  |


-------
### Extensions
With more time I would look into:
* **k8s events** - I was unable to implement event-watching to replicate `kubectl get events -w` in the time
available. I have left some commented out code in main.go for reference, but my approach was to be as follows:
    * Start a goroutine with an event listener in another thread (so it could run in parallel with the main function);
    * Send any events back into the main goroutine via a channel to be recorded / printed to stdout;
    * Close down both routines once the last messages about the deletion of the CRD pods was received.
* **Testing** - This module has very low test coverage and would ideally be close to 100%. There is a trade off to be
made between achieving a high unit test coverage and re-writing tests that the module authors (namely seldon and k8s)
have already implemented on the imported modules themselves. Therefore in this case I would look at writing BDD tests,
for example using [godog](https://github.com/cucumber/godog).
* **Running Environemnt** - It would be easier for others to use this module if a consistent runtime environment could
be ensured.
One way of ensuring consistency would be by building the binary and deploying that to users, however this isn't appropriate
in this case as the user must supply the k8s config file and may want to adjust the CRD.
A second option I would investigate would be including a Dockerfile that the user could run once they have copied their
k8s config into the directory. I'd need to think about the various networking implications of this approach.
