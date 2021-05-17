node {
    def gitHubBaseAddress = "github.com"
    def goHome = "/usr/local/go/bin"
    def buildDir = "/var/lib/jenkins/workspace/tfc-operator"  //

    def scriptHome = "${buildDir}/scripts"
	
    def gitAddress = "${gitHubBaseAddress}/tmax-cloud/tfc-operator.git"


    def version = "${params.majorVersion}.${params.minorVersion}.${params.tinyVersion}.${params.hotfixVersion}"
    def preVersion = "${params.preVersion}"

    def imageTag = "b${version}"
				
    def userName = "gyeongyeol-choi"
	
    def credentialsId = "gyeongyeol-choi"
    def userEmail = "gyeongyeol_choi@tmax.co.kr"
	
    stage('git clone') {	
	if(!fileExists(buildDir)){
	  sh "echo create build directory"
	  dir(buildDir){	
	      new File(buildDir).mkdir()
	      git branch: "${params.buildBranch}",
	      credentialsId: '${credentialsId}',
          url: "https://${gitAddress}"
	  }
	}
        sh "echo build directory is existed"
    }
	
    stage('git pull') { 
        dir(buildDir){
            // git pull

            sh "git checkout ${params.buildBranch}"
	    sh "git fetch --all"
            sh "git reset --hard origin/${params.buildBranch}"
            sh "git pull origin ${params.buildBranch}"

            sh '''#!/bin/bash
            export PATH=$PATH:/usr/local/go/bin
            export GO111MODULE=on
            go build -o bin/manager main.go
            '''
        }
    }
    
    stage('make manifests') {
	    sh "sed -i 's#{imageTag}#${imageTag}#' ./config/manager/kustomization.yaml"
        sh "sudo kubectl kustomize ./config/default/ > bin/tfc-operator-v${version}.yaml"
        sh "sudo kubectl kustomize ./config/crd/ > bin/crd-v${version}.yaml"
        sh "sudo tar -zvcf bin/tfc-operator-manifests-v${version}.tar.gz bin/tfc-operator-v${version}.yaml bin/crd-v${version}.yaml"
        
        sh "sudo mkdir -p manifests/v${version}"
        sh "sudo cp bin/*v${version}.yaml manifests/v${version}/"
    }

    stage('image build/push') {
        sh "sudo docker build --tag tmaxcloudck/tfc-operator:${imageTag} ."
        sh "sudo docker push tmaxcloudck/tfc-operator:${imageTag}"
        sh "sudo docker rmi tmaxcloudck/tfc-operator:${imageTag}"
    }

    stage('make-changelog') {
        sh "echo targetVersion: ${version}, preVersion: ${preVersion}"
        sh "sudo sh ${scriptHome}/make-changelog.sh ${version} ${preVersion}"
    }

    stage('git commit & push') {
        dir("${buildDir}") {
		
	
	    sh "git config --global user.name ${userName}"
            sh "git config --global user.email ${userEmail}"
	    sh "git config --global credential.helper store"		
		
            sh "git checkout ${params.buildBranch}"
            sh "git add -A"
			sh "git reset ./config/manager/kustomization.yaml"
            def commitMsg = "[Distribution] Release commit for tfc-operator v${version}"
            sh (script: "git commit -m \"${commitMsg}\" || true")
            sh "git tag v${version}"
	    sh "sudo git push -u origin +${params.buildBranch}"
            sh "sudo git push origin v${version}"
		
		
	    sh "git fetch --all"
            sh "git reset --hard origin/${params.buildBranch}"
	    sh "git pull origin ${params.buildBranch}"
        }
    }
}
