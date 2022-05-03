# TFC-Operator
TFC-Operator is the Kubenretes Native Operator with GitOPs of Terraform. Supports Github/Gitlab (Public/Private), and uses HCL (Hashicorp Configuration Language) to provision resources in Cloud Platform (e.g. AWS, Azure, GCP, vSPhere.. etc)

# Quick Start
This guide contains information on creating an EC2 instance on AWS through Terraform Claims (CR).

## 1. Public Repo
(1) Push the HCL (Hashicorp Configuration Language) code that defines Terraform resources to the git repository (Github / GitLab)

### Example.
```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  access_key = "<AWS_ACCESS_KEY_ID>"
  secret_key = "<AWS_SECRET_ACCESS_KEY>"
  region = "us-east-1"
}

resource "aws_instance" "app_server" {
  ami           = "ami-0003076ab1664cd23"
  instance_type = "t2.micro"

  tags = {
    Name = "ExampleAppServerInstance"
  }
}
```

(2) Create a Terraform claim based on the previously created Git Repo information
```bash
kubectl apply -f samples/public/01_claim.yaml
```

## 2. Private Repo
(1) Push the HCL (Hashicorp Configuration Language) code that defines Terraform resources to the git repository (Github / GitLab) 
### Example. (same as 1-(1))
```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  access_key = "<AWS_ACCESS_KEY_ID>"
  secret_key = "<AWS_SECRET_ACCESS_KEY>"
  region = "us-east-1"
}

resource "aws_instance" "app_server" {
  ami           = "ami-0003076ab1664cd23"
  instance_type = "t2.micro"

  tags = {
    Name = "ExampleAppServerInstance"
  }
}
```

(2) Create a secret holding an access token that can access the private repo
```bash
kubectl apply -f samples/private/01_secret.yaml
```
(3) Create a Terraform claim based on the previously created Git Repo information
```bash
kubectl apply -f samples/private/02_claim.yaml
```

