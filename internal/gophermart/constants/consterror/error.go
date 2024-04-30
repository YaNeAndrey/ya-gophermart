package consterror

import "errors"

var IncorrectPortNumber = errors.New("port accepts values from the range [1:65535]")
var IncorrectEndpointFormat = errors.New("need address in a form host:port")

var DuplicateLogin = errors.New("login is already registered")
var LoginNotFound = errors.New("login not registered")
var DuplicateUserOrder = errors.New("you have already registered an order with this number")
var DuplicateAnotherUserOrder = errors.New("another user registered an order with the same number")
var OrderNotFound = errors.New("order not found")
var InsufficientFunds = errors.New("insufficient funds")
var CountRequestToAccrual = errors.New("the number of requests to accrual service has been exceeded")
