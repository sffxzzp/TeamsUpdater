package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type (
	Teams struct {
		Ver    string `json:"version"`
		Client *http.Client
		URL    string
		NewVer string
	}
	RetData struct {
		IsUpdateAvailable bool   `json:"isUpdateAvailable"`
		NugetPackagePath  string `json:"nugetPackagePath"`
		ReleasesPath      string `json:"releasesPath"`
		URL               string `json:"url"`
		ScenarioCode      int    `json:"scenarioCode"`
		DeltaPackagePath  string `json:"deltaPackagePath"`
		DeltaReleasesPath string `json:"deltaReleasesPath"`
	}
)

func (t *Teams) httpGet(url string) []byte {
	retry := 3
	for retry > 0 {
		res, err := t.Client.Get(url)
		if err != nil {
			retry--
			continue
		}
		defer res.Body.Close()
		rbytes, err := io.ReadAll(res.Body)
		if err != nil {
			retry--
			continue
		}
		return rbytes
	}
	return []byte{}
}

func (t *Teams) checkUpdates() bool {
	fmt.Println("Checking for updates...")
	data := t.httpGet(fmt.Sprintf("https://teams.microsoft.com/desktopclient/update/%s/windows/OSBit?ring=general", t.Ver))
	var r RetData
	json.Unmarshal(data, &r)
	if r.IsUpdateAvailable {
		t.URL = r.NugetPackagePath
		t.NewVer = strings.Split(strings.Split(t.URL, "production-windows/")[1], "/")[0]
		return true
	} else {
		return false
	}
}

func (t *Teams) pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (t *Teams) download() error {
	fmt.Printf("Downloading: %s\n", t.URL)
	if t.pathExists("..\\app") {
		os.Rename("..\\app\\", "..\\app-"+t.Ver+"\\")
	}
	os.Mkdir("..\\app\\", 0755)
	res := t.httpGet(t.URL)
	zipReader, err := zip.NewReader(bytes.NewReader(res), int64(len(res)))
	if err != nil {
		return err
	}
	for _, file := range zipReader.File {
		fmt.Println("extracting", file.Name)
		unzipPath := filepath.Join("..\\app\\", file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(unzipPath, file.Mode())
			continue
		} else {
			dir := filepath.Dir(unzipPath)
			if !t.pathExists(dir) {
				os.MkdirAll(dir, 0755)
			}
		}
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		targetFile, err := os.OpenFile(unzipPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}
	fmt.Println("Download complete!")
	return nil
}

// a function that writes new version to config.json
func (t *Teams) writecfg() {
	cfg := fmt.Sprintf("{\"Version\": \"%s\"}", t.NewVer)
	os.WriteFile("config.json", []byte(cfg), 0644)
}

func newTeams() *Teams {
	teams := &Teams{
		Client: &http.Client{},
	}
	cCont, err := os.ReadFile("config.json")
	if err != nil {
		cCont = []byte("{\"Version\": \"0.0.0.0\"}")
	}
	json.Unmarshal(cCont, &teams)
	return teams
}

func main() {
	teams := newTeams()
	if teams.checkUpdates() {
		fmt.Printf("A newer version found, %s -> %s\n", teams.Ver, teams.NewVer)
		fmt.Println("Please close Teams and press enter to continue...")
		fmt.Scanln()
		teams.download()
		teams.writecfg()
		fmt.Println("Update complete!")
	} else {
		fmt.Println("No update available.")
	}
}
