pipeline {
    agent any

    environment {
        BACKEND_DIR = 'backend'
        FRONTEND_DIR = 'frontend'
        COMPOSE_FILE = 'docker-compose.yml'
    }

    stages {
         stage('Build Backend') {
            steps {
                dir("${BACKEND_DIR}") {
                    sh 'go mod tidy'
                    sh 'go build -o main .'
                }
            }
        }

        stage('Build Frontend') {
            steps {
                dir("${FRONTEND_DIR}") {
                    sh 'npm install'
                    sh 'npm run build'
                }
            }
        }

        stage('Docker Compose Build') {
            steps {
                sh 'docker-compose build'
            }
        }

        stage('Docker Compose Up') {
            steps {
                sh 'docker-compose up -d'
            }
        }

        stage('Health Check') {
            steps {
                script {
                    sleep 5 // Give services time to boot
                    def response = sh(script: "curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/api/todos", returnStdout: true).trim()
                    if (response != "200") {
                        error("Health check failed with response code ${response}")
                    }
                }
            }
        }
    }

    post {
        always {
            sh 'docker-compose down'
        }
    }
}
