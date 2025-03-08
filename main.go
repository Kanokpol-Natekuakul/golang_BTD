package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Item struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ConnectDatabase() {
	dsn := "root:123456789@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(&Item{})
	DB = database
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) // แก้ปัญหา Proxy Trusted

	// Route เริ่มต้น
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API!"})
	})

	// CRUD Routes
	r.GET("/items", GetAllItems)
	r.GET("/items/:id", GetItemByID)
	r.POST("/items", CreateItem)
	r.PUT("/items/:id", UpdateItem)
	r.DELETE("/items/:id", DeleteItem)

	ConnectDatabase()
	r.Run(":8080")
}

// CRUD Functions
func GetAllItems(c *gin.Context) {
	var items []Item
	DB.Find(&items)
	c.JSON(http.StatusOK, items)
}

func GetItemByID(c *gin.Context) {
	id := c.Param("id")
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateItem(c *gin.Context) {
	var newItem Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&newItem)
	c.JSON(http.StatusCreated, newItem)
}

func UpdateItem(c *gin.Context) {
	var item Item

	// ค้นหาข้อมูลในฐานข้อมูลตาม ID
	if err := DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// รับ JSON ที่ส่งมาแล้วอัปเดตค่า
	var updatedData Item
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Name = updatedData.Name
	item.Price = updatedData.Price
	DB.Save(&item)

	c.JSON(http.StatusOK, item)
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
