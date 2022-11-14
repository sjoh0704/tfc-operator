package controllers

import (
	"bytes"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var capacity int = 5
var commitID string

type Status struct {
	Action   OptionalString
	Apply    OptionalString
	Branch   OptionalString
	Commit   OptionalString
	Destroy  OptionalString
	Phase    OptionalString
	PrePhase OptionalString
	Reason   OptionalString
	State    OptionalString
	URL      OptionalString
}

type OptionalInt struct {
	Value int
	//Null  bool
	Set bool
}

type OptionalString struct {
	Value string
	//Null  bool
	Set bool
}

/* Pod Exec Logs (Output / Error) */
var stdout bytes.Buffer
var stderr bytes.Buffer

/* Clientset for Kubernetes */
var clientset *kubernetes.Clientset

/* Terraform Pods List */
var podNames []string

/* Kubernetes Config (.config) */
var config *rest.Config

var err error
