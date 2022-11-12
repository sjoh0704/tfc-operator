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

package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var TFApplyClaimWebhookLogger = logf.Log.WithName("tfapplyclaim-resource")

func (r *TFApplyClaim) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-claim-tmax-io-v1alpha1-tfapplyclaim,mutating=true,failurePolicy=fail,groups=claim.tmax.io,resources=tfapplyclaims,verbs=update,versions=v1alpha1,name=mutation.webhook.tfapplyclaim,admissionReviewVersions=v1beta1;v1,sideEffects=NoneOnDryRun

var _ webhook.Defaulter = &TFApplyClaim{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *TFApplyClaim) Default() {
	TFApplyClaimWebhookLogger.Info("default", "name", r.Name, "namespace", r.Namespace)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-claim-tmax-io-v1alpha1-tfapplyclaim,mutating=false,failurePolicy=fail,groups=claim.tmax.io,resources=tfapplyclaims;tfapplyclaims/status,versions=v1alpha1,name=validation.webhook.tfapplyclaim,admissionReviewVersions=v1beta1;v1,sideEffects=NoneOnDryRun

var _ webhook.Validator = &TFApplyClaim{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TFApplyClaim) ValidateCreate() error {
	TFApplyClaimWebhookLogger.Info("validate create", "name", r.Name, "namespace", r.Namespace)

	if r.Spec.Destroy {
		return errors.New("Cannot set spec.destroy as TRUE when create TFApplyClaim at first")
	}

	if r.Spec.Type == "private" && r.Spec.Secret == "" {
		return errors.New("In order to use the private type, need to fill out the secret name in your namespace.")
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TFApplyClaim) ValidateUpdate(old runtime.Object) error {
	// TFApplyClaimWebhookLogger.Info("validate update", "name", r.Name, "namespace", r.Namespace)

	oldTfc := old.(*TFApplyClaim).DeepCopy()
	if oldTfc.Status.Phase == "Error" {
		return errors.New("Cannot change it when it is an error phase.")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TFApplyClaim) ValidateDelete() error {
	TFApplyClaimWebhookLogger.Info("validate delete", "name", r.Name, "namespace", r.Namespace)

	// 프로비저닝된 상태에서 tfapplyclaim 삭제 요청시 금지
	if r.Status.Phase == "Applied" {
		return errors.New("Destroying action must precede deleting tfapplyclaim resource")
	}

	return nil
}
