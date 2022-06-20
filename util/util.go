package util

import (
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
