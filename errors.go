package yandex

type Error struct {
	Code    int    `json:"code"`
	ApiCode string `json:"api_code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GetApiCode() string {
	return e.ApiCode
}

func (e *Error) GetCode() int {
	return e.Code
}
