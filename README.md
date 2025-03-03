## 🛠 Cài đặt môi trường
1. **Visual Studio Code**: [Download tại đây](https://code.visualstudio.com/)
2. **Node.js**: 
   - Cài đặt bản LTS từ [trang chủ Node.js](https://nodejs.org/)
   - ![Node.js Logo](https://upload.wikimedia.org/wikipedia/commons/d/d9/Node.js_logo.svg)
3. **Golang**:
   - Tải bản mới nhất từ [golang.org](https://go.dev/dl/)
   - ![Golang Logo](https://go.dev/images/go-logo-blue.svg)

---

# Tên Dự Án 
Dự án 1

## Mô tả
Xây dựng ứng dụng web Golang dùng Iris framework có một ô text area cho người dùng viết promp
ấn nút gửi đi, ứng dụng gọi vào Groq API và trả về kết quả cho người dùng xem. 

## Cách chạy
- Sử dụng VScode
- Vào mục Backend tạo file .env rồi thêm GROQ_API_KEY = "your api key"
![Demo](./Image/1.png)
![Demo](./Image/6.png)
- Bật terminal thự mục rồi gõ: 
    + ![Demo](./Image/2.png)
    + go mod init project, để giúp quản lý dependencies (thư viện, package) cho dự án
    + go run main.go, để chạy dự án

- Vào mục Frontend mở terminal rồi gõ:
    + ![Demo](./Image/4.png)
    + Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope Process, để cấp quyền nếu bạn dùng window
    + ng serve, để chạy dự án 

- Vào mạng truy cập http://localhost:4200/ để test 
![Demo](./Image/3.png)

# Tên Dự Án 
Dự án 2

## Mô tả
Xây dựng ứng dụng web sinh file SSML từ hội thoại

## Cách chạy
- Sử dụng VScode
- Vào mục Frontend mở terminal rồi gõ:
    + ![Demo](./Image/4.png)
    + Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope Process, để cấp quyền nếu bạn dùng window
    + ng serve, để chạy dự án 

- Vào mạng truy cập http://localhost:4200/ để test 
![Demo](./Image/5.png)

# Tên Dự Án 
Dự án 3

## Mô tả
lập trình một ứng dụng web dùng Iris framework gọi vào Groq API, model deepseek-r1-distill-llama-70b để tự động hóa việc sinh ra 1 cuộc trò chuyện rồi lọc ra các từ quan trọng và dịch ra tiếng anh rồi lưu vào postgresql

## Cách chạy
- Sử dụng VScode
- Vào mục Backend tạo file .env rồi thêm GROQ_API_KEY = "your api key"
![Demo](./Image/1.png)
![Demo](./Image/6.png)
- Vào mục Backend/Postgresql/db.go để sửa user="tên người dùng" password="mật khẩu" dbname="Tên database" theo database ở Postgresql
![Demo](./Image/7.png)
- Bật terminal thự mục rồi gõ: 
    + go mod init project, để giúp quản lý dependencies (thư viện, package) cho dự án
    + go run main.go, để chạy dự án

- Vào mục Frontend mở terminal rồi gõ:
    + ![Demo](./Image/4.png)
    + Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope Process, để cấp quyền nếu bạn dùng window
    + ng serve, để chạy dự án 

- Vào mạng truy cập http://localhost:4200/ để test 
![Demo](./Image/8.png)
![Demo](./Image/9.png)