package main

import (
	"os"
	"log"
	"path"
	"io/ioutil"
	"fmt"
	"time"
	"net"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Connect
func connect(user, password, host string, port int) (*sftp.Client, error) {

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client

		err error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil

}

// UploadFile
func uploadFile(localFilePath, remoteDir string) {

	srcFile, err := os.Open(localFilePath)

	if err != nil {
		log.Fatal(err)
	}

	defer srcFile.Close()

	var remoteFileName = getChangedFileName(localFilePath)

	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))

	if err != nil {
		log.Fatal("dst error: ", err, path.Join(remoteDir, remoteFileName))
	}

	defer dstFile.Close()

	buf, err := ioutil.ReadAll(srcFile)

	if err != nil {
		panic(err)
	}

	dstFile.Write(buf)

	log.Print("uploaded file: ", localFilePath)
}

// UploadDirectory
func uploadDirectory(localPath, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)

	if err != nil {
		panic(err)
	}

	changedFileName := getChangedFileName(localPath)

	sftpClient.Mkdir(path.Join(remotePath, changedFileName))

	log.Print("make remote dir: ", path.Join(remotePath, getChangedFileName(localPath)))

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, changedFileName, backupDir.Name())

		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			uploadDirectory(localFilePath, remotePath)
		} else {
			uploadFile(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
}

// Remove whatever  a file or directory
func remove(filename, remoteDir string) {
	remoteFileName := path.Join(remoteDir, getChangedFileName(filename))
	fileInfo, err := sftpClient.Stat(remoteFileName)

	if err != nil {
		log.Println(err)
		return
	}

	if fileInfo.IsDir() {
		removeDirectory(remoteFileName, remoteDir)
	} else {
		removeFile(remoteFileName)
	}
}

// Remove file
func removeFile(filename string) {

	err := sftpClient.Remove(filename)

	if err != nil {
		log.Fatal("Can' remove file :", filename)
	} else {
		log.Print("Removed file: ", filename)
	}
}

// remove Directory
func removeDirectory(dir , remoteDir string) {
	fileInfos, err := sftpClient.ReadDir(dir)

	if err != nil {
		log.Fatal("Readdir err: ", err)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			fmt.Println("remove directory: ", path.Join(dir, fileInfo.Name()))
			removeDirectory(path.Join(dir, fileInfo.Name()), remoteDir)
		} else {
			fmt.Println(path.Join(remoteDir, getChangedFileName(dir)))
			removeFile(path.Join(remoteDir, getChangedFileName(dir), fileInfo.Name()))
		}
	}

	sftpClient.RemoveDirectory(dir)
}