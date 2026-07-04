// Jenkinsfile - định nghĩa pipeline dạng "Declarative Pipeline".
// Toàn bộ quy trình CI được viết bằng code và lưu cùng source (Pipeline as Code).
pipeline {
    // Chạy trên node Jenkins bất kỳ (Jenkins native trên máy, hoặc container).
    agent any

    environment {
        APP_DIR = 'app'          // Thư mục chứa go.mod
        BIN_NAME = 'app'         // Tên file binary build ra
        CGO_ENABLED = '1'        // Bật CGO để dùng được `go test -race`
        GOROOT = '/usr/local/go' // Nơi cài Go
        // Tiến trình Jenkins native thường thiếu Go trong PATH -> thêm vào đây.
        // Dùng được cho cả Jenkins native (macOS) lẫn container (đều cài Go ở /usr/local/go).
        PATH = "/usr/local/go/bin:${PATH}"
    }

    options {
        timestamps()                                  // Thêm mốc thời gian vào log
        timeout(time: 15, unit: 'MINUTES')            // Tự huỷ nếu chạy quá lâu
        buildDiscarder(logRotator(numToKeepStr: '10')) // Chỉ giữ 10 build gần nhất
    }

    stages {
        stage('Môi trường') {
            steps {
                echo 'Thông tin Jenkins agent:'
                echo "  NODE_NAME       = ${env.NODE_NAME}"
                echo "  NODE_LABELS     = ${env.NODE_LABELS}"
                echo "  EXECUTOR_NUMBER = ${env.EXECUTOR_NUMBER}"
                echo "  WORKSPACE       = ${env.WORKSPACE}"
                echo "  JENKINS_URL     = ${env.JENKINS_URL}"
                sh '''
                    echo "Hostname = $(hostname)"
                    echo "OS       = $(uname -s)"
                    echo "Arch     = $(uname -m)"
                '''
                echo 'Kiểm tra công cụ build...'
                sh 'go version'
                sh 'echo "GOPATH=$GOPATH" && echo "PATH=$PATH"'
            }
        }

        stage('Tải dependencies') {
            steps {
                dir("${APP_DIR}") {
                    sh 'go mod download'
                    sh 'go mod verify'
                }
            }
        }

        stage('Kiểm tra format') {
            steps {
                dir("${APP_DIR}") {
                    // gofmt -l liệt kê các file chưa được format đúng.
                    // Nếu có file nào -> fail để buộc code luôn sạch.
                    sh '''
                        unformatted=$(gofmt -l .)
                        if [ -n "$unformatted" ]; then
                            echo "Các file chưa format đúng:"
                            echo "$unformatted"
                            exit 1
                        fi
                        echo "Tất cả file đã được format đúng."
                    '''
                }
            }
        }

        stage('Vet (phân tích tĩnh)') {
            steps {
                dir("${APP_DIR}") {
                    sh 'go vet ./...'
                }
            }
        }

        stage('Test') {
            steps {
                dir("${APP_DIR}") {
                    // -race: phát hiện race condition; -coverprofile: sinh dữ liệu coverage.
                    sh 'go test -v -race -coverprofile=coverage.out ./...'
                    // In tỉ lệ coverage tổng ra log.
                    sh 'go tool cover -func=coverage.out | tail -n 1'
                    // Sinh báo cáo coverage dạng HTML để lưu lại làm artifact.
                    sh 'go tool cover -html=coverage.out -o coverage.html'
                }
            }
        }

        stage('Build') {
            steps {
                dir("${APP_DIR}") {
                    // Build binary Linux tĩnh vào thư mục ../bin ở gốc workspace.
                    sh 'CGO_ENABLED=0 go build -o ../bin/${BIN_NAME} .'
                    sh 'ls -lh ../bin'
                }
            }
        }
    }

    post {
        always {
            // Lưu binary và báo cáo coverage để tải về từ giao diện Jenkins.
            archiveArtifacts artifacts: 'bin/**, app/coverage.*', allowEmptyArchive: true, fingerprint: true
        }
        success {
            echo 'Pipeline THANH CONG - build va test deu pass.'
        }
        failure {
            echo 'Pipeline THAT BAI - hay xem log stage bi do o tren.'
        }
        cleanup {
            cleanWs() // Dọn workspace sau mỗi lần chạy (cần plugin ws-cleanup).
        }
    }
}
