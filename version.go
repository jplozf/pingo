package main

const MajorVersion = "0"

// This variable will be overwritten at compile time
var Version = "dev"

func GetDisplayVersion() string {
	if Version == "dev" {
		return MajorVersion + ".x-dev"
	}
	return Version
}
