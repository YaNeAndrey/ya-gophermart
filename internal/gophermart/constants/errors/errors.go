package error

import "errors"

var IncorrectPortNumber = errors.New("port accepts values from the range [1:65535]")
var IncorrectEndpointFormat = errors.New("need address in a form host:port")
