package internal

import "errors"

var _ error = (*ErrCode)(nil)

var errMap = map[ErrCode]string{
	ErrInternalServer: "internal server error",
}

type ErrCode uint16

func (e ErrCode) Error() string {
	return errMap[e]
}

// ErrInternalServer 内部服务错误，断开连接
const ErrInternalServer ErrCode = 10001

var ErrHandShake = errors.New("handshake failed")
