#!/usr/bin/groovy
import java.text.SimpleDateFormat

podTemplate(label: 'jenkins-pipeline', containers: [
    containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:2.62', args: '${computer.jnlpmac} ${computer.name}', workingDir: '/home/jenkins', resourceRequestCpu: '200m', resourceLimitCpu: '200m', resourceRequestMemory: '256Mi', resourceLimitMemory: '256Mi'),
    containerTemplate(name: 'golang', image: 'golang:1.7.5', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'docker', image: 'docker:17.06.0', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'kubectl', image: 'lachlanevenson/k8s-kubectl:v1.7.4', command: 'cat', ttyEnabled: true)
],
volumes:[
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock')
])
    {
        node ('jenkins-pipeline') {
            println "DEBUG: Pipeline starting"
        
            // grab repo from source control
            checkout scm

            // configuration parameters and variables for pipeline
            // def pwd = pwd()
            def repo = "chzbrgr71"
            def appMajorVersion = "1.0"
            def acrServer = "briar123.azurecr.io"
            def acrJenkinsCreds = "acr_creds" //this is set in Jenkins global credentials
            sh 'git rev-parse HEAD > git_commit_id.txt'
            try {
                env.GIT_COMMIT_ID = readFile('git_commit_id.txt').trim()
                env.GIT_SHA = env.GIT_COMMIT_ID.substring(0, 7)
            } catch (e) {
                error "${e}"
            }
            def buildName = env.JOB_NAME
            def buildNumber = env.BUILD_NUMBER
            def imageTag = env.BRANCH_NAME + '-' + env.GIT_SHA
            def date = new Date()
            sdf = new SimpleDateFormat("MM/dd/yyyy HH:mm:ss")
            def buildDate = sdf.format(date)
            def appVersion = "${appMajorVersion}.${env.BUILD_NUMBER}"
            def apiImage = "${repo}/smackapi:${imageTag}"
            def webImage = "${repo}/smackweb:${imageTag}"

            // write out variables for debug purposes
            println "DEBUG: env.GIT_COMMIT_ID ==> ${env.GIT_COMMIT_ID}"
            println "DEBUG: env.GIT_SHA ==> ${env.GIT_SHA}"
            println "DEBUG: env.BRANCH_NAME ==> ${env.BRANCH_NAME}"
            println "DEBUG: env.JOB_NAME ==> ${env.JOB_NAME}"
            println "DEBUG: env.BUILD_NUMBER ==> ${env.BUILD_NUMBER}"
            println "DEBUG: appVersion ==> " + appVersion
            println "DEBUG: buildDate ==> " + buildDate
            println "DEBUG: imageTag ==> " + imageTag
            println "DEBUG: apiImage ==> " + apiImage
            println "DEBUG: webImage ==> " + webImage

            println "DEBUG: code compile and test stage starting"
            stage ('BUILD: code compile and test') {
                container('golang') {
                    sh "go get github.com/gorilla/mux"
                    sh "cd smackapi && go build"
                    sh "cd smackapi && go test -v"
                    sh "cd smackweb && go build"
                    sh "cd smackweb && go test -v"
                }
            }

            stage ('BUILD: containerize and publish TO repository') {
                println "DEBUG: build and push containers stage starting"
                container('docker') {
                    // Login to ACR
                    withCredentials([[$class          : 'UsernamePasswordMultiBinding', credentialsId: acrJenkinsCreds,
                                    usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                        println "DEBUG: docker login ${acrServer} -u ${env.USERNAME} -p ${env.PASSWORD}"
                        sh "docker login ${acrServer} -u ${env.USERNAME} -p ${env.PASSWORD}"
                        // optionally push to Docker Hub with a custom Jenkins env variable
                    }

                    // build containers
                    sh "cd smackapi && docker build --build-arg BUILD_DATE='${buildDate}' --build-arg VERSION=${appVersion} --build-arg VCS_REF=${env.GIT_SHA} -t ${apiImage} ."                    
                    sh "cd smackweb && docker build --build-arg BUILD_DATE='${buildDate}' --build-arg VERSION=${appVersion} --build-arg VCS_REF=${env.GIT_SHA} -t ${webImage} ."

                    // push images to repo (ACR)
                    def apiACRImage = acrServer + "/" + apiImage
                    env.ENV_API_IMAGE = "${apiACRImage}"
                    sh "docker tag ${apiImage} ${apiACRImage}"
                    sh "docker push ${apiACRImage}"
                    println "DEBUG: pushed image ${apiACRImage}"
                    def webACRImage = acrServer + "/" + webImage
                    env.ENV_WEB_IMAGE = "${webACRImage}"
                    sh "docker tag ${webImage} ${webACRImage}"
                    sh "docker push ${webACRImage}"
                    println "DEBUG: pushed image ${webACRImage}"
                }
            }

            // use kubernetes plug-in to release or update app
            stage ('DEPLOY: update application on kubernetes') {
                println "DEBUG: deploy new containers to kubernetes stage"
                container('kubectl') {
                    sh "kubectl apply -f kube-jenkins.yaml"
                    sh "kubectl set image deployment/smackapi-deploy smackapi=${env.ENV_API_IMAGE} --namespace=default"
                    sh "kubectl set image deployment/smackweb-deploy smackweb=${env.ENV_WEB_IMAGE} --namespace=default"
                }
            }
        }
    }
