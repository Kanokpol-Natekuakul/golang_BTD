# Go REST API with Gin and GORM

REST API ที่สร้างด้วย Go, Gin Framework และ GORM สำหรับจัดการข้อมูล Items พร้อมฐานข้อมูล MySQL

## ฟีเจอร์

- CRUD Operations (Create, Read, Update, Delete)
- RESTful API Design
- JSON Response Format
- MySQL Database Integration
- Auto Migration
- Error Handling

## เทคโนโลยีที่ใช้

- **Go** 1.24.1 - Programming Language
- **Gin** - HTTP Web Framework
- **GORM** - ORM Library
- **MySQL** - Database
- **JSON** - Data Format

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Welcome message |
| GET | `/items` | ดึงข้อมูล items ทั้งหมด |
| GET | `/items/:id` | ดึงข้อมูล item ตาม ID |
| POST | `/items` | สร้าง item ใหม่ |
| PUT | `/items/:id` | อัปเดต item ตาม ID |
| DELETE | `/items/:id` | ลบ item ตาม ID |

## โครงสร้างฐานข้อมูล

### Table: items

| Column | Type | Description |
|--------|------|-------------|
| id | INT (Primary Key, Auto Increment) | รหัสสินค้า |
| name | VARCHAR(255) | ชื่อสินค้า |
| price | DECIMAL(10,2) | ราคา |
| created_at | TIMESTAMP | วันที่สร้าง |
| updated_at | TIMESTAMP | วันที่อัปเดต |

## การติดตั้งและรัน

### ข้อกำหนดเบื้องต้น

- Go 1.24.1 หรือใหม่กว่า
- MySQL Server
- Git

### ขั้นตอนการติดตั้ง

1. **โคลนโปรเจค**
   ```bash
   git clone <repository-url>
   cd myapi
   ```

2. **ติดตั้ง Dependencies**
   ```bash
   go mod download
   ```

3. **ตั้งค่าฐานข้อมูล MySQL**
   ```bash
   # เข้า MySQL
   mysql -u root -p
   
   # รันสคริปต์ฐานข้อมูล
   source database.sql
   ```

4. **แก้ไขการเชื่อมต่อฐานข้อมูล** (ถ้าจำเป็น)
   
   เปิดไฟล์ `main.go` และแก้ไข DSN connection string:
   ```go
   dsn := "root:your_password@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
   ```

5. **รันแอปพลิเคชัน**
   ```bash
   go run main.go
   ```

   เซิร์ฟเวอร์จะรันที่ `http://localhost:8080`

## การใช้งาน API

### 1. ดึงข้อมูล Items ทั้งหมด
```bash
curl -X GET http://localhost:8080/items
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Laptop",
    "price": 29999.99,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

### 2. ดึงข้อมูล Item ตาม ID
```bash
curl -X GET http://localhost:8080/items/1
```

### 3. สร้าง Item ใหม่
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "price": 1999.50
  }'
```

### 4. อัปเดต Item
```bash
curl -X PUT http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "price": 2499.99
  }'
```

### 5. ลบ Item
```bash
curl -X DELETE http://localhost:8080/items/1
```

## โครงสร้างโปรเจค

```
myapi/
├── main.go           # ไฟล์หลักของแอปพลิเคชัน
├── database.sql      # สคริปต์สร้างฐานข้อมูล
├── go.mod           # Go modules configuration
├── go.sum           # Go modules checksums
└── README.md        # ไฟล์นี้
```

## รายละเอียดของโค้ด

### Models
```go
type Item struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Database Connection
```go
func ConnectDatabase() {
    dsn := "root:123456789@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }
    database.AutoMigrate(&Item{})
    DB = database
}
```

## การพัฒนาต่อ

### เพิ่ม Middleware
```go
// CORS Middleware
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
    c.Header("Access-Control-Allow-Headers", "Content-Type")
    
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    
    c.Next()
})
```

### เพิ่ม Validation
```go
type CreateItemRequest struct {
    Name  string  `json:"name" binding:"required,min=1,max=255"`
    Price float64 `json:"price" binding:"required,gt=0"`
}
```

### เพิ่ม Pagination
```go
func GetAllItems(c *gin.Context) {
    var items []Item
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    offset := (page - 1) * limit
    
    DB.Offset(offset).Limit(limit).Find(&items)
    c.JSON(http.StatusOK, items)
}
```

## การทดสอบ

### ใช้ Postman
1. Import collection หรือสร้าง requests ใหม่
2. ตั้งค่า Base URL: `http://localhost:8080`
3. ทดสอบแต่ละ endpoint

### ใช้ curl commands
ใช้ตัวอย่าง curl commands ที่ให้ไว้ข้างต้น

## การแก้ปัญหา

### ปัญหาการเชื่อมต่อฐานข้อมูล
1. ตรวจสอบว่า MySQL server รันอยู่
2. ตรวจสอบ username/password ใน DSN
3. ตรวจสอบว่าฐานข้อมูล `mydatabase` ถูกสร้างแล้ว

### ปัญหา Port ซ้ำ  
```bash
# หา process ที่ใช้ port 8080
lsof -i :8080

# ฆ่า process
kill -9 <PID>
```

### ปัญหา Go Modules
```bash
# ล้าง module cache
go clean -modcache

# ดาวน์โหลด dependencies ใหม่
go mod download
```

## Dependencies

ดูรายละเอียดใน `go.mod`:

- gin-gonic/gin v1.10.0
- gorm.io/gorm v1.25.12  
- gorm.io/driver/mysql v1.5.7
- go-sql-driver/mysql v1.9.0

## การ Deploy

### Docker (แนะนำ)
```dockerfile
FROM golang:1.24.1-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
```

### Build Binary
```bash
# Build สำหรับ Linux
GOOS=linux GOARCH=amd64 go build -o api-linux main.go

# Build สำหรับ Windows  
GOOS=windows GOARCH=amd64 go build -o api-windows.exe main.go
```

