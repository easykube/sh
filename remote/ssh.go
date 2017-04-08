package remote

import (
	"fmt"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshSession struct {
	name       string
	host       string
	port       int
	user       string
	password   string
	client     *ssh.Client
	session    *ssh.Session
	sftpClient *sftp.Client
}

func NewSSHSession() *SshSession {
	return &SshSession{}
}

//初始化设置
func (this *SshSession) Init(conf *Config) {
	this.Close()
	this.name = conf.Name
	this.host = conf.Host
	this.port = 22
	if conf.Port > 0 {
		this.port = conf.Port
	}

	this.user = conf.User
	this.password = conf.Password

}

//获取配置
func (this *SshSession) Conf() *Config {
	conf := NewConfig()
	conf.UseWinRm = false
	conf.Name = this.name
	conf.Host = this.host
	conf.Port = this.port
	conf.User = this.user
	conf.Password = this.password
	conf.IsLocal = false
	return conf

}

//是否本地会话
func (this *SshSession) IsLocal() bool {
	return false
}

func (this *SshSession) Open() error {

	addr := fmt.Sprintf("%s:%d", this.host, this.port)
	cfg := &ssh.ClientConfig{
		User: this.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(this.password)},
	}
	client, err := ssh.Dial("tcp", addr, cfg)
	this.client = client
	if err != nil {
		return err
	} else {
		session, err := client.NewSession()
		this.session = session
		if err != nil {
			os.Exit(2)
		}
		return err
	}
}

func (this *SshSession) Close() error {
	var err1 error
	var err2 error
	if this.session != nil {
		err1 = this.session.Close()
	}

	if this.sftpClient != nil {
		err1 = this.sftpClient.Close()
	}

	if err1 == nil && err2 == nil {
		return nil
	}

	this.session = nil
	this.sftpClient = nil

	return nil
}

func (this *SshSession) Run(cmd string) (string, error) {
	out, err := this.session.CombinedOutput(cmd)
	//fmt.Printf("RESULT:  %s:\n%s", this.host, string(result))
	return string(out), err
}

func (this *SshSession) initSftp() error {
	if this.sftpClient == nil {
		sftpClient, err := sftp.NewClient(this.client)
		if err != nil {
			return err
		}
		this.sftpClient = sftpClient

	}
	return nil

}

//显示目录
func (this *SshSession) ListDir(remoteDir string) ([]os.FileInfo, error) {
	err := this.initSftp()
	if err != nil {
		return nil, err
	}
	return this.sftpClient.ReadDir(remoteDir)
}

func (this *SshSession) UpFile(src string, dst string) error {
	err := this.initSftp()
	if err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := this.sftpClient.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 4096)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}

	return nil
}

func (this *SshSession) DownFile(src string, dst string) error {

	err := this.initSftp()
	if err != nil {
		return err
	}

	srcFile, err := this.sftpClient.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 4096)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}

	return nil

}
