# Deployment Guide - Golang HoBi

## 1. Setup Aiven MySQL Database

### Bước 1: Tạo Aiven Account

1. Truy cập [Aiven.io](https://aiven.io)
2. Đăng ký tài khoản mới
3. Tạo project mới

### Bước 2: Tạo MySQL Service

1. Trong Aiven Console, click "Create Service"
2. Chọn "MySQL"
3. Chọn plan phù hợp (Startup plan free)
4. Chọn region gần nhất
5. Đặt tên service: `golang-hobi-db`
6. Click "Create Service"

### Bước 3: Lấy thông tin kết nối

1. Vào service MySQL vừa tạo
2. Tab "Overview" → "Connection Information"
3. Ghi lại các thông tin:
   - Host
   - Port
   - Database name
   - Username
   - Password

## 2. Setup Google OAuth

### Bước 1: Tạo Google Cloud Project

1. Truy cập [Google Cloud Console](https://console.cloud.google.com)
2. Tạo project mới
3. Enable Google+ API

### Bước 2: Tạo OAuth Credentials

1. Vào "APIs & Services" → "Credentials"
2. Click "Create Credentials" → "OAuth 2.0 Client IDs"
3. Chọn "Web application"
4. Thêm Authorized redirect URIs:
   - `http://localhost:8080/auth/callback` (development)
   - `https://your-app-name.onrender.com/auth/callback` (production)
5. Ghi lại Client ID và Client Secret

## 3. Deploy lên Render

### Bước 1: Push code lên GitHub

```bash
git add .
git commit -m "Prepare for deployment"
git push origin main
```

### Bước 2: Tạo Render Service

1. Truy cập [Render.com](https://render.com)
2. Đăng ký/đăng nhập với GitHub
3. Click "New" → "Web Service"
4. Connect GitHub repository
5. Chọn repository `Golang-HoBi`

### Bước 3: Cấu hình Render Service

- **Name**: `golang-hobi`
- **Environment**: `Docker`
- **Region**: `Oregon` (hoặc region gần nhất)
- **Branch**: `main`
- **Root Directory**: (để trống)
- **Build Command**: (để trống - sử dụng Dockerfile)
- **Start Command**: (để trống - sử dụng Dockerfile)

### Bước 4: Cấu hình Environment Variables

Trong Render Dashboard, thêm các environment variables:

```
DB_HOST=your-aiven-mysql-host.aivencloud.com
DB_PORT=12345
DB_USER=avnadmin
DB_PASSWORD=your-aiven-mysql-password
DB_NAME=defaultdb
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
SESSION_SECRET=your-random-session-secret
PORT=8080
```

### Bước 5: Deploy

1. Click "Create Web Service"
2. Render sẽ tự động build và deploy
3. Chờ deployment hoàn thành (5-10 phút)

## 4. Database Migration

### Bước 1: Import database schema

1. Vào Aiven Console → MySQL service
2. Tab "Console" → "Query Editor"
3. Copy và paste nội dung file `db.sql`
4. Execute để tạo tables

### Bước 2: Verify connection

1. Truy cập URL của app: `https://your-app-name.onrender.com`
2. Kiểm tra health check: `https://your-app-name.onrender.com/health`

## 5. Troubleshooting

### Lỗi Database Connection

- Kiểm tra environment variables trong Render
- Verify Aiven MySQL service đang running
- Check firewall rules và network access

### Lỗi Build

- Kiểm tra Dockerfile syntax
- Verify go.mod và go.sum files
- Check build logs trong Render

### Lỗi OAuth

- Verify Google OAuth redirect URIs
- Check environment variables GOOGLE_CLIENT_ID và GOOGLE_CLIENT_SECRET

## 6. Monitoring

### Render Logs

- Vào Render Dashboard → Service → Logs
- Monitor application logs và errors

### Aiven Monitoring

- Vào Aiven Console → MySQL service → Metrics
- Monitor database performance và connections

## 7. SSL/HTTPS

- Render tự động cung cấp SSL certificate
- App sẽ accessible qua HTTPS: `https://your-app-name.onrender.com`

## 8. Custom Domain (Optional)

1. Trong Render Dashboard → Service → Settings
2. Tab "Custom Domains"
3. Add custom domain và configure DNS
