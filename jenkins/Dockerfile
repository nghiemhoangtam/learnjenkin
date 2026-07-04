# Jenkins LTS + Go, đóng gói sẵn để pipeline có thể chạy `go build` / `go test`
# mà không cần cài thêm gì.
FROM jenkins/jenkins:lts-jdk21

# Chuyển sang root để cài phần mềm hệ thống.
USER root

# Phiên bản Go muốn cài (khớp với môi trường dev cho nhất quán).
ARG GO_VERSION=1.24.3

# Cài gcc (cần cho `go test -race`) và tải Go về /usr/local/go.
# dpkg --print-architecture trả về "amd64" hoặc "arm64", khớp tên file của Go.
RUN apt-get update \
    && apt-get install -y --no-install-recommends gcc ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && ARCH="$(dpkg --print-architecture)" \
    && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" -o /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

# Đưa Go vào PATH và cấu hình GOPATH nằm trong jenkins_home (được lưu bền vững).
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/var/jenkins_home/go"
ENV GOCACHE="/var/jenkins_home/.cache/go-build"

# Cho phép git thao tác với thư mục mà Jenkins checkout (tránh lỗi "dubious ownership").
RUN git config --system --add safe.directory '*'

# Cài sẵn plugin để bỏ qua bước chọn plugin thủ công trong wizard.
COPY plugins.txt /usr/share/jenkins/ref/plugins.txt
RUN jenkins-plugin-cli --plugin-file /usr/share/jenkins/ref/plugins.txt

# Trả về user jenkins (không chạy bằng root).
USER jenkins
