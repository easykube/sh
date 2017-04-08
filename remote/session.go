package remote

import "os"

//会话配置
type Config struct {
	//主机标识名
	Name string `yaml:"name"`
	//是否本地主机
	IsLocal bool `yaml:"local"`
	//主机名或ip
	Host string `yaml:"host"`
	//端口
	Port int `yaml:"port"`
	//登录的用户
	User string `yaml:"user"`
	//密码
	Password string `yaml:"password"`
	//是否使用WinRm远程连接，否则使用ssh
	UseWinRm bool `yaml:"winrm"`
}

//创建配置
func NewConfig() *Config {
	return &Config{}
}

//会话
type Session interface {
	//初始化设置
	Init(conf *Config)
	//获取配置
	Conf() *Config
	//是否本地会话
	IsLocal() bool
	//打开会话
	Open() error
	//关闭会话
	Close() error
	//执行命令
	Run(cmd string) (string, error)

	//运行远程命令文件
	//RunFile(scriptFile string) error

	//显示目录
	ListDir(remoteDir string) ([]os.FileInfo, error)
	//上传文件
	UpFile(src string, dst string) error
	//下载文件
	DownFile(src string, dst string) error
}

//创建远程会话
func NewSession(conf *Config) Session {
	var s Session
	if conf.UseWinRm {
		s = NewWinRmSession()
		s.Init(conf)
	} else {
		s = NewSSHSession()
		s.Init(conf)
	}

	return s
}
