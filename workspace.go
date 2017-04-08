package sh

/**
**/
import (
	"github.com/easykube/sh/remote"
	"github.com/easykube/util"
)

//工作区
type WorkSpace struct {
	//配置
	conf *remote.Config

	//远程session
	session remote.Session
}

//新建工作区
func NewWorkSpace() *WorkSpace {
	return &WorkSpace{}
}

//加载配置文件并初始化
func (this *WorkSpace) Load(file string) {

}

//初始化
func (this *WorkSpace) Init(conf *remote.Config) {
	this.conf = conf
	if !conf.IsLocal {
		if this.session == nil {
			this.session = remote.NewSession(conf)
		}
		err := this.session.Open()
		util.LogError("work space init", err)
	}
}

//运行命令
func (this *WorkSpace) Run(cmdLine string) {
	if this.session != nil {
		out, err := this.session.Run(cmdLine)
		util.LogError("WorkSpace run", err)
		println(out)
	}

}

//上传文件
func (this *WorkSpace) Put(src, dst string) {
	if this.session != nil {
		err := this.session.UpFile(src, dst)
		util.LogError("WorkSpace up", err)
	}

}

//本地运行
func (this *WorkSpace) Local(cmdLine string) {
	out, err := util.ExecCmdLine(cmdLine)
	util.LogError("WorkSpace local", err)
	println(out)
}
