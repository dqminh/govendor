package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type dependency struct {
	Vcs  string `json:"vcs"`  // type of the repository
	Repo string `json:"repo"` // repository of the dependency
	Rev  string `json:"rev"`  // revision of the dependency
	Path string `json:"path"` // where we stored the dependency in GOPATH
}

func (d *dependency) install(srcPath string) error {
	log.Printf("Installing %v\n", d.Repo)
	if err := os.Chdir(srcPath); err != nil {
		return fmt.Errorf("Failed to navigate to srcPath")
	}
	if err := os.RemoveAll(d.Path); err != nil {
		return fmt.Errorf("Failed to remove %s", d.Path)
	}
	if err := os.MkdirAll(d.Path, 0755); err != nil {
		return fmt.Errorf("Failed to mkdir %s", d.Path)
	}
	switch d.Vcs {
	case "git":
		if err := exec.Command("git", "clone", "--quiet", "--no-checkout",
			d.Repo, d.Path).Run(); err != nil {
			return fmt.Errorf("Failed to clone git %s: %v", d.Path, err)
		}
		if err := os.Chdir(d.Path); err != nil {
			return fmt.Errorf("Failed to change to %s", d.Path)
		}
		if err := exec.Command("git", "reset", "--quiet", "--hard",
			d.Rev).Run(); err != nil {
			return fmt.Errorf("Failed to change to git rev %s: %v", d.Rev, err)
		}
	case "hg":
		if err := exec.Command("hg", "clone", "--quiet", "--updaterev", d.Rev,
			d.Repo, d.Path).Run(); err != nil {
			return fmt.Errorf("Failed to clone hg %s: %v", d.Path, err)
		}
	}
	return nil
}

func main() {
	var depPath string
	if len(os.Args) == 1 {
		depPath = "deps.json"
	} else {
		depPath = os.Args[1]
	}

	// read dependency file
	deps, err := readDependencies(depPath)
	if err != nil {
		log.Fatalf("Invalid dependency file %s: %v\n", depPath, err)
	}

	// write .env file
	writeEnv()

	// create _vendor directory in current dir
	srcPath, err := createVendor()
	if err != nil {
		log.Fatalf("Failed to create _vendor directory: %v\n", err)
	}

	// start installing dependencies
	for _, dep := range deps {
		if err := dep.install(srcPath); err != nil {
			log.Fatal("Failed to install %v: %v", dep.Repo, err)
		}
	}
	fmt.Println("Dependencies written into _vendor/src")
}

func createVendor() (string, error) {
	srcPath := "_vendor/src"
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		return "", errors.New("Failed to create _vendor directory")
	}
	return filepath.Abs(srcPath)
}

func readDependencies(depPath string) ([]*dependency, error) {
	depData, err := ioutil.ReadFile(depPath)
	if err != nil {
		return nil, err
	}
	var deps []*dependency
	if err := json.Unmarshal(depData, &deps); err != nil {
		return nil, err
	}
	return deps, nil
}

var envTips = `Written "export GOPATH=$(pwd)/_vendor:$GOPATH" into .env
You can autoload .env file with "https://github.com/kennethreitz/autoenv"
`

func writeEnv() {
	_, err := os.Stat(".env")
	if err == nil {
		log.Println(".env exists. Skipping...")
		return
	}
	if os.IsNotExist(err) {
		err := ioutil.WriteFile(".env",
			[]byte(`export GOPATH=$(pwd)/_vendor:$GOPATH`), 0755)
		if err != nil {
			log.Fatal("Failed to write .env file")
		}
		fmt.Println(envTips)
	}
}
