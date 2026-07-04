# Lab học Jenkins: Build & Test ứng dụng Go

Lab này giúp bạn hiểu Jenkins qua một ví dụ thực tế: dùng Jenkins để **build** và **test** một ứng dụng Go, với toàn bộ quy trình được định nghĩa bằng code (`Jenkinsfile`).

## 1. Kiến trúc lab

```
learnjenkin/
├── app/                      # Ứng dụng Go
│   ├── go.mod
│   ├── main.go               # HTTP server nhỏ (endpoint /health, /add, /divide)
│   └── calculator/
│       ├── calculator.go     # Logic (Add, Subtract, Multiply, Divide)
│       └── calculator_test.go# Unit test
├── jenkins/
│   ├── Dockerfile            # Jenkins LTS + Go + gcc cài sẵn
│   └── plugins.txt           # Danh sách plugin
├── docker-compose.yml        # Chạy Jenkins bằng 1 lệnh
├── Jenkinsfile               # Định nghĩa pipeline CI (build + test)
└── README.md
```

**Ý tưởng:** Chúng ta chạy Jenkins trong Docker. Image Jenkins được build kèm sẵn Go, nên trong pipeline có thể gọi thẳng `go test`, `go build` mà không phải cài thêm gì.

## 2. Yêu cầu

- Docker + Docker Compose (đã có: Docker 28.x, Compose v2.x)
- Go (tùy chọn, chỉ cần nếu muốn chạy test dưới máy local): `go 1.24+`

## 3. Các bước thực hiện

### Bước 0 — (Tùy chọn) Chạy thử app dưới máy local

Trước khi đưa lên Jenkins, hãy chắc chắn code chạy được:

```bash
cd app
go test ./...        # chạy unit test
go run .             # chạy server, mở http://localhost:8090/health
```

### Bước 1 — Khởi động Jenkins

Từ thư mục gốc của lab:

```bash
docker compose up -d --build
```

- Lần đầu sẽ hơi lâu vì phải tải image Jenkins + Go + cài plugin.
- Kiểm tra container đã chạy: `docker compose ps`
- Xem log: `docker compose logs -f jenkins`

### Bước 2 — Mở khoá Jenkins (Unlock)

1. Mở trình duyệt: <http://localhost:8080>
2. Jenkins yêu cầu mật khẩu admin ban đầu. Lấy bằng lệnh:

   ```bash
   docker exec jenkins-lab cat /var/jenkins_home/secrets/initialAdminPassword
   ```

3. Dán mật khẩu vào ô, bấm **Continue**.
4. Chọn **Install suggested plugins** (plugin cần thiết đã cài sẵn nên bước này sẽ nhanh).
5. Tạo tài khoản admin đầu tiên (nhớ user/mật khẩu này).
6. Giữ nguyên Jenkins URL mặc định → **Save and Finish** → **Start using Jenkins**.

### Bước 3 — Kiểm tra Jenkins thấy Go

Vào **Manage Jenkins → Tools** hoặc đơn giản là tạo pipeline test nhanh (Bước 4). Bạn cũng có thể xác nhận Go có trong container:

```bash
docker exec jenkins-lab go version
```

### Bước 4 — Đưa code vào Git

Jenkins pipeline "as code" lấy `Jenkinsfile` + source từ một Git repository. Có **2 lựa chọn**:

**Lựa chọn A — Dùng repo local đã mount sẵn (không cần GitHub, nhanh nhất để học):**

`docker-compose.yml` đã mount thư mục lab vào container tại `/workspace`. Ta chỉ cần biến nó thành một git repo:

```bash
git init
git add .
git commit -m "khởi tạo lab jenkins"
```

Khi tạo job (Bước 5), dùng Repository URL là: `/workspace`

**Lựa chọn B — Dùng GitHub (giống thực tế):**

1. Tạo repo trên GitHub, push code lên.
2. Khi tạo job, dùng URL `https://github.com/<user>/<repo>.git`.

### Bước 5 — Tạo Pipeline job

1. Trang chủ Jenkins → **New Item**.
2. Nhập tên: `go-app-ci` → chọn **Pipeline** → **OK**.
3. Kéo xuống mục **Pipeline**:
   - **Definition**: chọn **Pipeline script from SCM**
   - **SCM**: chọn **Git**
   - **Repository URL**: `/workspace` (Lựa chọn A) hoặc URL GitHub (Lựa chọn B)
   - **Branch Specifier**: đổi thành `*/main` (hoặc `*/master` tùy nhánh của bạn)
   - **Script Path**: giữ nguyên `Jenkinsfile`
4. **Save**.

> Mẹo: Nếu chỉ muốn thử nhanh, chọn **Definition = Pipeline script** rồi dán trực tiếp nội dung `Jenkinsfile`. Nhưng cách này sẽ không tự lấy source Go, nên để build+test thật hãy dùng **from SCM** như trên.

### Bước 6 — Chạy build

1. Vào job `go-app-ci` → bấm **Build Now**.
2. Xem tiến trình ở **Stage View** (biểu đồ các stage) hoặc bấm vào số build → **Console Output** để xem log chi tiết.
3. Khi thành công, mục **Build Artifacts** sẽ có `bin/app` (binary) và `app/coverage.html` (báo cáo coverage) để tải về.

## 4. Giải thích Jenkinsfile

Pipeline gồm các stage tuần tự:

| Stage | Việc làm | Lệnh chính |
|-------|----------|-----------|
| Môi trường | In phiên bản Go, biến môi trường | `go version` |
| Tải dependencies | Tải & kiểm tra module | `go mod download`, `go mod verify` |
| Kiểm tra format | Bắt lỗi code chưa `gofmt` | `gofmt -l .` |
| Vet | Phân tích tĩnh tìm lỗi tiềm ẩn | `go vet ./...` |
| Test | Chạy unit test + đo coverage | `go test -race -coverprofile` |
| Build | Biên dịch ra binary | `go build -o ../bin/app .` |

Khối `post` chạy sau cùng: lưu artifact (`archiveArtifacts`), báo kết quả và dọn dẹp workspace.

## 5. Thử nghiệm để hiểu sâu hơn

Sau khi build xanh lần đầu, hãy thử "phá" để thấy Jenkins phản ứng:

- **Làm fail test:** sửa `calculator.go` cho `Add` trả về `a + b + 1`, commit, build lại → stage **Test** đỏ.
- **Làm fail format:** thêm khoảng trắng/thụt lề sai vào file `.go`, commit → stage **Kiểm tra format** đỏ.
- **Tự động build khi có commit:** trong cấu hình job, bật **Poll SCM** với lịch `H/2 * * * *` (2 phút kiểm tra 1 lần), hoặc dùng webhook nếu chạy GitHub.

## 6. Lệnh vận hành thường dùng

```bash
docker compose up -d --build      # khởi động / rebuild Jenkins
docker compose ps                 # xem trạng thái
docker compose logs -f jenkins    # xem log
docker compose down               # dừng (giữ dữ liệu trong volume)
docker compose down -v            # dừng và XOÁ sạch dữ liệu Jenkins
docker exec jenkins-lab go version
```

## 7. Dọn dẹp

```bash
docker compose down -v   # xoá container + volume jenkins_home
```
