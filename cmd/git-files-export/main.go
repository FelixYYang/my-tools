package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	wDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	args := os.Args
	exportName := "gitExport"
	/*if len(args) > 1 && len(args[1]) != 0 {
		exportName = args[1]
	}*/
	commitCurrent := "HEAD"
	commitPrevious := "HEAD~"
	if len(args) > 1 {
		commitCurrent = args[1]
		if len(args) > 2 {
			commitPrevious = args[2]
		}
	}

	exportPath := filepath.Join(wDir, exportName)
	cmdIsGitDir := exec.Command("git", []string{
		"rev-parse",
		"--is-inside-work-tree",
	}...)
	isGitDir, err := cmdIsGitDir.CombinedOutput()
	if err != nil {
		log.Println(string(isGitDir))
		log.Fatalln(err)
	}
	if !strings.HasPrefix(string(isGitDir), "true") {
		log.Fatalln(string(isGitDir))
	}

	cmdGetTopDir := exec.Command("git", strings.Split("rev-parse --show-toplevel", " ")...)
	topDirBytes, err := cmdGetTopDir.CombinedOutput()
	if err != nil {
		log.Println(string(topDirBytes))
		log.Fatalln(err)
	}
	toDir := string(bytes.TrimRight(topDirBytes, "\n"))

	cmdGetFileList := exec.Command("git", strings.Split(fmt.Sprintf("diff --name-only %s %s --", commitPrevious, commitCurrent), " ")...)
	fileListBytes, err := cmdGetFileList.CombinedOutput()
	if err != nil {
		log.Println(string(fileListBytes))
		log.Fatalln(err)
	}
	filesBytes := bytes.Split(bytes.TrimRight(fileListBytes, "\n"), []byte("\n"))
	println(exportPath)
	for _, filesByte := range filesBytes {
		fileRelativePath := string(filesByte)
		fileP := filepath.Join(exportPath, fileRelativePath)
		fileDir := filepath.Dir(fileP)
		_, err := os.Stat(fileDir)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
				log.Fatalln(err)
			}
		}
		sourceFile, err := os.Open(filepath.Join(toDir, fileRelativePath))
		if err != nil {
			log.Println(err)
			continue
		}
		file, err := os.Create(fileP)
		if err != nil {
			log.Fatalln(err)
		}
		if _, err := io.Copy(file, sourceFile); err != nil {
			log.Fatalln(err)
		}
		_ = file.Close()
		_ = sourceFile.Close()
		log.Println("filepath" + fileP)
	}
}
