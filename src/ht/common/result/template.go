package result

const (
	SUCCESS = Template{200, "Success"}
)

type Template struct {
	code int
	msg  string
}
