package internal

import "errors"

var _ error = (*ErrCode)(nil)

var errMap = map[ErrCode]string{
	ErrInternalServer: "internal server error",
}

// ErrCode https://developer.mozilla.org/zh-CN/docs/Web/API/CloseEvent
type ErrCode uint16

func (e ErrCode) Error() string {
	return errMap[e]
}

type XError struct {
	ECode ErrCode
	Err   error
}

func NewXError(eCode ErrCode, err error) *XError {
	return &XError{
		ECode: eCode,
		Err:   err,
	}
}

func (e *XError) Error() string {
	return e.Err.Error()
}

const (
	// CloseNormal 正常关闭; 无论为何目的而创建，该链接都已成功完成任务
	CloseNormal ErrCode = 1000

	// ErrCloseGoingAway 终端离开，可能因为服务端错误，也可能因为浏览器正从打开连接的页面跳转离开
	ErrCloseGoingAway ErrCode = 1001

	// ErrCloseProtocol 由于协议错误而中断连接
	ErrCloseProtocol ErrCode = 1002

	// ErrCloseUnSupported 由于接收到不允许的数据类型而断开连接 (如仅接收文本数据的终端接收到了二进制数据)
	ErrCloseUnSupported ErrCode = 1003

	// ErrCloseNoStatus 表示没有收到预期的状态码
	ErrCloseNoStatus ErrCode = 1005

	// ErrCloseAbNormal 用于期望收到状态码时连接非正常关闭 (也就是说，没有发送关闭帧)
	ErrCloseAbNormal ErrCode = 1006

	// ErrUnsupportedData 由于收到了格式不符的数据而断开连接 (如文本消息中包含了非 UTF-8 数据)
	ErrUnsupportedData = 1007

	// ErrPolicyViolation 由于收到不符合约定的数据而断开连接。这是一个通用状态码，用于不适合使用 1003 和 1009 状态码的场景
	ErrPolicyViolation ErrCode = 1008

	// ErrCloseTooLarge 由于收到过大的数据帧而断开连接
	ErrCloseTooLarge ErrCode = 1009

	// ErrMissingExtension 客户端期望服务器商定一个或多个拓展，但服务器没有处理，因此客户端断开连接
	ErrMissingExtension = 1010

	// ErrInternalServer 客户端由于遇到没有预料的情况阻止其完成请求，因此服务端断开连接
	ErrInternalServer ErrCode = 1011

	// ErrServiceRestart 服务器由于重启而断开连接
	ErrServiceRestart ErrCode = 1012

	// ErrTryAgainLater 服务器由于临时原因断开连接，如服务器过载因此断开一部分客户端连接
	ErrTryAgainLater ErrCode = 1013

	// ErrTlsHandshake 表示连接由于无法完成 TLS 握手而关闭 (例如无法验证服务器证书)
	ErrTlsHandshake ErrCode = 1015
)

var (
	ErrHandShake  = errors.New("handshake failed")
	ErrTextEncode = errors.New("invalid text encode, must be utf-8 encode")
)
