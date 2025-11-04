package schemas

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CreateFileRequest represents the JSON schema for creating a file
type CreateFileRequest struct {
	Filename string `json:"filename" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
}

// CreateTaskRequest represents the JSON schema for creating a task
type CreateTaskRequest struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	ParserID       uint   `json:"parser_id" validate:"required"`
	CronExpression string `json:"cron_expression"`
	IsRecurring    bool   `json:"is_recurring"`
	IsActive       bool   `json:"is_active"`
}

// CreateParserRequest represents the JSON schema for creating a parser
type CreateParserRequest struct {
	ParserName        string `json:"parser_name" validate:"required"`
	ParserDescription string `json:"parser_description"`
}
