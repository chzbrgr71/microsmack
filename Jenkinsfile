#!/usr/bin/groovy
import java.text.SimpleDateFormat

podTemplate(label: 'jenkins-pipeline', containers: [
    containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:2.62', args: '${computer.jnlpmac} ${computer.name}', workingDir: '/home/jenkins', resourceRequestCpu: '200m', resourceLimitCpu: '200m', resourceRequestMemory: '256Mi', resourceLimitMemory: '256Mi'),
    containerTemplate(name: 'golang', image: 'golang:1.7.5', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'docker', image: 'docker:17.06.0', command: 'cat', ttyEnabled: true)
],
volumes:[
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock')
])
    {
        node ('jenkins-pipeline') {
            println "DEBUG: Pipeline starting"
            def pwd = pwd()
            checkout scm

            // prep config values for pipeline
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
            println "DEBUG: env.GIT_COMMIT_ID ==> ${env.GIT_COMMIT_ID}"
            println "DEBUG: env.GIT_SHA ==> ${env.GIT_SHA}"
            println "DEBUG: env.BRANCH_NAME ==> ${env.BRANCH_NAME}"
            println "DEBUG: env.JOB_NAME ==> ${env.JOB_NAME}"
            println "DEBUG: env.BUILD_NUMBER ==> ${env.BUILD_NUMBER}"
            println "DEBUG: imageTag ==> " + imageTag
            println "DEBUG: buildDate ==> " + buildDate

            println "DEBUG: Start code compile stage"
            stage ('BUILD: code compile and test') {
                container('golang') {
                    sh "go get github.com/gorilla/mux"
                    sh "cd smackapi && go build"
                    sh "cd smackapi && go test -v"
                    sh "cd smackweb && go build"
                    sh "cd smackweb && go test -v"
                }
            }

            stage ('BUILD: containerize and publish TO ACR') {
                container('docker') {
                    // for now, push to Docker Hub. Set in "Manage Jenkins, Configure System, Environment Variables"
                    sh "docker login -u chzbrgr71 -p ${DOCKER_PWD}"
                    sh "cd smackapi && docker build --build-arg BUILD_DATE='${buildDate}' --build-arg VERSION=1.0.${env.BUILD_NUMBER} --build-arg VCS_REF=${env.GIT_SHA} -t chzbrgr71/smackapi:${imageTag} -f ./smackapi/Dockerfile ."
                    sh "docker push chzbrgr71/smackapi:${imageTag}"
                    sh "cd smackweb && docker build --build-arg BUILD_DATE='${buildDate}' --build-arg VERSION=1.0.${env.BUILD_NUMBER} --build-arg VCS_REF=${env.GIT_SHA} -t chzbrgr71/smackweb:${imageTag} -f ./smackweb/Dockerfile ."
                    sh "docker push chzbrgr71/smackweb:${imageTag}"
                }
            }
        }
    }

// --------------
// Utility functions. These would normally be in an external library in a seperate repo

def kubectlTest() {
    // Test that kubectl can correctly communication with the Kubernetes API
    println "checking kubectl connnectivity to the API"
    sh "kubectl get nodes"

}

def helmLint(String chart_dir) {
    // lint helm chart
    println "running helm lint ${chart_dir}"
    sh "helm lint ${chart_dir}"

}

def helmConfig() {
    //setup helm connectivity to Kubernetes API and Tiller
    println "DEBUG: initiliazing helm client"
    sh "helm init"
    println "DEBUG: checking client/server version"
    sh "helm version"
}


def helmDeploy(Map args) {
    // configure helm client and confirm tiller process is installed
    helmConfig()

    def String namespace

    // if namespace isn't parsed into the function set the namespace to the name
    if (args.namespace == null) {
        namespace = args.name
    } else {
        namespace = args.namespace
    }

    if (args.dry_run) {
        println "DEBUG: Running dry-run deployment"

        //sh "helm upgrade --dry-run --install ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory},ingress.hostname=${args.hostname} --namespace=${namespace}"
        sh "helm upgrade --dry-run --install ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory},ingress.hostname=${args.hostname}"
    } else {
        println "DEBUG: Running deployment"

        // reimplement --wait once it works reliable
        //sh "helm upgrade --install ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory},ingress.hostname=${args.hostname} --namespace=${namespace}"
        sh "helm upgrade --install ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory},ingress.hostname=${args.hostname}"

        // sleeping until --wait works reliably
        sleep(20)

        echo "DEBUG: Application ${args.name} successfully deployed. Use helm status ${args.name} to check"
    }


}

def helmDelete(Map args) {
        println "DEBUG: Running helm delete ${args.name}"

        sh "helm delete ${args.name}"
}

def helmTest(Map args) {
    println "DEBUG: Running: helm test ${args.name} --cleanup"

    sh "helm test ${args.name} --cleanup"
}

def containerBuildPub(Map args) {

    println "Running Docker build/publish: ${args.host}/${args.acct}/${args.repo}:${args.tags}"
      
    docker.withRegistry("https://${args.host}", "${args.auth_id}") {

        def img = docker.image("${args.acct}/${args.repo}")
        sh "docker version"
        sh "docker build --build-arg VCS_REF=${env.GIT_SHA} -t ${args.acct}/${args.repo} ${args.dockerfile}"        
        //for (int i = 0; i < args.tags.size(); i++) {
            //img.push(args.tags.get(i))
        //}
        img.push(args.tags.get(0))
    
        return img.id
        }
    }

def getContainerRepoAcct(config) {

    println "setting container registry creds according to Jenkinsfile.json"
    def String acct

    if (env.BRANCH_NAME == 'master') {
        acct = config.container_repo.master_acct
    } else {
        acct = config.container_repo.alt_acct
    }

    return acct
}

@NonCPS
def getMapValues(Map map=[:]) {
    // jenkins and workflow restriction force this function instead of map.values(): https://issues.jenkins-ci.org/browse/JENKINS-27421
    def entries = []
    def map_values = []

    entries.addAll(map.entrySet())

    for (int i=0; i < entries.size(); i++){
        String value =  entries.get(i).value
        map_values.add(value)
    }

    return map_values
}  
