pipeline {
    agent any
    stages {
        stage('Abort Previous Builds') {
            steps {
                sh '''
                rm -rf Compile
                rm -rf data
                rm -rf microservice
                rm -rf restler_bin
                rm -rf restler-fuzzer
                rm -rf Test
                mkdir data
                mkdir restler_bin
                go clean
                '''
            }
        }
        stage('Build') {
            steps {
                sh '''
                go build
                '''
            }
        }
        stage('Curl Test') {
            options { timeout(time: 20, unit: "SECONDS") }
            parallel {
                stage('microservice') {
                    steps {
                        sh '''
                        ./microservice -s &
                        PID=$! && sleep 10 && kill -9 $PID
                        '''
                    }
                }
                stage('curl') {
                    steps {
                        sh '''
                        sleep 2
                        R1=$(curl -s -k https://localhost/v1/get?name=Bob | tr -d '"')
                        R2=$(curl -s -k -X POST https://localhost/v1/set?name=Bob | tr -d '"')
                        R3=$(curl -s -k https://localhost/v1/get?name=Bob | tr -d '"')
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
        stage('Robot Test') {
            options { timeout(time: 20, unit: "SECONDS") }
            parallel {
                stage('microservice') {
                    steps {
                        sh '''
                        ./microservice -s &
                        PID=$! && sleep 10 && kill -9 $PID
                        '''
                    }
                }
                stage('robot') {
                    steps {
                        sh '''
                        sleep 2
                        REQUESTS_CA_BUNDLE=./root_server_cert.pem robot test.robot
                        '''
                    }
                }
            }
        }
        stage('RESTler Test') {
            options { timeout(time: 60, unit: "SECONDS") }
            parallel {
                stage('microservice') {
                    steps {
                        sh '''
                        ./microservice &
                        PID=$! && sleep 50 && kill -9 $PID
                        '''
                    }
                }
                stage('restler') {
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                            export CWD=$(pwd)
                            apt-get update && apt-get install -y dotnet-sdk-6.0
                            git clone https://github.com/microsoft/restler-fuzzer.git
                            python ./restler-fuzzer/build-restler.py --dest_dir "${CWD}/restler_bin"
                            ./restler_bin/restler/Restler compile --api_spec openapi.yaml
                            ./restler_bin/restler/Restler test \
                                --grammar_file Compile/grammar.py \
                                --dictionary_file Compile/dict.json \
                                --settings Compile/engine_settings.json \
                                --target_ip 127.0.0.1 \
                                --target_port 80 \
                                --no_ssl
                            '''
                        }
                    }
                }
            }
        }
    }
}
