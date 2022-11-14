package controllers

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tmax-cloud/tfc-operator/api/v1alpha1"
	claimv1alpha1 "github.com/tmax-cloud/tfc-operator/api/v1alpha1"

	"os"
)

func statusUpdate(t *claimv1alpha1.TFApplyClaim, update Status) {
	if update.Action.Set {
		t.Status.Action = update.Action.Value
	}
	if update.Apply.Set {
		t.Status.Apply = update.Apply.Value
	}
	if update.Branch.Set {
		t.Status.Branch = update.Branch.Value
	}
	if update.Commit.Set {
		t.Status.Commit = update.Commit.Value
	}
	if update.Destroy.Set {
		t.Status.Destroy = update.Destroy.Value
	}
	if update.Phase.Set {
		t.Status.Phase = update.Phase.Value
	}
	if update.PrePhase.Set {
		t.Status.PrePhase = update.PrePhase.Value
	}
	if update.Reason.Set {
		t.Status.Reason = update.Reason.Value
	}
	if update.State.Set {
		t.Status.State = update.State.Value
	}
	if update.URL.Set {
		t.Status.URL = update.URL.Value
	}

}

// deploymentForProvider returns a provider Deployment object
func (r *TFApplyClaimReconciler) deploymentForApply(m *claimv1alpha1.TFApplyClaim) (*appsv1.Deployment, error) {
	ls := labelsForApply(m.Name)
	replicas := int32(1) //m.Spec.Size
	image_path := os.Getenv("TFC_WORKER")
	if image_path == "" {
		return nil, fmt.Errorf("TFC_WORKER environment variable doesn't exist")
	}

	envList := []corev1.EnvVar{}

	// private type의 경우, secret env을 추가한다.
	if m.Spec.Type == "private" {
		envList = append(envList, corev1.EnvVar{
			Name: "GIT_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: m.Spec.Secret},
					Key:                  "token",
				},
			},
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           image_path, //"tmaxcloudck/tfc-worker:v0.0.1",
						Name:            "ubuntu",
						Command:         []string{"/bin/sleep", "3650d"},
						ImagePullPolicy: "Always",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "ubuntu",
						}},
						Env: envList,
					}},
				},
			},
		},
	}

	// Set Provider instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep, nil
}

// labelsForProvider returns the labels for selecting the resources
// belonging to the given Provider CR name.
func labelsForApply(name string) map[string]string {
	return map[string]string{"app": "tfapplyclaim", "tfapplyclaim_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func dequeuePlan(slice []v1alpha1.Plan, capacity int) []v1alpha1.Plan {
	fmt.Println(slice[1:])
	fmt.Println(slice[:capacity-1])
	return slice[:capacity-1]
}

// getPodName defines the pod name of terraform pod
func (r *TFApplyClaimReconciler) getPodName(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) error {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	// Update the Provider status with the pod names
	// List the pods for this provider's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(tfapplyclaim.Namespace),
		client.MatchingLabels(labelsForApply(tfapplyclaim.Name)),
		client.MatchingFields{"status.phase": "Running"},
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "TFApplyClaim.Namespace", tfapplyclaim.Namespace, "TFApplyClaim.Name", tfapplyclaim.Name)
		return err
	}
	podNames = getPodNames(podList.Items)

	if len(podNames) < 1 {
		log.Info("Not yet create Terraform Pod...")
		return fmt.Errorf("Not yet create Terraform Pod...")
	} else if len(podNames) > 1 {
		log.Info("Not yet terminate Previous Terraform Pod...")
		return fmt.Errorf("Not yet terminate Previous Terraform Pod...")
	} else {
		log.Info("Ready to Execute Terraform Pod!")
	}
	return nil
}

// createClientSet initializes the ClientSet of Kubernetes by InClusterConfig
func (r *TFApplyClaimReconciler) createClientSet(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) error {
	// creates the in-cluster config
	var err error
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())

	config, err = rest.InClusterConfig()
	if err != nil {
		log.Error(err, "Failed to create in-cluster config")
		statusUpdate(tfapplyclaim, Status{
			PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
			Phase:    OptionalString{"Error", true},
			Action:   OptionalString{"", true},
			Reason:   OptionalString{"Failed to create in-cluster config", true}})
		return err
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Failed to create clientset")
		statusUpdate(tfapplyclaim, Status{
			PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
			Phase:    OptionalString{"Error", true},
			Action:   OptionalString{"", true},
			Reason:   OptionalString{"Failed to create clientset", true}})
		return err
	}

	return nil
}

// adjustPodCount adjusts the count of Terraform Working Pods. If there is no need to create pods, it closes the pods.
func (r *TFApplyClaimReconciler) adjustPodCount(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim, terminate bool) error {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err := r.Get(ctx, tfapplyclaim.GetNamespacedName(), found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForApply(tfapplyclaim)
		if err != nil {
			log.Error(err, "Failed to create deployment")
			return err
		}

		log.Info("Creating a new Deployment", "Deployment.Namespace: ", dep.Namespace, "Deployment.Name: ", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace: ", dep.Namespace, "Deployment.Name: ", dep.Name)
			return err
		}
		// Deployment created successfully - return and requeue
		return fmt.Errorf("Deployment created successfully - Requeue")
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return err
	}

	// Ensure the deployment size is the same as the spec
	size := int32(1)

	if terminate {
		log.Info("There's no need to Create Terraform Pod...")
		size = 0
	}
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return err
		}
	}

	return nil
}

// check if deployment available
func (r *TFApplyClaimReconciler) checkDeploymentAvailable(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) error {
	dep := &appsv1.Deployment{}
	if err := r.Get(ctx, tfapplyclaim.GetNamespacedName(), dep); err != nil {
		return fmt.Errorf("Failed to get deployment")
	}

	for _, condition := range dep.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
			return nil
		}
	}

	return fmt.Errorf("Failed to get deployment")

}
