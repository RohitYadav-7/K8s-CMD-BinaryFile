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
```sh
./k8sCMDbinary --list 
```
### 2. Describepod
```sh
./k8sCMDbinary --describepod mypod1 --namespace default
```
### 3. Cleanup
```sh
./k8sCMDbinary --cleanup
```
### 4. Showevents
```sh
./k8sCMDbinary --showevents 
```
### 5. Showsecrets
```sh
./k8sCMDbinary --showsecrets
```
### 6. Editreplicas
```sh
./k8sCMDbinary --editreplicas=default:deployment1:3
```
### 7. Editingress
```sh
./k8sCMDbinary --editingress --namespace default --ingressname example-ingress
```
### 8. Editdeployment
```sh
./k8sCMDbinary --editdeployment=default/deployment3
```


