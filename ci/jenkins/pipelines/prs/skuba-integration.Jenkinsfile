/**
 * This pipeline verifies basic skuba deployment, bootstrapping, and adding nodes to a cluster on GitHub Pr
 */

void setBuildStatus(String context, String description, String state) {
    def body = "{\"state\": \"${state}\", " +
               "\"target_url\": \"${BUILD_URL}/display/redirect\", " +
               "\"description\": \"${description}\", " +
               "\"context\": \"${context}\"}"
    def headers = '-H "Content-Type: application/json" -H "Accept: application/vnd.github.v3+json"'
    def url = "https://${GITHUB_TOKEN}@api.github.com/repos/SUSE/skuba/statuses/${GIT_COMMIT}"

    sh(script: "curl -X POST ${headers} ${url} -d '${body}'", label: "Sending commit status")
}

pipeline {
    agent { node { label 'caasp-team-private' } }

    environment {
        OPENRC = credentials('ecp-openrc')
        GITHUB_TOKEN = credentials('github-token')
        PLATFORM = 'openstack'
    }

    stages {
        stage('Setting GitHub in-progress status') { steps {
            setBuildStatus('jenkins/skuba-integration', 'in-progress', 'pending')
        } }

        stage('Git Clone') { steps {
            deleteDir()
            checkout([$class: 'GitSCM',
                      branches: [[name: "*/${BRANCH_NAME}"], [name: '*/master']],
                      doGenerateSubmoduleConfigurations: false,
                      extensions: [[$class: 'LocalBranch'],
                                   [$class: 'WipeWorkspace'],
                                   [$class: 'RelativeTargetDirectory', relativeTargetDir: 'skuba']],
                      submoduleCfg: [],
                      userRemoteConfigs: [[refspec: '+refs/pull/*/head:refs/remotes/origin/PR-*',
                                           credentialsId: 'github-token',
                                           url: 'https://github.com/SUSE/skuba']]])

            dir("${WORKSPACE}/skuba") {
                sh(script: "git checkout ${BRANCH_NAME}", label: "Checkout PR Branch")
            }
        }}

        stage('Getting Ready For Cluster Deployment') { steps {
            sh(script: 'make -f skuba/ci/Makefile pre_deployment', label: 'Pre Deployment')
            sh(script: 'make -f skuba/ci/Makefile pr_checks', label: 'PR Checks')
        } }

        stage('Cluster Deployment') { steps {
            sh(script: 'make -f skuba/ci/Makefile deploy', label: 'Deploy')
            archiveArtifacts("skuba/ci/infra/${PLATFORM}/terraform.tfstate")
            archiveArtifacts("skuba/ci/infra/${PLATFORM}/terraform.tfvars")
        } }

        stage('Run end-to-end tests') { steps {
           dir("skuba") {
             sh(script: 'make build-ginkgo', label: 'build ginkgo binary')
             sh(script: "make setup-ssh && SKUBA_BIN_PATH=\"${WORKSPACE}/go/bin/skuba\" GINKGO_BIN_PATH=\"${WORKSPACE}/skuba/ginkgo\" IP_FROM_TF_STATE=TRUE PLATFORM=openstack make test-e2e", label: "End-to-end tests")
       } } }
    }
    post {
        always {
            sh(script: 'make --keep-going -f skuba/ci/Makefile post_run', label: 'Post Run')
        }
        cleanup {
            dir("${WORKSPACE}") {
                deleteDir()
            }
        }
        failure {
            setBuildStatus('jenkins/skuba-integration', 'failed', 'failure')
        }
        success {
            setBuildStatus('jenkins/skuba-integration', 'success', 'success')
        }
    }
}