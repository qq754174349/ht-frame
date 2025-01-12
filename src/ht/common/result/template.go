package result

var (
	SUCCESS   = Template{200, "Success"}
	FAILURE   = Template{500, "Failure"}
	NoReg     = Template{520, "未注册"}
	RepeatReg = Template{521, "重复注册"}
)

type Template struct {
	Code int
	Msg  string
}
