#!/usr/bin/env groovy

def orgDir = 'github.com/kyroy'
def projectDir = "${orgDir}/kdtree"
def goOrgDir = "/go/src/${orgDir}"
def goProjectDir = "/go/src/${projectDir}"

def testLog = 'jenkins_test.log'
def testXML = 'jenkins_test.xml'
def coverLog = 'jenkins_cover.log'
def coverXML = 'jenkins_cover.xml'

pipeline {
    agent { docker 'kyroy/go-test' }
    stages {
        stage('Init') {
            steps {
                sh "mkdir -p ${goOrgDir}"
                sh "ln -s \$PWD ${goProjectDir}"

                sh "cd ${goProjectDir} && go get -t ./..."
            }
        }

        stage('Lint') {
            steps {
                echo "cf ${goProjectDir} && gometalinter --disable-all --enable=vet --enable=golint --vendor ./..."
            }
        }

        stage('Test') {
            steps {
                // test
                sh "cd ${goProjectDir} && 2>&1 go test -coverprofile=${coverLog} -covermode count -v ./... | tee ${testLog}"

                // generate metadata
                sh "go2xunit -fail -input ${testLog} -output ${testXML}"
                sh "cd ${goProjectDir} && gocover-cobertura < ${coverLog} > ${coverXML}"

                step([$class             : 'CoberturaPublisher',
                      autoUpdateHealth   : false,
                      autoUpdateStability: false,
                      coberturaReportFile: coverXML,
                      failUnhealthy      : false,
                      failUnstable       : false,
                      maxNumberOfBuilds  : 0,
                      onlyStable         : false,
                      sourceEncoding     : 'ASCII',
                      zoomCoverageChart  : false])
            }
        }
    }
    post {
        always {
            junit allowEmptyResults: true, testResults: '*test.xml'
            deleteDir()
        }
        changed {
            notifyStatusChange notificationRecipients: 'dennis.kuhnert@sap.com,mail@kyroy.com', componentName: 'kyroy/kdtree'
        }
    }
}
