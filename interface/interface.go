package iface

import "fadingrose/app-analyzer/internal/types"

type APK = types.APK

// ApkLoaderInterface defines the interface for loading APKs
// and includes methods for reading the APK list and loading
// APKs from external sources.
type ApkLoaderInterface interface {
	// ReadApkList returns a list of APKs.
	ReadApkList() []APK

	// LoadApkFromExternal provides methods for loading APKs from external sources.
	LoadApkFromExternal
}

// LoadApkFromExternal defines the interface for loading APKs
// from different external sources such as QR codes and URLs.
type LoadApkFromExternal interface {
	// LoadApkFromQR loads an APK from a QR code.
	LoadApkFromQR(imageUrl string) APK

	// LoadApkFromURL loads an APK from a specified URL.
	LoadApkFromURL(url string) APK

	// LoadApkFromFile loads an APK from a local file path.
	LoadApkFromFile(path string, apkType string) APK
}
