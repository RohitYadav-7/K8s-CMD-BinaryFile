# kubernetesCMD with GO-lang using Binary file 
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

**Clone the github code:**

```sh
git clone https://github.com/RohitYadav-7/K8s-CMD-BinaryFile.git
```

**Now bulid your binary file:**

```sh
 go build -o <yourbinaryfilename> main.go
```

**Example**

```sh
 go build -o k8sCMDbinary main.go 
 ```

## functions

### 1. List

**Description**
// TODO(user): List all runing pods in your Kubernetes cluster.

```sh
./k8sCMDbinary --list 
```
### 2. Describepod

**Description**
// TODO(user): Deacribe a specific pod in your Kubernetes cluster.

```sh
./k8sCMDbinary --describepod podname --namespace default
```
### 3. Cleanup 

**Description**
// TODO(user): all pods get deleted except the running pods in your Kubernetes cluster.

```sh
./k8sCMDbinary --cleanup
```
### 4. Showevents
**Description**
// TODO(user): Showevents in your Kubernetes cluster.
```sh
./k8sCMDbinary --showevents 
```
### 5. Showsecrets
**Description**
// TODO(user): Showsecrets in your Kubernetes cluster.
```sh
./k8sCMDbinary --showsecrets
```
### 6. Editreplicas
**Description**
// TODO(user): Editreplicas in your Kubernetes cluster.
```sh
./k8sCMDbinary --editreplicas=default:deploymentname:<replicanumbers>
```
### 7. Editingress
**Description**
// TODO(user): Editingress file in your Kubernetes cluster.
```sh
./k8sCMDbinary --editingress --namespace default --ingressname ingressName
```
### 8. Editdeployment
**Description**
// TODO(user): Editdeployment in your Kubernetes cluster.
```sh
./k8sCMDbinary --editdeployment=namespacename/deploymentname 
```
### 9. Describedeployment
**Description**
// TODO(user): describedeployment in your Kubernetes cluster.
```sh
./k8sCMDbinary  --describedeployment=namespacename/deploymentname    
```

### 10. Viewing Pod Logs
**Description**
// TODO(user): Viewing Pod Logs in your Kubernetes cluster.
```sh
./k8sCMDbinary  --viewlogs --pod podname --namespace namespacename --container containername
```


