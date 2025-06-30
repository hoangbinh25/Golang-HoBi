# Golang-HoBi

Một ứng dụng e-commerce được xây dựng bằng Go, hỗ trợ đăng ký/đăng nhập, quản lý sản phẩm, giỏ hàng, đặt hàng và thanh toán.

## 🚀 Features

- **Authentication**: Đăng ký, đăng nhập với email/password và Google OAuth
- **Email Verification**: Xác thực email khi đăng ký
- **Password Reset**: Quên mật khẩu và reset qua email
- **Product Management**: CRUD sản phẩm, danh mục
- **Shopping Cart**: Thêm, sửa, xóa sản phẩm trong giỏ hàng
- **Order Management**: Đặt hàng, theo dõi trạng thái đơn hàng
- **Admin Panel**: Quản lý sản phẩm, đơn hàng, xác nhận giao hàng
- **Responsive Design**: Giao diện responsive cho mobile và desktop

## 🛠 Tech Stack

- **Backend**: Go (Golang)
- **Database**: MySQL (Aiven)
- **Authentication**: Google OAuth 2.0
- **Email**: SMTP
- **Frontend**: HTML, CSS, JavaScript
- **Deployment**: Docker, Render
- **Session Management**: Gorilla Sessions

## 📋 Prerequisites

- Go 1.23+
- MySQL 8.0+
- Git

## 🚀 Quick Start

### Local Development

1. **Clone repository**

```bash
git clone https://github.com/your-username/Golang-HoBi.git
cd Golang-HoBi
```

2. **Setup environment variables**

```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Install dependencies**

```bash
go mod download
```

4. **Setup database**

```bash
# Import db.sql vào MySQL database
mysql -u root -p go_ecommerce < db.sql
```

5. **Run application**

```bash
go run main.go
```

App sẽ chạy tại: `http://localhost:8080`

### Environment Variables

Tạo file `.env` với các biến sau:

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=go_ecommerce

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Session
SESSION_SECRET=your-session-secret

# Server
PORT=8080
```

## 🚀 Deployment

Xem hướng dẫn chi tiết tại [DEPLOYMENT.md](./DEPLOYMENT.md)

### Deploy lên Render với Aiven Database

1. **Setup Aiven MySQL Database**
2. **Configure Google OAuth**
3. **Deploy lên Render**
4. **Import database schema**

## 📁 Project Structure

```
Golang-HoBi/
├── assets/                 # Static files (CSS, JS, images)
├── controllers/           # HTTP handlers
│   ├── adminController/   # Admin panel controllers
│   ├── cartController/    # Shopping cart controllers
│   ├── checkoutController/ # Checkout controllers
│   ├── userController/    # User management controllers
│   └── oauth/            # OAuth authentication
├── initializers/         # Database, session, OAuth setup
├── models/              # Data models
├── views/               # HTML templates
├── main.go             # Application entry point
├── Dockerfile          # Docker configuration
├── render.yaml         # Render deployment config
└── db.sql             # Database schema
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Nếu gặp vấn đề, vui lòng:

1. Kiểm tra [Issues](../../issues) để xem có ai đã báo cáo chưa
2. Tạo issue mới với thông tin chi tiết về lỗi
3. Liên hệ qua email: binhhp.work@gmail.com

## 🔄 Changelog

### v1.0.0

- Initial release
- Basic e-commerce functionality
- User authentication
- Product management
- Shopping cart
- Order processing
