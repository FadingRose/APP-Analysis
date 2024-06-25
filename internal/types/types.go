package types

import "encoding/xml"

// APK represents the details of an Android application package.
type APK struct {
	Path             string            // Path is the file path to the APK.
	Name             string            // Name is the name of the application.
	Version          string            // Version is the version of the application.
	Size             int64             // Size is the size of the APK file in bytes.
	MD5              string            // MD5 is the MD5 checksum of the APK file.
	PackageName      string            // PackageName is the unique identifier for the application.
	Developer        string            // Developer is the name of the application's developer.
	Signature        string            // Signature is the cryptographic signature of the APK.
	Permissions      []string          // Permissions requested by the application.
	Activities       []string          // Activities defined in the application.
	Services         []string          // Services defined in the application.
	Providers        []string          // Content providers defined in the application.
	Receivers        []string          // Broadcast receivers defined in the application.
	MetaData         map[string]string // Meta-data defined in the application.
	MinSdkVersion    string            // Minimum SDK version required by the application.
	TargetSdkVersion string            // Target SDK version for the application.
	Type             string            // Type is the type of the application.
	Urls             []string          // URLs extracted from the application.
}

// Manifest represents the struct for parsing AndroidManifest.xml to APK
type Manifest struct {
	PackageName string `xml:"package,attr"`
	VersionName string `xml:"versionName,attr"`
	VersionCode string `xml:"versionCode,attr"`
	Application struct {
		XMLName    xml.Name `xml:"application"`
		Label      string   `xml:"label,attr"`
		Activities []struct {
			Name string `xml:"name,attr"`
		} `xml:"activity"`
		Services []struct {
			Name string `xml:"name,attr"`
		} `xml:"service"`
		Providers []struct {
			Name string `xml:"name,attr"`
		} `xml:"provider"`
		Receivers []struct {
			Name string `xml:"name,attr"`
		} `xml:"receiver"`
		MetaData []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"meta-data"`
	} `xml:"application"`
	Permissions []struct {
		Name string `xml:"name,attr"`
	} `xml:"uses-permission"`
	UsesSdk struct {
		MinSdkVersion    string `xml:"minSdkVersion,attr"`
		TargetSdkVersion string `xml:"targetSdkVersion,attr"`
	} `xml:"uses-sdk"`
}

// /<path_to_apk>/a.apk

// read file -> open it -> pointer
