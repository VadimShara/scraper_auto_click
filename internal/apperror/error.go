package apperror

import "encoding/json"

var(
	EndOfCatalogue = NewAppError(nil, "The end of catalogue of site", "", "US-000001")
	ErrorRetrievingText = NewAppError(nil, "The element text could not be retrieved", "", "US-000002")
	ErrorCreationObjectDB = NewAppError(nil, "Database object creation error", "", "US-000003")
	ErrorElementHTMLSearch = NewAppError(nil, "Element HTML not found", "", "US-000004")
	ErrorOpeningSitePage = NewAppError(nil, "Site page does not open", "", "US-000005")
	ErrorDB = NewAppError(nil, "Error in database", "", "US-000006")
)

type AppError struct {
	Err					error 
	Message				string
	DeveloperMessage	string
	Code				string
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError{
	return &AppError{
		Err: err,
		Message: message, 
		DeveloperMessage: developerMessage,
		Code: code,
	}
}
