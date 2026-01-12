package constants

const (
	CodeSuccess                     = 1000
	CodeLoginSuccess                = 1001
	CodeLogoutSuccess               = 1002
	CodeChangePasswordSuccess       = 1003
	CodeForgotPasswordSuccess       = 1004
	CodeVerifyForgotPasswordSuccess = 1005
	CodeResetPasswordSuccess        = 1006
	CodeUpdateInfoSuccess           = 1007
	CodeCreateUserSuccess           = 1008
	CodeCreateDepartmentSuccess     = 1009
	CodeBadRequest                  = 4000
	CodeLoginFailed                 = 4001
	CodeInvalidToken                = 4002
	CodeUnAuth                      = 4003
	CodeNoRefreshToken              = 4004
	CodeUserNotFound                = 4005
	CodeInvalidPassword             = 4006
	CodeEmailDoesNotExist           = 4007
	CodeTooManyAttempts             = 4008
	CodeInvalidOTP                  = 4009
	CodeEmailAlreadyExists          = 4010
	CodePhoneAlreadyExists          = 4011
	CodeDepartmentNotFound          = 4012
	CodeUsernameAlreadyExists       = 4013
	CodeForbidden                   = 4014
	CodeNameAlreadyExists           = 4015
	CodeInvalidID                   = 4016
	CodeInternalError               = 5000

	ExchangeEmail       = "email.send"
	QueueNameAuthEmail  = "email.send.auth"
	RoutingKeyAuthEmail = "email.send.auth"
)
