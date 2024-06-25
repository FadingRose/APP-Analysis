package apkloader

import (
	"crypto/md5"
	"encoding/hex"
	"fadingrose/app-analyzer/internal/types"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tuotoo/qrcode"
)

type APK = types.APK

// ApkLoader is responsible for loading and managing APKs.
type ApkLoader struct {
	ApkList []APK // ApkList stores the list of loaded APKs.
}

// NewApkLoader initializes and returns a new ApkLoader instance.
func NewApkLoader() ApkLoader {
	return ApkLoader{
		ApkList: []APK{}, // Initializes an empty APK list.
	}
}

// AddApk adds a new APK to the ApkList.
func (a *ApkLoader) AddApk(apk APK) {
	a.ApkList = append(a.ApkList, apk)
}

// ReadApkList returns the list of loaded APKs.
func (a *ApkLoader) ReadApkList() []APK {
	return a.ApkList
}

// PrintApkList prints the details of all APKs in ApkList.
func (a *ApkLoader) PrintApkList() {
	for i, apk := range a.ApkList {
		fmt.Printf("APK #%d:\n", i+1)
		fmt.Printf("	Path: %s\n", apk.Path)
		fmt.Printf("	Name: %s\n", apk.Name)
		fmt.Printf("	Version: %s\n", apk.Version)
		fmt.Printf("	Size: %d bytes\n", apk.Size)
		fmt.Printf("	MD5: %s\n", apk.MD5)
		fmt.Printf("	Package Name: %s\n", apk.PackageName)
		fmt.Printf("	Developer: %s\n", apk.Developer)
		fmt.Printf("	Signature: %s\n", apk.Signature)
		fmt.Printf("	Permissions: %v\n", apk.Permissions)
		fmt.Printf("	Activities: %v\n", apk.Activities)
		fmt.Printf("	Services: %v\n", apk.Services)
		fmt.Printf("	Providers: %v\n", apk.Providers)
		fmt.Printf("	Receivers: %v\n", apk.Receivers)
		fmt.Printf("	MetaData: %v\n", apk.MetaData)
		fmt.Printf("	Min SDK Version: %s\n", apk.MinSdkVersion)
		fmt.Printf("	Target SDK Version: %s\n", apk.TargetSdkVersion)
	}
}

// LoadApkFromQR loads an APK from a QR code.
// Note: This is a stub implementation. Actual QR code processing logic needs to be added.
func (a *ApkLoader) LoadApkFromQR(imageUrl string) APK {
	// Read QR code from image
	fi, err := os.Open(imageUrl)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	// Decode QR code
	qrMatrix, err := qrcode.Decode(fi)
	if err != nil {
		panic(err)
	}

	// Fetch APK from URL
	return a.LoadApkFromURL(qrMatrix.Content)
}

// LoadApkFromURL loads an APK from a specified URL.
func (a *ApkLoader) LoadApkFromURL(url string) APK {
	// Process the URL
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	apkData, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Create a temporary file to store the APK
	tempFile, err := os.CreateTemp("", "*.apk")
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()

	tempFile.Write(apkData)
	tempFilePath := tempFile.Name()

	// Calculate MD5
	hash := md5.Sum(apkData)
	md5Str := hex.EncodeToString(hash[:])

	apk := APK{
		Name:    tempFile.Name(),
		Version: "1.0", // Version should be extracted from APK metadata
		Size:    int64(len(apkData)),
		MD5:     md5Str,
		Path:    tempFilePath,
	}
	a.ApkList = append(a.ApkList, apk)
	return apk
}

// LoadApkFromFile loads an APK from a local file path.
func (a *ApkLoader) LoadApkFromFile(path string, apkType string) APK {
	fileData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Calculate MD5
	hash := md5.Sum(fileData)
	md5Str := hex.EncodeToString(hash[:])

	apk := APK{
		Name:    "LocalApp",
		Version: "1.0", // Version should be extracted from APK metadata
		Size:    int64(len(fileData)),
		MD5:     md5Str,
		Path:    path,
		Type:    apkType,
	}
	apk = ExtractAPKInfo(apk)
	a.ApkList = append(a.ApkList, apk)
	return apk
}
