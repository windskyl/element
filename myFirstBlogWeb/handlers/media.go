package handlers

import (
	"io"
	"log"
	"mime/multipart"
	"myFirstBlogWeb/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// helper to save uploaded file and return saved filename
func saveUploadedFile(fileHeader *multipart.FileHeader, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".png"
	}
	name := utils.GenerateUUID() + ext
	outPath := filepath.Join(destDir, name)
	out, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}
	return name, nil
}

// UploadImage 接收 multipart/form-data 字段名 file，保存到 ./public/uploads 并返回访问 URL
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传的文件"})
		return
	}
	uploadsDir := filepath.Join("public", "uploads")
	savedName, err := saveUploadedFile(file, uploadsDir)
	if err != nil {
		log.Printf("UploadImage save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	// 返回可直接访问的 URL
	url := "/uploads/" + savedName
	c.JSON(http.StatusOK, gin.H{"url": url, "filename": savedName})
}

// DeleteImage 删除指定上传文件（需谨慎）
func DeleteImage(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名缺失"})
		return
	}
	path := filepath.Join("public", "uploads", filepath.Clean(name))
	if err := os.Remove(path); err != nil {
		log.Printf("DeleteImage error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": name})
}
