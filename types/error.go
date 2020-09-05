package types

import (
	"errors"
)

var (
	ErrHttpStatusNotOK             = errors.New("Http status is not 200")
	ErrNoSubDomain                 = errors.New("No subDomain record")
	ErrSubDomainNotMatch           = errors.New("RR.DomainName != subDomain")
	ErrResponseIdNotMatchRequestId = errors.New("Response id not match request id")

	ErrCannotGetIpFromIpCommnad = errors.New("Cannot get ip from ip command")
	ErrCannotGetIpFromIpSource  = errors.New("Cannot get ip from ip source")

	ErrNotNeedUpdate = errors.New("currentDNSIP == new ip,do nothing")
	ErrReqFieldEmpty = errors.New("some request field is empty")
)
