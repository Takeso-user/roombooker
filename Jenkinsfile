pipeline {
    agent any

    environment {
        GO_VERSION = '1.22'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Lint') {
            steps {
                sh 'go vet ./...'
                sh 'go fmt ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o bin/server ./cmd/server'
            }
        }

        stage('Unit Tests') {
            steps {
                sh 'go test -v -race -coverprofile=coverage.out ./...'
            }
        }

        stage('Integration Tests') {
            steps {
                sh 'go test -v -tags=integration ./...'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh 'docker build -t roombooker:latest .'
            }
        }

        stage('Push Docker Image') {
            steps {
                sh 'docker tag roombooker:latest localhost:5000/roombooker:latest'
                sh 'docker push localhost:5000/roombooker:latest'
            }
        }

        stage('Deploy') {
            steps {
                sh 'docker-compose down'
                sh 'docker-compose up -d'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
            junit '**/test-results.xml'
        }
    }
}
