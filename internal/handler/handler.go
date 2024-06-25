package handler

import (
	"fadingrose/app-analyzer/internal/logger"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tuotoo/qrcode"
)

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	filePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Run analysis in a new goroutine
	report := AnalyzeFile(filePath, c)

	c.JSON(http.StatusOK, gin.H{"message": report})
}

func UploadURLHandler(c *gin.Context) {
	var req string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL"})
		return
	}

	parsedURL, err := url.ParseRequestURI(req)
	logger.Info.Println("URL:", parsedURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL"})
		return
	}

	response, err := http.Get(parsedURL.String())
	if err != nil || response.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to download URL"})
		return
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read URL content"})
		return
	}

	filePath := filepath.Join("uploads", filepath.Base(parsedURL.Path))
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save URL content"})
		return
	}

	go AnalyzeFile(filePath, c)
	c.JSON(http.StatusOK, gin.H{"message": "URL content downloaded successfully, analysis started."})
}

func UploadQRCodeHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file"})
		return
	}

	filePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save QR code"})
		return
	}

	go func() {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open QR code"})
			return
		}
		defer file.Close()

		qrCode, err := qrcode.Decode(file)
		if err != nil {
			fmt.Printf("Error decoding QR code: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to decode QR code"})
			return
		}

		parsedURL, err := url.ParseRequestURI(qrCode.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL in QR code"})
			return
		}

		response, err := http.Get(parsedURL.String())
		if err != nil || response.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to download URL from QR code"})
			return
		}
		defer response.Body.Close()

		content, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read URL content from QR code"})
			return
		}

		filePath := filepath.Join("uploads", filepath.Base(parsedURL.Path))
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save URL content from QR code"})
			return
		}

		go AnalyzeFile(filePath, c)
		c.JSON(http.StatusOK, gin.H{"message": "QR code content downloaded successfully, analysis started."})
	}()
}

func AnalyzeFile(filePath string, c *gin.Context) string {
	// Simulate analysis delay
	logger.Info.Println("Analyzing file:", filePath)

	// move /uploads to /cache
	// mkdir /cache

	if err := os.Rename("uploads", "cache"); err != nil {
		logger.Fatal.Println("Failed to move uploads to cache:", err)
		return ""
	}
	appName := strings.Split(filePath, "/")[1]
	appPath := filepath.Join("cache", appName)
	appPath = "./" + appPath
	logger.Info.Println("App path:", appPath)

	apk2urlHandler(appPath)

	manifestHandler(appPath)

	inferHandler()

	report := reportHandler()
	return report
}

func apk2urlHandler(appPath string) {
	script := fmt.Sprintf("./apk2url.sh %s", appPath)
	logger.Info.Println("Running apk2url script:", script)
	cmd := exec.Command("sh", "-c", script)
	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to run apk2url script:", err)
		panic(err)
	}

	// appPath ./cache/sdjsq_58419.apk
	// move ./endpoints to ./cache/endpoints
	cmd = exec.Command("mv", "./endpoints", "./cache/endpoints")
	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to move endpoints:", err)
		panic(err)
	}

	appName := strings.Split(appPath, "/")[2]
	appName = strings.Split(appName, ".")[0]
	// move {appName}-decompiled to ./cache/{appName}-decompiled
	cmd = exec.Command("mv", fmt.Sprintf("%s-decompiled", appName), fmt.Sprintf("cache/%s-decompiled", appName))
	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to move decompiled:", err)
		panic(err)
	}

}

func manifestHandler(appPath string) {
	appName := strings.Split(appPath, "/")[2]
	appName = strings.Split(appName, ".")[0]
	// move AndroidManifest.xml to ./cache/AndroidManifest.xml
	script := "./manifest.py"
	manifestPath := fmt.Sprintf("./cache/%s-decompiled/%s_jadx/resources/AndroidManifest.xml", appName, appName)
	endpointsPath := fmt.Sprintf("./cache/endpoints/%s_endpoints.txt", appName)
	uniPath := fmt.Sprintf("./cache/endpoints/%s_uniqurls.txt", appName)

	logger.Info.Println("Running manifest script:", script, manifestPath, endpointsPath, uniPath)

	cmd := exec.Command("python", script, manifestPath, endpointsPath, uniPath)

	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to run manifest script:", err)
		panic(err)
	}
}

func inferHandler() {

	logger.Info.Println("Running inference scripts")
	target := "../models/target.json"
	outputDir := "./cache"

	// mv ./target.json to ../models/target.json
	cmd := exec.Command("mv", "./target.json", "../models/target.json")
	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to move target.json:", err)
		panic(err)
	}

	binary_script := "../models/infer_binary.py"
	binary_model := "../models/model.bin.keras"

	// run binary script
	cmd = exec.Command("python", binary_script, binary_model, target, outputDir)

	// var stdout, stderr bytes.Buffer

	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to run binary script:", err)
		panic(err)
	}

	multi_script := "../models/infer_multi.py"
	multi_model := "../models/model.multi.keras"

	// run multi script

	cmd = exec.Command("python", multi_script, multi_model, target, outputDir)

	logger.Info.Println("Running multi script:", multi_script, multi_model, target, outputDir)

	if err := cmd.Run(); err != nil {
		logger.Fatal.Println("Failed to run multi script:", err)
		panic(err)
	}

}

type Report struct {
	Binary_White     float32
	Binary_Malicious float32
	Multi_Black      float32
	Multi_Gamble     float32
	Multi_Scam       float32
	Multi_Sex        float32
}

func (r Report) String() string {
	// convert float32 to percentage

	return fmt.Sprintf("Binary: %.2f%% white, %.2f%% malicious\nMulti: %.2f%% black, %.2f%% gamble, %.2f%% scam, %0.f%% sex", r.Binary_White, r.Binary_Malicious, r.Multi_Black, r.Multi_Gamble, r.Multi_Scam, r.Multi_Sex)

}

func reportHandler() string {
	logger.Info.Println("Generating report")
	binary_path := "./cache/predictions_binary.txt"
	multi_path := "./cache/predictions_multi.txt"

	binary, err := os.ReadFile(binary_path)
	if err != nil {
		logger.Fatal.Println("Failed to read binary predictions:", err)
		panic(err)
	}

	multi, err := os.ReadFile(multi_path)
	if err != nil {
		logger.Fatal.Println("Failed to read multi predictions:", err)
		panic(err)
	}

	binary_report := strings.Split(string(binary), " ")
	multi_report := strings.Split(string(multi), " ")

	report := Report{
		Binary_White:     to_float32(binary_report[0]) * 100,
		Binary_Malicious: to_float32(binary_report[1]) * 100,
		Multi_Black:      to_float32(multi_report[0]) * 100,
		Multi_Gamble:     to_float32(multi_report[1]) * 100,
		Multi_Scam:       to_float32(multi_report[2]) * 100,
		Multi_Sex:        to_float32(multi_report[3]) * 100,
	}

	return report.String()
}

func to_float32(s string) float32 {
	// delete '\n'
	s = strings.ReplaceAll(s, "\n", "")
	float64Val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		logger.Fatal.Println("Failed to convert string to float32:", err)
		panic(err)
	}
	return float32(float64Val)
}
