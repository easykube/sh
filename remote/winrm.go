package remote

/**
服务端配置，设置为认证方式为Basic,不用加密传输，这时可以直接使用用户名和密码登录，传输为明文
    在cmd中
    winrm quickconfig
    y
    winrm set winrm/config/service/Auth '@{Basic="true"}'
    winrm set winrm/config/service '@{AllowUnencrypted="true"}'
    winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="1024"}'

    在Powershell中
	设置
	Set-Item -Path "WSMan:\localhost\Service\Auth\Basic" -Value $true
	Set-Item -Path "WSMan:\localhost\Service\AllowUnencrypted" -Value $true
　　 查看
	Get-ChildItem WSMan:\localhost\Service\Auth | Where {$_.Name -eq "Basic"}
	Get-ChildItem WSMan:\localhost\Service | Where {$_.Name -eq "AllowUnencrypted"}
**/

import (
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"fmt"

	"github.com/easykube/sh/remote/winrmcp"
	"github.com/masterzen/winrm"
)

type WimRmSession struct {
	name     string
	host     string
	user     string
	password string
	port     int
	client   *winrm.Client
	config   *winrmcp.Config
	winrmcp  *winrmcp.Winrmcp
}

type fileInfo struct {
	name  string
	size  int64
	isdir bool
	mode  os.FileMode
	mtime time.Time
	sys   interface{}
}

// Name returns the base name of the file.
func (fi *fileInfo) Name() string { return fi.name }

// Size returns the length in bytes for regular files; system-dependent for others.
func (fi *fileInfo) Size() int64 { return fi.size }

// Mode returns file mode bits.
func (fi *fileInfo) Mode() os.FileMode { return fi.mode }

// ModTime returns the last modification time of the file.
func (fi *fileInfo) ModTime() time.Time { return fi.mtime }

// IsDir returns true if the file is a directory.
func (fi *fileInfo) IsDir() bool { return fi.isdir }

func (fi *fileInfo) Sys() interface{} { return fi.sys }

func NewWinRmSession() *WimRmSession {
	return &WimRmSession{}
}

//初始化设置
func (this *WimRmSession) Init(conf *Config) {
	this.Close()
	this.name = conf.Name
	this.host = conf.Host
	this.port = 5985
	if conf.Port > 0 {
		this.port = conf.Port
	}

	this.user = conf.User
	this.password = conf.Password

}

//获取配置
func (this *WimRmSession) Conf() *Config {
	conf := NewConfig()
	conf.UseWinRm = true
	conf.Name = this.name
	conf.Host = this.host
	conf.Port = this.port
	conf.User = this.user
	conf.Password = this.password
	conf.IsLocal = false
	return conf

}

//是否本地会话
func (this *WimRmSession) IsLocal() bool {
	return false

}

func (this *WimRmSession) initwincp() {
	if this.winrmcp == nil {
		this.winrmcp = winrmcp.NewWinrmcp2(this.client, this.config)
	}
}

func (this *WimRmSession) Open() error {

	this.config = &winrmcp.Config{}
	this.config.Auth.Password = this.password
	this.config.Auth.User = this.user
	this.config.Https = false
	this.config.Insecure = true
	this.config.MaxOperationsPerShell = 15
	this.config.OperationTimeout = 60

	endpoint := winrm.NewEndpoint(this.host, this.port, false, true, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, this.user, this.password)
	if err != nil {
		return err
	}
	this.client = client
	return nil
}

func (this *WimRmSession) Close() error {
	if this.winrmcp != nil {
		this.winrmcp = nil
	}

	if this.client != nil {
		this.client = nil
	}
	this.client = nil
	return nil
}

func (this *WimRmSession) Run(cmd string) (string, error) {
	out1, out2, out3, err := this.client.RunWithString(cmd, "")
	//fmt.Printf("RESULT:  %s:\n%s %s,%d", this.host, out1, out2, out3)
	out := fmt.Sprintf("%s\n    error:\n%s:\n    exitcode:%d", out1, out2, out3)
	return out, err
}

func (this *WimRmSession) ListDir(remoteDir string) ([]os.FileInfo, error) {
	this.initwincp()
	list, err := this.winrmcp.List(remoteDir)
	if err != nil {
		return nil, err
	}
	files := make([]os.FileInfo, len(list))
	for i, f := range list {
		fi := &fileInfo{}
		fi.name = f.Path
		fi.size = int64(f.Length)
		if strings.Index(f.Mode, "d") == 0 {
			fi.isdir = true
		} else {
			fi.isdir = false
		}

		files[i] = fi

	}
	return files, nil
}

func (this *WimRmSession) Write(toPath string, src io.Reader) error {
	this.initwincp()
	return this.winrmcp.Write(toPath, src)
}

func (this *WimRmSession) UpFile(src string, dst string) error {
	this.initwincp()
	return this.winrmcp.Copy(src, dst)
}

func (this *WimRmSession) DownFile(src string, dst string) error {
	return errors.New("没有实现DownFile")
}
