package result

var (
	SUCCESS = Template{200, "Success"}
	FAILURE = Template{500, "服务器繁忙"}
	NoLog   = Template{600, "未登录"}
	NoReg   = Template{610, "未注册"}
)

type Template struct {
	Code int
	Msg  string
}
