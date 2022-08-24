<!--
order: 1
-->

# Kubernetes local setup with Kind

### Pre-requisite : 
- `kubectl` For installation, follow the steps mentioned on the kubernetes [docs](https://kubernetes.io/docs/tasks/tools/)
- `kind` For installation, follow the steps mentioned on the Kind [docs](https://kind.sigs.k8s.io/docs/user/quick-start/#installation). 

### Steps for kind setup and change kube configs :

Once the installation is complete, follow the given steps to setup a local cluster :
```commandline
kind create cluster
```

Once the setup is complete, you will be able to see : 
```commandline
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.24.0) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
```

To see the cluster details in the config file of kube : 
```commandline
kubectl config view
```
And the output in the will be something like this : 
```commandline
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://127.0.0.1:59771
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
current-context: ""
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

To check the current-context :
```commandline
kubectl config current-context
```

If the current-context is set to nil or some other cluster in the above config, then switch the context to kind-kind using :
```commandline
kubectl config use-context kind-kind
```

### Ways to switch namespace :  
- First way is to make changes in the namespace present in the yaml files of the setup to be run. 
  So you can continue to use the default namespace of kind cluster.
    ```
    namespace: default
    ```
  Note : You have to change it at all the places where this namespace is used. Ex: http://gaia-genesis.dev-native.svc.cluster.local used in validator.yml
- Second way is to make a new namespace by using :
   ```commandline
   kubectl create namespace dev-native
   ```
  And then adding namespace in ~/.kube/config/<kube-config>.yaml of `cluster: kind-kind`
   ```
   namespace: dev-native
   ```
  So that when you use `kubectl config view`, context looks like :
   ```
   contexts:
   - context:
       cluster: kind-kind
       namespace: dev-native
       user: kind-kind
     name: kind-kind 
   ```

### Ex. : Setting up gaia chain
Once all the steps mentioned is complete then setup of gaia chain can be done by running inside the k8s directory.
```commandline
make apply PROCESS=gaia
```

### Ex. : Setting up pstake chain
Once all the steps mentioned is complete then setup of gaia chain can be done by running inside the k8s directory.
```commandline
make apply PROCESS=pstae
```
