package utils

//Error 自定义错误接口
type Error interface {
	error
	Status() int
}
//StatusError have Error and Status method 
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

//Status code
func (se StatusError) Status() int {
	return se.Code
}
