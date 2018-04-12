#!/usr/bin/env groovy

/*
 * Copyright 2018 Dennis Kuhnert
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
