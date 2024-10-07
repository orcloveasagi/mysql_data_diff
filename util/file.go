package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func AbsPath(relPaths ...string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Printf("get absolute path error : %v", err)
		return "nil", err
	}
	for _, dir := range relPaths {
		workingDir = filepath.Join(workingDir, dir)
	}
	return workingDir, nil
}

func JoinPath(paths ...string) (result string) {
	for _, dir := range paths {
		result = filepath.Join(result, dir)
	}
	return result
}

func ReadFile(dirs ...string) ([]byte, error) {
	path, err := AbsPath(dirs...)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func MakeDir(relativeDir ...string) (absoluteDir string, err error) {
	absoluteDir, err = AbsPath(relativeDir...)
	err = os.MkdirAll(absoluteDir, os.ModePerm)
	if err != nil {
		log.Printf("mkdir error : %v", err)
		return absoluteDir, err
	}
	return absoluteDir, err
}

func MakeFile(dir, file string) (*os.File, error) {
	path := JoinPath(dir, file)
	create, err := os.Create(path)
	if err != nil {
		log.Printf("make file error : %v", err)
		return nil, err
	}
	return create, nil
}

func WriteLineToFile(data []string, filePath ...string) (dir string, err error) {
	dir, err = MakeDir(filePath[:len(filePath)-1]...)
	if err != nil {
		return dir, err
	}
	file, err := MakeFile(dir, filePath[len(filePath)-1])
	if err != nil {
		return dir, err
	}
	defer Close(file)
	for _, line := range data {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			log.Printf("write file error: %v", err)
			return dir, err
		}
	}

	log.Printf("%s generated", file.Name())
	return dir, err
}

func WriteFile(data []byte, filepath ...string) error {
	path, err := AbsPath(filepath...)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}

func OpenExplorer(path string) error {
	var cmd string
	switch runtime.GOOS {
	case "windows":
		cmd = "explorer"
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	default:
		return fmt.Errorf("not support platform")
	}
	return exec.Command(cmd, path).Start()
}
