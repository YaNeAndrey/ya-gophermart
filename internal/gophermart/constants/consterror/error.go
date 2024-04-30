package consterror

import "errors"

var ErrIncorrectPortNumber = errors.New("port accepts values from the range [1:65535]")
var ErrIncorrectEndpointFormat = errors.New("need address in a form host:port")

var ErrDuplicateLogin = errors.New("login is already registered")
var ErrLoginNotFound = errors.New("login not registered")
var ErrDuplicateUserOrder = errors.New("you have already registered an order with this number")
var ErrDuplicateAnotherUserOrder = errors.New("another user registered an order with the same number")
var ErrOrderNotFound = errors.New("order not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrCountRequestToAccrual = errors.New("the number of requests to accrual service has been exceeded")
