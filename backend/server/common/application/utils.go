package application

import "os"

func IsProd() bool {
	_, isSet := os.LookupEnv("IS_PROD")
	return isSet
}

func IsCloudBuild() bool {
	_, isSet := os.LookupEnv("IS_CLOUD_BUILD")
	return isSet
}
