package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	SshKeyPath  string
	SshPassword string
)

func Connect(ip, port string) ([]byte, error) {
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{hostAuthFunc()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}
	// 实例化client
	addr := fmt.Sprintf("%s:%s", ip, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("创建ssh client 失败%s\n", err)
		return nil, err
	}
	defer sshClient.Close()

	// 创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("创建ssh session 失败%s\n", err)
		return nil, err
	}
	defer session.Close()

	// 执行命令
	combo, _ := session.CombinedOutput("occtl --json show users")
	//if err != nil{
	//	log.Fatal("远程执行cmd 失败",err)
	//}
	return combo, nil
}

// publicKeyAuthFunc 获取私钥 用来免密登录
func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	key, err := os.ReadFile(kPath)
	if err != nil {
		log.Printf("ssh key file read failed %s\n", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Printf("ssh key signer failed%s\n", err)
	}
	return ssh.PublicKeys(signer)
}

// passwordAuthFunc 使用密码登录
func passwordAuthFunc(pwd string) ssh.AuthMethod {
	return ssh.Password(pwd)
}

// hostAuthFunc 主机登录类型判断 兼容秘钥 和 密码邓丽
func hostAuthFunc() ssh.AuthMethod {
	_, err := os.ReadFile(SshKeyPath)
	if err != nil {
		return passwordAuthFunc(SshPassword)
	}
	return publicKeyAuthFunc(SshKeyPath)
}
