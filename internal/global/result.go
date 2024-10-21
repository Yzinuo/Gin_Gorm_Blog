package global



/*
响应设计方案：不使用 HTTP 码来表示业务状态, 采用业务状态码的方式
- 只要能到达后端的请求, HTTP 状态码都为 200
- 业务状态码为 0 表示成功, 其他都表示失败
- 当后端发生 panic 并且被 gin 中间件捕获时, 才会返回 HTTP 500 状态码
*/

var(
	SUCCESS = 0
	FAIL = 500
	_code  = map[int]struct{}{}
	_msg   = make(map[int]string)
)

type Result struct{
	Code 	int
	Msg 	string
}

func (e Result) GetCode() int {
	return e.Code
}
func (e Result) GetMsg() string {
	return e.Msg
}

func RegisterErrorCode(code int, msg string) (Result)  {
	if _,ok := _code[code]; ok {
		panic("code has been registered")
	}
	if msg == ""{
		panic("msg cannot be empty")
	}

	_code[code] = struct{}{}

	_msg[code] = msg

	return Result{
		Code: code,
		Msg: msg,
	}
}

func GetMsg(code int) string{
	return _msg[code]
}

var(
	OkReresult = RegisterErrorCode(SUCCESS,"ok")
	FailResult = RegisterErrorCode(FAIL,"fail")
	//常见错误
	ErrRequest  = RegisterErrorCode(9001, "请求参数格式错误")
	ErrDbOp     = RegisterErrorCode(9004, "数据库操作异常")
	ErrRedisOp  = RegisterErrorCode(9005, "Redis 操作异常")

	// 登录相关错误
	ErrPassword     = RegisterErrorCode(1002, "密码错误")
	ErrUserNotExist = RegisterErrorCode(1003, "该用户不存在")
	ErrOldPassword  = RegisterErrorCode(1010, "旧密码不正确")

	// jwt认证错误
	ErrTokenNotExist    = RegisterErrorCode(1201, "TOKEN 不存在，请重新登陆")
	ErrTokenRuntime     = RegisterErrorCode(1202, "TOKEN 已过期，请重新登陆")
	ErrTokenWrong       = RegisterErrorCode(1203, "TOKEN 不正确，请重新登陆")
	ErrTokenType        = RegisterErrorCode(1204, "TOKEN 格式错误，请重新登陆")
	ErrTokenCreate      = RegisterErrorCode(1205, "TOKEN 生成失败")
	ErrPermission       = RegisterErrorCode(1206, "权限不足")
	// 数据库相关错误

	// Tag and Category
	ErrCateHasArt = RegisterErrorCode(3003, "删除失败，分类下存在文章")

	// 上传或获取文件
	ErrFileReceive = RegisterErrorCode(9101, "文件接收失败")

	// 注册错误
)

