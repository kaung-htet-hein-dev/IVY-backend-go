package constants

var (
	// Error
	ErrEmailAlreadyRegistered    = "Email is already registered."
	ErrUserNotFound              = "user not found"
	ErrUnauthorized              = "Unauthorized"
	ErrInvalidToken              = "Invalid Token"
	ErrErrorRecordNotFound       = "Record not found"
	ErrValidationFailed          = "Validation failed"
	ErrProjectAlreadyExist       = "Project already exists. Choose another name."
	ErrProjectNotFound           = "Project not found"
	ErrFailedUpdateProject       = "failed to update project"
	ErrInvalidProjectID          = "invalid project ID"
	ErrFailedDeleteProject       = "failed to delete project"
	ErrProjectDeleted            = "project deleted successfully"
	ErrDuplicatedData            = "Duplicated entry found."
	ErrProjectUnauthorizedAccess = "You are not authorized to update this project"
	// Success
	SuccUserRegistered  = "User registered successfully"
	SuccLoginSuccessful = "login successful"
	Successful          = "Successful"
	SuccProjectCreated  = "Project created"
	SuccProjectUpdated  = "Project updated successfully"
	SuccTaskCreated     = "Task created successfully"
)
