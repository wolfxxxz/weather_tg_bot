package apperrors

import "fmt"

type AppError struct {
	Message string
	Code    string
}

func NewAppError() *AppError {
	return &AppError{}
}

var (
	EnvConfigLoadError = AppError{
		Message: "Failed to load env file",
		Code:    EnvInitErr,
	}
	EnvConfigParseError = AppError{
		Message: "Failed to parse env file",
		Code:    EnvParseErr,
	}
	BotInitializationError = AppError{
		Message: "Failed to init new bot",
		Code:    BotInitErr,
	}
	BotSendMessageError = AppError{
		Message: "Failed to send the message",
		Code:    BotSendMsgErr,
	}
	WeatherClientError = AppError{
		Message: "Failed send request",
		Code:    HttpSendRequestErr,
	}
	BotMapperEncodingErr = AppError{
		Message: "Failed encoding mapper",
		Code:    MapperEncodingErr,
	}
	MongoDataExistsError = AppError{
		Message: "Failed send mongoDB",
		Code:    UserRepoErr,
	}
	MongoSaveUserFailedError = AppError{
		Message: "Failed save mongoDB",
		Code:    UserRepoErr,
	}
	MongoGetFailedError = AppError{
		Message: "Failed Get mongoDB",
		Code:    UserRepoErr,
	}
	MongoUpdateModFailedError = AppError{
		Message: "Failed Update mongoDB",
		Code:    UserRepoErr,
	}
	MongoInitFailedError = AppError{
		Message: "Failed Init mongoDB",
		Code:    InitMongoErr,
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message: fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:    appError.Code,
	}
}

func IsAppError(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
