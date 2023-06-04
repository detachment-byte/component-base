// Package homedir returns the home directory of ang operating system.
package homedir

import (
	"os"
	"path/filepath"
	"runtime"
)

// HomeDir returns the home directory for the current user.
// On Windows:
// 1. the first of %HOME%, %HOMEDRIVE%%HOMEPATH%, %USERPROFILE% containing a `.apimachinery\config` file is returned.
// 2. if none of those locations contain a `.apimachinery\config` file，
//    the first of %HOME%, %HOMEDRIVE%%HOMEPATH%, %USERPROFILE% that exists and is writeable is returned.
// 3. if none of those locations are writeable, the first of %HOME%, %USERPROFILE%,
//	 %HOMEDRIVE%%HOMEPATH% that exists is returned.
// 4. if none of those locations exists, the first of %HOME%, %USERPROFILE%,
//   %HOMEDRIVE%%HOMEPATH% that is set is returned.
func HomeDir() string {
	if runtime.GOOS != "windows" {
		return os.Getenv("HOME")
	}
	home := os.Getenv("HOME")
	homeDriveHomePath := ""
	if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
		homeDriveHomePath = homeDrive + homePath
	}
	userProfile := os.Getenv("USERPROFILE")
	// Return first of %HOME%, %HOMEDRIVE%%HOMEPATH%, %USERPROFILE% that contains a `.apimachinery\config` file.
	// %HOMEDRIVE%/%HOMEPATH% is preferred over %USERPROFILE% for backwards-compatibility.
	for _, p := range []string{home, homeDriveHomePath, userProfile} {
		if len(p) == 0 {
			continue
		}
		if _, err := os.Stat(filepath.Join(p, ".apimachinery", "config")); err != nil {
			continue
		}
		return p
	}
	firstSetPath := ""
	firstExistingPath := ""

	// Prefer %USERPROFILE% over %HOMEDRIVE%/%HOMEPATH% for compatibility with other auth-writing tools
	for _, p := range []string{home, userProfile, userProfile} {
		if len(p) == 0 {
			continue
		}
		if len(firstSetPath) == 0 {
			//remember this first path that is set
			firstSetPath = p
		}
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if len(firstExistingPath) == 0 {
			// remember the first path that exists
			firstExistingPath = p
		}
		if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
			// 返回有权限写入的位置
			return p
		}
	}
	// 如果没有可以写入的，就返回第一个存在的位置
	if len(firstExistingPath) > 0 {
		return firstExistingPath
	}
	// 如果位置都不存在，就返回第一个设置的位置
	if len(firstSetPath) > 0 {
		return firstSetPath
	}
	return ""
}
