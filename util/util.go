package util

import (
	"fmt"
	"strings"

	claimv1alpha1 "github.com/tmax-cloud/tfc-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// LowestNonZeroResult compares two reconciliation results
// and returns the one with lowest requeue time.
func LowestNonZeroResult(i, j ctrl.Result) ctrl.Result {
	switch {
	case i.IsZero():
		return j
	case j.IsZero():
		return i
	case i.Requeue:
		return i
	case j.Requeue:
		return j
	case i.RequeueAfter < j.RequeueAfter:
		return i
	default:
		return j
	}
}

func GetTerraformVariables(tfapply *claimv1alpha1.TFApplyClaim) string {

	if tfapply.Spec.Variable == "" {
		return ""
	}

	variableList := strings.Split(tfapply.Spec.Variable, ",")

	cmd := ""
	for _, v := range variableList {
		value := strings.Trim(v, " ")
		cmd += fmt.Sprintf("export TF_VAR_%s; ", value)
	}

	return cmd
}
