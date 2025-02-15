//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Preferences
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package settings

import (
	"EmbyExplorer_for_Windows/assets"
	"encoding/json"
	"github.com/lxn/walk"
	"io"
	"os"
	"path/filepath"
)

type Settings struct {
	WindowRect       walk.Rectangle
	EmbySecure       bool
	EmbyServer       string
	EmbyPort         string
	EmbyUser         string
	EmbyPassword     []byte
	LastExportFolder string
}

var settings Settings

func SetPreferencesDetail(secure bool, server string, port string, user string, password string) {
	settings.EmbySecure = secure
	settings.EmbyServer = server
	settings.EmbyPort = port
	settings.EmbyUser = user
	settings.EmbyPassword = []byte(password)
}

func SetWindowBounds(rect walk.Rectangle) {
	settings.WindowRect = rect
}

func GetWindowBounds() walk.Rectangle {
	return settings.WindowRect
}

func GetConnectionSettings() (secure bool, server string, port string, user string, password string) {
	return settings.EmbySecure, settings.EmbyServer, settings.EmbyPort,
		settings.EmbyUser, string(settings.EmbyPassword)
}

func SetLastExportFolder(path string) {
	settings.LastExportFolder = path
}

func GetLastExportFolder() string {
	return settings.LastExportFolder
}

func Valid() bool {
	return settings.EmbyServer != "" &&
		settings.EmbyPort != "" &&
		settings.EmbyUser != "" &&
		len(settings.EmbyPassword) > 0
}

func SavePreferences() {
	j, err := json.Marshal(settings)
	if err == nil {
		dir, _ := os.UserConfigDir()
		dir = filepath.Join(dir, assets.AppName)
		_, err := os.Stat(dir)
		if err != nil {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				panic(err)
			}
		}
		fname := filepath.Join(dir, preferencesFileName)
		_ = os.WriteFile(fname, j, 0644)
	}
}

func LoadPreferences() error {
	dir, err := os.UserConfigDir()
	dir = filepath.Join(dir, assets.AppName)
	fname := filepath.Join(dir, preferencesFileName)
	j, err := os.Open(fname)
	if err == nil {
		byteValue, _ := io.ReadAll(j)
		_ = j.Close()
		_ = json.Unmarshal(byteValue, &settings)
		return nil
	}
	return err
}

const preferencesFileName = "org.janbuchholz.embyexplorerforwindows.json"
