# go-web3-sample-app
Web3 sample app in Go

## Requirements

* [gvm](https://github.com/moovweb/gvm) Go 1.19
    ```bash
    gvm install go1.19 --prefer-binary --with-build-tools --with-protobuf
    gvm use go1.19 --default
    ```
* [curl](https://help.ubidots.com/en/articles/2165289-learn-how-to-install-run-curl-on-windows-macosx-linux)
* [wget](https://www.gnu.org/software/wget/)
* [jq](https://github.com/stedolan/jq/wiki/Installation)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
* [kind >=0.16.0](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

## Kubernetes deployment

### Qleet Kind cluster

#### Deploy workload

```bash
make kind-deploy
```

#### Map the port of the sample app:

```bash
kubectl port-forward -n web3 svc/web3-sample-app 8080:8080
```

#### Open sample app in a browser

[http://localhost:8080/](http://localhost:8080/)

#### Undeploy workload

```bash
make kind-undeploy
```

### Any Kubernetes cluster

#### Deploy workload

```bash
kubectl apply -f ./k8s --namespace=web3 --validate=false
```

#### Delete workload

```bash
kubectl delete -f ./k8s --namespace=web3
```

## Build
Run `build` target
```bash
make build
```

## Release
Run `release` target
```bash
make release
```

## Help

```text
$ make
Usage: make COMMAND
Commands :
help          - List available tasks
clean         - Cleanup
test          - Run tests
build         - Build binary
run           - Run binary
get           - Download and install dependency packages
release       - Create and push a new tag
update        - Update dependencies to latest versions
version       - Print current version(tag)
image-build   - Build a Docker image
image-run     - Run a Docker image
image-stop    - Stop a Docker image
kind-deploy   - Deploy to a local KinD cluster
kind-undeploy - Undeploy from a local KinD cluster
kind-redeploy - Redeploy to a local KinD cluster
```
