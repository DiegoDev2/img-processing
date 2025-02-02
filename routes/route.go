package routes

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/DiegoDev2/img-processing/middleware"
	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

var allowedExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

func RegisterRoutes(g *gin.Engine) {
	g.Use(middleware.RateLimiter)
	g.Use(middleware.CORSMiddleware())
	g.MaxMultipartMemory = 50 << 20

	g.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			log.Println("Error retrieving file:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		ext := filepath.Ext(file.Filename)
		if !allowedExtensions[ext] {
			log.Println("Invalid file extension:", ext)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only PNG, JPG, and WEBP are allowed"})
			return
		}

		src, err := file.Open()
		if err != nil {
			log.Println("Error opening file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		buf := make([]byte, 512)
		_, err = src.Read(buf)
		if err != nil {
			log.Println("Error reading file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
		log.Printf("File header: %x\n", buf)

		src.Seek(0, 0)

		var img image.Image
		switch ext {
		case ".png":
			img, err = png.Decode(src)
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(src)
		case ".webp":
			img, err = webp.Decode(src)
		default:
			log.Println("Unsupported file extension:", ext)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
			return
		}

		if err != nil {
			log.Println("Error decoding image:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
			return
		}

		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		if width > 2560 || height > 2048 {
			log.Println("Resizing image to fit within 2560x2048")
			img = resize.Resize(2560, 2048, img, resize.Lanczos3)
		}

		if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
			log.Println("Error creating uploads directory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
			return
		}

		webpPath := filepath.Join("./uploads", filepath.Base(file.Filename)+".webp")

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			outFile, err := os.Create(webpPath)
			if err != nil {
				log.Println("Error creating output file:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output file"})
				return
			}
			defer outFile.Close()

			options := &webp.Options{Quality: 80}
			if err := webp.Encode(outFile, img, options); err != nil {
				log.Println("Error encoding WebP:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode WebP"})
				return
			}
		}()

		wg.Wait()

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded, optimized, and converted to WebP successfully",
			"path":    webpPath,
		})
	})
}
