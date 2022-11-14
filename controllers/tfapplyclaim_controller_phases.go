/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	claimv1alpha1 "github.com/tmax-cloud/tfc-operator/api/v1alpha1"
	"github.com/tmax-cloud/tfc-operator/util"
)

func (r *TFApplyClaimReconciler) ReadyClaimPhase(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) (ctrl.Result, error) {

	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())

	log.Info("Start ReadyClaim")

	repoType := tfapplyclaim.Spec.Type
	secretName := tfapplyclaim.Spec.Secret
	secret := &corev1.Secret{}

	// Check the Secret (Git Credential) for Terraform HCL Code

	key := types.NamespacedName{Name: secretName, Namespace: tfapplyclaim.Namespace}

	if repoType == "private" {
		if err := r.Get(ctx, key, secret); errors.IsNotFound(err) {
			log.Error(err, "Credential secret doesn't exist")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{"", true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Credential secret doesn't exist", true}})
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		} else if err != nil {
			// Error reading the object - requeue the request.
			log.Error(err, "Failed to get Secret")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{"", true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to get Secret", true}})
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		_, exists_token := secret.Data["token"]

		if !exists_token {
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{"", true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Invalid Secret (token)", true}})
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		// error 상태인 tfapplyclaim을 수정한 경우
		if tfapplyclaim.Status.Phase == "Error" && (tfapplyclaim.Status.Reason == "Credential secret doesn't exist" ||
			tfapplyclaim.Status.Reason == "Failed to get Secret" ||
			tfapplyclaim.Status.Reason == "Invalid Secret (token)") {
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{"Error", true},
				Phase:    OptionalString{"Awaiting", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Changed to correct spec", true}})
			return ctrl.Result{}, nil
		}
	}

	if tfapplyclaim.Status.Phase == "" {
		statusUpdate(tfapplyclaim, Status{
			PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
			Phase:    OptionalString{"Awaiting", true}})
		return ctrl.Result{}, nil
	}

	// // If done Apply / Destroy, terraform pods is terminated.
	// if tfapplyclaim.Status.Phase == "Applied" || tfapplyclaim.Status.Phase == "Destroyed" && tfapplyclaim.Spec.Destroy == false {
	// 	_ = r.adjustPodCount(ctx, tfapplyclaim, true)
	// }

	return ctrl.Result{}, nil
}

// Status Update: action - reject
func (r *TFApplyClaimReconciler) RejectClaimPhase(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) (ctrl.Result, error) {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	log.Info("Start RejectClaim")
	statusUpdate(tfapplyclaim, Status{
		PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
		Phase:    OptionalString{"Rejected", true}})
	return ctrl.Result{}, nil
}

// Status Update: action - approve
func (r *TFApplyClaimReconciler) ApproveClaimPhase(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) (ctrl.Result, error) {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	log.Info("Start ApproveClaim")

	err = r.adjustPodCount(ctx, tfapplyclaim, false)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	err = r.getPodName(ctx, tfapplyclaim)
	if err != nil {
		log.Info("Waiting for pod to be created. requeue after 5 sec")
		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	err = r.createClientSet(ctx, tfapplyclaim)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Go Client - POD EXEC
	if tfapplyclaim.Status.Action == "Approve" && (tfapplyclaim.Status.Phase == "Awaiting" || tfapplyclaim.Status.Phase == "Rejected") {

		// 1. Git Clone Repository
		stdout.Reset()
		stderr.Reset()

		err = util.ExecClone(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		if err != nil && !strings.Contains(stdout.String(), "already exists") {
			log.Error(err, "Failed to Clone Git Repository")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Clone Git Repository", true}})
			return ctrl.Result{}, err
		}

		if tfapplyclaim.Spec.Branch != "" {
			stdout.Reset()
			stderr.Reset()

			err = util.ExecBranchCheckout(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

			fmt.Println(stdout.String())
			fmt.Println(stderr.String())

			if err != nil && !strings.Contains(stdout.String(), "already exists") {
				log.Error(err, "Failed to Checkout Git Branch")
				statusUpdate(tfapplyclaim, Status{
					PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
					Phase:    OptionalString{"Error", true},
					Action:   OptionalString{"", true},
					Reason:   OptionalString{"Failed to Checkout Git Branch", true}})
				return ctrl.Result{}, err
			}
		}

		// 2. Terraform Initialization
		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformDownload(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		if err != nil {
			log.Error(err, "Failed to Download Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Download Terraform", true}})
			return ctrl.Result{}, err
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformInit(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		if err != nil {
			log.Error(err, "Failed to Initialize Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Initialize Terraform", true}})
			return ctrl.Result{}, err
		} else {
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Approved", true}})
		}
	}
	tfapplyclaim.Status.Action = ""
	return ctrl.Result{}, nil
}

// Status Update: action - plan
func (r *TFApplyClaimReconciler) PlanClaimPhase(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) (ctrl.Result, error) {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	log.Info("Start PlanClaim")

	// 3. Terraform Plan
	if (tfapplyclaim.Status.Phase == "Approved" || tfapplyclaim.Status.Phase == "Planned") && tfapplyclaim.Status.Action == "Plan" {
		err = r.adjustPodCount(ctx, tfapplyclaim, false)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		if err := r.checkDeploymentAvailable(context.TODO(), tfapplyclaim); err != nil {
			log.Error(err, "Wait for deployment to be available")
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		err = r.getPodName(ctx, tfapplyclaim)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		err = r.createClientSet(ctx, tfapplyclaim)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Git Pull
		stdout.Reset()
		stderr.Reset()

		err := util.ExecGitPull(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		if err != nil {
			log.Error(err, "Failed to Pull Git Repository")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Pull Git Repository", true}})
			return ctrl.Result{}, err
		}

		// Get Commit ID
		stdout.Reset()
		stderr.Reset()

		err = util.ExecGetCommitID(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		log.Info(stdout.String())

		if err != nil {
			log.Error(err, "Failed to Get Commit ID")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Get Commit ID", true}})
			return ctrl.Result{}, err
		} else {
			commitID = strings.TrimRight(stdout.String(), "\r\n")
		}

		err = util.ExecTerraformInit(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		log.Info(stdout.String())
		log.Info(stderr.String())

		if err != nil {
			log.Error(err, "Failed to Initialize Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Initialize Terraform", true}})
			return ctrl.Result{}, err
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformPlan(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		stdoutStderr := stdout.String() + "\n" + stderr.String()

		// add plan to plans
		var plan claimv1alpha1.Plan

		plan.LastExectionTime = time.Now().Format("2006-01-02 15:04:05") // yyyy-MM-dd HH:mm:ss
		plan.Commit = commitID
		plan.Log = stdoutStderr

		if len(tfapplyclaim.Status.Plans) == capacity {
			tfapplyclaim.Status.Plans = dequeuePlan(tfapplyclaim.Status.Plans, capacity)
		}
		tfapplyclaim.Status.Plans = append([]claimv1alpha1.Plan{plan}, tfapplyclaim.Status.Plans...)

		if err != nil {
			log.Error(err, "Failed to Plan Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Plan Terraform", true}})
			return ctrl.Result{}, err
		} else {
			tfapplyclaim.Status.PrePhase = tfapplyclaim.Status.Phase
			tfapplyclaim.Status.Phase = "Planned"
		}

	}
	tfapplyclaim.Status.Action = ""
	return ctrl.Result{}, nil
}

// Status Update: action - apply
func (r *TFApplyClaimReconciler) ApplyClaimPhase(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) (ctrl.Result, error) {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	log.Info("Start ApplyClaim")
	// 4. Terraform Apply
	if (tfapplyclaim.Status.Phase == "Approved" || tfapplyclaim.Status.Phase == "Planned") && tfapplyclaim.Status.Action == "Apply" {

		err = r.adjustPodCount(ctx, tfapplyclaim, false)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		if err := r.checkDeploymentAvailable(context.TODO(), tfapplyclaim); err != nil {
			log.Error(err, "Wait for deployment to be available")
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		err = r.getPodName(ctx, tfapplyclaim)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		err = r.createClientSet(ctx, tfapplyclaim)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Get Commit ID
		stdout.Reset()
		stderr.Reset()

		err = util.ExecGetCommitID(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
		fmt.Println(stdout.String())

		if err != nil {
			log.Error(err, "Failed to Get Commit ID")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Get Commit ID", true}})
			return ctrl.Result{}, err
		} else {
			statusUpdate(tfapplyclaim, Status{
				Commit: OptionalString{strings.TrimRight(stdout.String(), "\r\n"), true},
				URL:    OptionalString{tfapplyclaim.Spec.URL, true},
				Branch: OptionalString{tfapplyclaim.Spec.Branch, true}})
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformApply(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		stdoutStderr := stdout.String() + "\n" + stderr.String()

		tfapplyclaim.Status.Apply = stdoutStderr

		if err != nil {
			log.Error(err, "Failed to Apply Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Apply Terraform", true}})
			return ctrl.Result{}, err
		}

		var matched string
		var added, changed, destroyed int

		lines := strings.Split(string(stdoutStderr), "\n")

		for i, line := range lines {
			if strings.Contains(line, "Apply complete!") {
				matched = lines[i]
				s := strings.Split(string(matched), " ")

				added, _ = strconv.Atoi(s[3])
				changed, _ = strconv.Atoi(s[5])
				destroyed, _ = strconv.Atoi(s[7])
			}
		}

		stdout.Reset()
		stderr.Reset()

		// Read Terraform State File
		err = util.ExecReadState(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
		fmt.Println(stdout.String())

		if err != nil {
			log.Error(err, "Failed to Read tfstate file")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Read tfstate file", true}})
			return ctrl.Result{}, err
		} else {
			statusUpdate(tfapplyclaim, Status{
				State:    OptionalString{stdout.String(), true},
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Applied", true}})
			tfapplyclaim.Status.Resource.Added = added
			tfapplyclaim.Status.Resource.Updated = changed
			tfapplyclaim.Status.Resource.Deleted = destroyed

			// Add finalizer first if not exist to avoid the race condition between init and delete
			if !controllerutil.ContainsFinalizer(tfapplyclaim, "claim.tmax.io/terraform-protection") {
				controllerutil.AddFinalizer(tfapplyclaim, "claim.tmax.io/terraform-protection")
			}
		}
	}
	tfapplyclaim.Status.Action = ""
	return ctrl.Result{}, nil
}

// Spec Update: destroy - true (This fuction is triggered by all users)
func (r *TFApplyClaimReconciler) DestroyClaimPhase(ctx context.Context, tfapplyclaim *claimv1alpha1.TFApplyClaim) (ctrl.Result, error) {
	log := r.Log.WithValues("tfapplyclaim", tfapplyclaim.GetNamespacedName())
	log.Info("Start DestroyClaim")
	// 5. Terraform Destroy (if required)
	if tfapplyclaim.Status.Phase == "Applied" && tfapplyclaim.Spec.Destroy == true {
		err = r.adjustPodCount(ctx, tfapplyclaim, false)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		if err := r.checkDeploymentAvailable(context.TODO(), tfapplyclaim); err != nil {
			log.Error(err, "Wait for deployment to be available")
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		err = r.getPodName(ctx, tfapplyclaim)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		err = r.createClientSet(ctx, tfapplyclaim)
		if err != nil {
			return ctrl.Result{}, err
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecClone(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		if err != nil && !strings.Contains(stdout.String(), "already exists") {
			log.Error(err, "Failed to Clone Git Repository")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Clone Git Repository", true}})
			return ctrl.Result{}, err
		}

		if tfapplyclaim.Spec.Branch != "" {
			stdout.Reset()
			stderr.Reset()

			err = util.ExecBranchCheckout(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
			fmt.Println(stdout.String())
			fmt.Println(stderr.String())

			if err != nil && !strings.Contains(stdout.String(), "already exists") {
				log.Error(err, "Failed to Checkout Git Branch")
				statusUpdate(tfapplyclaim, Status{
					PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
					Phase:    OptionalString{"Error", true},
					Action:   OptionalString{"", true},
					Reason:   OptionalString{"Failed to Checkout Git Branch", true}})
				return ctrl.Result{}, err
			}
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformDownload(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		if err != nil {
			log.Error(err, "Failed to Initialize Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Initialize Terraform", true}})
			return ctrl.Result{}, err
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformInit(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)

		fmt.Println(stdout.String())
		fmt.Println(stderr.String())

		if err != nil {
			log.Error(err, "Failed to Initialize Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Initialize Terraform", true}})
			return ctrl.Result{}, err
		}

		// Revert to Commit Point
		stdout.Reset()
		stderr.Reset()

		err = util.ExecRevertCommit(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
		fmt.Println(stdout.String())

		if err != nil {
			log.Error(err, "Failed to Revert Commit")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Revert Commit", true}})
			return ctrl.Result{}, err
		}

		// Recover Terraform State
		stdout.Reset()
		stderr.Reset()

		err = util.ExecRecoverState(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
		fmt.Println(stdout.String())

		if err != nil {
			log.Error(err, "Failed to Recover Terraform State")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Recover Terraform State", true}})
			return ctrl.Result{}, err
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecTerraformDestroy(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
		log.Info(stdout.String())
		log.Info(stderr.String())

		stdoutStderr := stdout.String() + "\n" + stderr.String()

		tfapplyclaim.Status.Destroy = stdoutStderr

		if err != nil {
			log.Error(err, "Failed to Destroy Terraform")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Action:   OptionalString{"", true},
				Reason:   OptionalString{"Failed to Destroy Terraform", true}})
			return ctrl.Result{}, err
		}

		var matched string
		var added, changed, destroyed int

		lines := strings.Split(string(stdoutStderr), "\n")

		for i, line := range lines {
			if strings.Contains(line, "Destroy complete!") {
				matched = lines[i]
				s := strings.Split(string(matched), " ")

				destroyed, _ = strconv.Atoi(s[3])
			}
		}

		stdout.Reset()
		stderr.Reset()

		err = util.ExecReadState(clientset, config, podNames[0], tfapplyclaim.Namespace, nil, &stdout, &stderr, tfapplyclaim)
		fmt.Println(stdout.String())

		if err != nil {
			log.Error(err, "Failed to Read tfstate file")
			statusUpdate(tfapplyclaim, Status{
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Error", true},
				Reason:   OptionalString{"Failed to Read tfstate file", true}})
			return ctrl.Result{}, err
		} else {
			statusUpdate(tfapplyclaim, Status{
				State:    OptionalString{stdout.String(), true},
				PrePhase: OptionalString{tfapplyclaim.Status.Phase, true},
				Phase:    OptionalString{"Destroyed", true}})

			tfapplyclaim.Status.Resource.Added = added
			tfapplyclaim.Status.Resource.Updated = changed
			tfapplyclaim.Status.Resource.Deleted = destroyed

			tfapplyclaim.Spec.Destroy = false

			// Remove finalizer if there is no need to maintain this resource
			controllerutil.RemoveFinalizer(tfapplyclaim, "claim.tmax.io/terraform-protection")
		}
	}
	tfapplyclaim.Status.Action = ""
	return ctrl.Result{}, nil
}
