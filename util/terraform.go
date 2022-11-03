package util

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	claimv1alpha1 "github.com/tmax-cloud/tfc-operator/api/v1alpha1"
)

func ExecClone(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	var protocol string

	if strings.Contains(apply.Spec.URL, "http://") {
		protocol = "http://"
	} else if strings.Contains(apply.Spec.URL, "https://") {
		protocol = "https://"
	}

	url := strings.TrimLeft(apply.Spec.URL, protocol)

	cmd := ""
	//github인 경우와 gitlab인 경우 token url 순서가 다름
	if strings.Contains(apply.Spec.URL, "github") {
		if apply.Spec.Type == "private" {
			url = protocol + "$GIT_TOKEN:x-oauth-basic" + "@" + url
		} else {
			url = protocol + "TMP_TOKEN:x-oauth-basic" + "@" + url
		}
		rep_token := "GIT_TOKEN=$(echo $GIT_TOKEN | sed -e 's/\\!/%21/g' -e 's/\\#/%23/g' -e 's/\\$/%24/g' -e 's/\\&/%26/g' -e \"s/'/%27/g\" -e 's/(/%28/g' -e 's/)/%29/g' -e 's/\\*/%2A/g'  -e 's/\\+/%2B/g' -e 's/\\,/%2C/g' -e 's/\\//%2F/g' -e 's/\\:/%3A/g' -e 's/\\;/%3B/g' -e 's/\\=/%3D/g' -e 's/\\?/%3F/g' -e 's/\\@/%40/g' -e 's/\\[/%5B/g'  -e 's/\\]/%5D/g');"
		cmd = rep_token + "git clone " + url + " " + HCL_DIR

	} else {
		if apply.Spec.Type == "private" {
			url = protocol + "x-oauth-basic:$GIT_TOKEN" + "@" + url
		} else {
			url = protocol + "x-oauth-basic:TMP_TOKEN" + "@" + url
		}
		cmd = "git clone " + url + " " + HCL_DIR

	}

	// SSL 미인증 VCS (e.g. Gitlab)을 위한 로직 처리
	cmd = "git config --global http.sslVerify false;" + cmd

	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecBranchCheckout(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" + "git checkout -t origin/" + apply.Spec.Branch
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecTerraformDownload(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	version := apply.Spec.Version

	cmd := "cd /tmp;" +
		fmt.Sprintf("wget %s/%s/%s_%s_linux_amd64.zip;", TERRAFORM_BINARY_URL, version, TERRAFORM_BINARY_NAME, version) +
		fmt.Sprintf("wget %s/%s/%s_%s_SHA256SUMS;", TERRAFORM_BINARY_URL, version, TERRAFORM_BINARY_NAME, version) +
		fmt.Sprintf("unzip -o -d /bin %s_%s_linux_amd64.zip;", TERRAFORM_BINARY_NAME, version) +
		"rm -rf /tmp/build;"

	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecTerraformInit(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	var cmd string

	versions := strings.Split(apply.Spec.Version, ".")

	majorVersion, _ := strconv.Atoi(versions[0])
	minorVersion, _ := strconv.Atoi(versions[1])

	if int(majorVersion) >= 1 || minorVersion >= 15 {
		cmd = "cd " + HCL_DIR + ";" + "terraform init"
	} else {
		cmd = "cd " + HCL_DIR + ";" + "terraform init -verify-plugins=false"
	}

	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecGitPull(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" + "git pull"
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecGetCommitID(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" +
		"git log --pretty=format:\"%H\" | head -n 1"

	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecCreateVariables(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" + "cat > terraform.tfvars.json << EOL\n" +
		apply.Spec.Variable + "\nEOL"

	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecTerraformPlan(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" + "terraform plan"
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecTerraformApply(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" + "terraform apply -auto-approve"
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecReadState(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" +
		"cat terraform.tfstate"
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecRevertCommit(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" +
		"git reset " + apply.Status.Commit
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecRecoverState(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" +
		"cat > terraform.tfstate << EOL\n" +
		apply.Status.State + "EOL"
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

func ExecTerraformDestroy(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer, apply *claimv1alpha1.TFApplyClaim) error {

	cmd := "cd " + HCL_DIR + ";" + "terraform destroy -auto-approve"
	err := execPodCmd(client, config, podName, podNamespace, cmd, nil, stdout, stderr)

	if err != nil {
		return err
	}
	return nil
}

// execPodCmd exec command on specific pod and wait the command's output.
func execPodCmd(client kubernetes.Interface, config *restclient.Config, podName string, podNamespace string,
	command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	cmd := []string{
		//"sh",
		"/bin/sh",
		"-c",
		command,
	}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(podNamespace).SubResource("exec")

	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if stdin == nil {
		option.Stdin = false
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}
	return nil
}

// ReadIDFromFile returns a Cloud Resource ID from Terraform State File
func ReadIDFromFile(filename string) (string, error) {
	var matched string // line with id
	var id string

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err, "Failed to read Terraform State File")
		return "", err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "\"id\"") {
			matched = lines[i]
		}
	}

	t := strings.Split(matched, "\"")

	if len(t) >= 4 {
		id = t[3]
	} else {
		err = errors.New("Index out of range")
		return "", err
	}

	return id, nil
}

// Initialize Terraform Working Directory
func InitTerraform_CLI(targetDir string, cloudType string) error {
	// Select the Terraform Plugin (cloudType: AWS, Azure, GCP)
	orgDir := HCL_DIR + "/" + ".terraform" + cloudType
	dstDir := targetDir + "/" + ".terraform"

	// Make the Destination Directory for plugin
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err = os.Mkdir(dstDir, 0755)
		if err != nil {
			return err
		}
		// Copy the Terraform Plugin (e.g. AWS, Azure, GCP) at Woring Directory
		err = copy(orgDir, dstDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute Terraform (Apply / Destroy)
func ExecuteTerraform_CLI(targetDir string, isDestroy bool) error {

	// Provision the Resources by Terraform
	cmd := exec.Command("bash", "-c", "terraform apply -auto-approve")

	// Swith the command from "apply" to "destroy"
	if isDestroy {
		// Destroy the Reosource by Terraform
		cmd = exec.Command("bash", "-c", "terraform destroy -auto-approve")
	}

	cmd.Dir = targetDir
	stdoutStderr, err := cmd.CombinedOutput()

	fmt.Printf("%s\n", stdoutStderr)

	return err
}

// Copy a Dierectory (preserve directory structure)
func copy(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			var data, err1 = ioutil.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(destination, relPath), data, 0777)
		}
	})
	return err
}
