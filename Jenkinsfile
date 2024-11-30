pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh '''
                go clean
                go build
                if [ -e data ]; then
                  rm -rf data
                fi
                mkdir data
                '''
            }
        }
        stage('Curl Test') {
            options { timeout(time: 20, unit: "SECONDS") }
            parallel {
                stage('microservice') {
                    steps {
                        sh '''
                            echo "++++ microservice ++++"
                            ./microservice -s &
                            PID=$! && sleep 15 && kill -9 $PID
                        '''
                    }
                }
                stage('curl') {
                    steps {
                        sh '''
                            echo "++++ cURL ++++"
                            sleep 5
                            R1=$(curl -k https://localhost/v1/get?name=Bob | tr -d '"')
                            R2=$(curl -k -X POST https://localhost/v1/set?name=Bob | tr -d '"')
                            R3=$(curl -k https://localhost/v1/get?name=Bob | tr -d '"')
                            if [ "$R1" != "{message:Internal Server Error}" ]; then
                                exit 1
                            fi
                            if [ "$R2" != "welcome Bob" ]; then
                                exit 1
                            fi
                            if [ "$R3" != "hello Bob" ]; then
                                exit 1
                            fi
                        '''
                    }
                }
            }
        }
    }
}