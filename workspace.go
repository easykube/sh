package sh

//工作区
type WorkSpace struct {
}

//新建工作区
func NewWorkSpace() *WorkSpace {
	return &WorkSpace{}
}

//初始化工作区
func (this *WorkSpace) Init() {

}

//运行命令
func (this *WorkSpace) Run(cmdLine string) {

}

//上传文件
func (this *WorkSpace) Put(dst, src string) {

}

//本地运行
func (this *WorkSpace) Local(cmdLine string) {

}
