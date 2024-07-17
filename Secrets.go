package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Define flags
	cleanupPtr := flag.Bool("cleanup", false, "Clean up resources")
	listPtr := flag.Bool("list", false, "List pods")
	showSecretsPtr := flag.Bool("showsecrets", false, "Show secrets")
	showEventsPtr := flag.Bool("showevents", false, "Show events")
	describePodPtr := flag.String("describepod", "", "Describe a specific pod")
	namespacePtr := flag.String("namespace", "default", "Namespace of the pod")
	editReplicasPtr := flag.String("editreplicas", "", "Edit replicas (format: <namespace>:<deployment-name>:<replica-count>)")
	ingressNamePtr := flag.String("ingressname", "", "Name of the Ingress resource")
	editIngressPtr := flag.Bool("editingress", false, "Edit Ingress resource")
	editDeploymentPtr := flag.String("editdeployment", "", "Edit deployment (namespace/deployment)")

	// Parse flags
	flag.Parse()

	// Create Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		fmt.Println("Error creating Kubernetes config:", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		os.Exit(1)
	}

	// Execute the appropriate function based on flags
	if *cleanupPtr {
		cleanup(clientset)
	} else if *listPtr {
		listPods(clientset)
	} else if *showSecretsPtr {
		showSecrets(clientset)
	} else if *showEventsPtr {
		showEvents(clientset)
	} else if *describePodPtr != "" {
		describePod(clientset, *describePodPtr, *namespacePtr)
	} else if *editReplicasPtr != "" {
		editReplicas(clientset, *editReplicasPtr)
	} else if *editIngressPtr {
		if *ingressNamePtr == "" {
			fmt.Println("Please provide the ingress name using --ingressname flag")
			os.Exit(1)
		}
		editIngress(clientset, *namespacePtr, *ingressNamePtr)
	} else if *editDeploymentPtr != "" {
		parts := strings.Split(*editDeploymentPtr, "/")
		if len(parts) != 2 {
			fmt.Println("Invalid format for editdeployment flag. Use namespace/deployment")
			os.Exit(1)
		}
		editDeployment(clientset, parts[0], parts[1])
	} else {
		flag.Usage()
	}
}

func cleanup(clientset *kubernetes.Clientset) {
	// Get the list of all pods in all namespaces
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	// Iterate over the list of pods and delete those in ImagePullBackOff state
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Reason == "ImagePullBackOff" {
				fmt.Printf("Deleting pod %s in namespace %s\n", pod.Name, pod.Namespace)
				err := clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				if err != nil {
					log.Printf("Error deleting pod %s in namespace %s: %v", pod.Name, pod.Namespace, err)
				}
				break
			}
		}
	}
}

func listPods(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing pods:", err)
		return
	}
	for _, pod := range pods.Items {
		fmt.Printf("Namespace: %s, Pod Name: %s\n", pod.Namespace, pod.Name)
	}
}

func showSecrets(clientset *kubernetes.Clientset) {
	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing secrets:", err)
		return
	}

	for _, secret := range secrets.Items {
		fmt.Printf("Secret Name: %s\n", secret.Name)
		for key, value := range secret.Data {
			decodedValue := base64.StdEncoding.EncodeToString(value)
			fmt.Printf("  Key: %s, Value: %s\n", key, decodedValue)
		}
	}
}

func showEvents(clientset *kubernetes.Clientset) {
	events, err := clientset.CoreV1().Events("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("Error retrieving events:", err)
		return
	}

	// Print event details
	for _, event := range events.Items {
		fmt.Printf("Event: %s\n", event.Message)
		fmt.Printf("Reason: %s\n", event.Reason)
		fmt.Printf("Type: %s\n", event.Type)
		fmt.Printf("Source: %s\n", event.Source.Component)
		fmt.Printf("Namespace: %s\n", event.Namespace)
		fmt.Printf("Involved Object: %s/%s\n", event.InvolvedObject.Kind, event.InvolvedObject.Name)
		fmt.Printf("First Timestamp: %s\n", event.FirstTimestamp)
		fmt.Printf("Last Timestamp: %s\n", event.LastTimestamp)
		fmt.Printf("Count: %s\n", event.Count)
		fmt.Println("----------------------------------------------------")
	}
}

func describePod(clientset *kubernetes.Clientset, podName, namespace string) {
	// Implement pod description logic
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error describing pod:", err)
		return
	}
	fmt.Printf("Pod Name: %s\nNamespace: %s\nNode Name: %s\nStatus: %s\n",
		pod.Name, pod.Namespace, pod.Spec.NodeName, pod.Status.Phase)
	fmt.Println("Labels:")
	for key, value := range pod.Labels {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println("Annotations:")
	for key, value := range pod.Annotations {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println("Containers:")
	for _, container := range pod.Spec.Containers {
		fmt.Printf("  Name: %s\n  Image: %s\n  Ports: %v\n", container.Name, container.Image, container.Ports)
	}
}
func editReplicas(clientset *kubernetes.Clientset, replicaParam string) {
	// Parse the replicaParam to get namespace, deployment name, and replica count
	parts := strings.Split(replicaParam, ":")
	if len(parts) != 3 {
		fmt.Println("Invalid format for editreplicas flag. Use format: <namespace>:<deployment-name>:<replica-count>")
		return
	}

	namespace := parts[0]
	deploymentName := parts[1]
	replicaCount, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Println("Invalid replica count:", err)
		return
	}

	// Fetch the deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error fetching deployment:", err)
		return
	}

	// Update the replica count
	replicas := int32(replicaCount)
	deployment.Spec.Replicas = &replicas

	// Apply the update
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error updating replicas:", err)
		return
	}

	fmt.Printf("Successfully updated replicas for deployment %s in namespace %s to %d\n", deploymentName, namespace, replicaCount)
}

func editIngress(clientset *kubernetes.Clientset, namespace, ingressName string) {
	// Fetch the Ingress resource
	ingress, err := clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error fetching ingress:", err)
		return
	}

	// Save the Ingress resource to a temporary file
	tmpFile, err := ioutil.TempFile("", "ingress-*.yaml")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write the Ingress resource to the file
	data, err := yaml.Marshal(ingress)
	if err != nil {
		fmt.Println("Error marshaling ingress to YAML:", err)
		return
	}
	if _, err := tmpFile.Write(data); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return
	}

	// Open the file in Vim for editing
	cmd := exec.Command("vim", tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error opening file in vim:", err)
		return
	}

	// Read the modified file
	modifiedData, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println("Error reading modified file:", err)
		return
	}

	// Unmarshal the modified data back into an Ingress resource
	var modifiedIngress networkingv1.Ingress
	if err := yaml.Unmarshal(modifiedData, &modifiedIngress); err != nil {
		fmt.Println("Error unmarshaling modified file:", err)
		return
	}

	// Ensure the namespace and name are set correctly
	modifiedIngress.Namespace = namespace
	modifiedIngress.Name = ingressName

	// Apply the changes back to the cluster
	_, err = clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), &modifiedIngress, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error updating ingress:", err)
		return
	}

	fmt.Println("Ingress resource updated successfully")
}

func editDeployment(clientset *kubernetes.Clientset, namespace, deploymentName string) {
	// Get the deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error getting deployment:", err)
		return
	}

	// Marshal the deployment to YAML
	data, err := yaml.Marshal(deployment)
	if err != nil {
		fmt.Println("Error marshaling deployment to YAML:", err)
		return
	}

	// Write the YAML to a temporary file
	tmpfile, err := ioutil.TempFile("", "deployment-*.yaml")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpfile.Name()) // Clean up the file afterward

	if _, err := tmpfile.Write(data); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return
	}

	// Open the temporary file in Vim
	cmd := exec.Command("vim", tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running Vim:", err)
		return
	}

	// Read the updated file
	updatedData, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		fmt.Println("Error reading updated file:", err)
		return
	}

	// Unmarshal the updated YAML to a deployment object
	updatedDeployment := &appsv1.Deployment{}
	if err := yaml.Unmarshal(updatedData, updatedDeployment); err != nil {
		fmt.Println("Error unmarshaling updated YAML:", err)
		return
	}

	// Update the deployment in the cluster
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), updatedDeployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error updating deployment:", err)
		return
	}

	fmt.Println("Deployment updated successfully.")
}
