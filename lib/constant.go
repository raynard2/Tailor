package lib

const (
	Limit          = 10
	MaxFileSize    = 100 // this is in mb
	INVALID_PARAMS = "Invalid Parameters."
	INVALID_BODY   = "Invalid Request Body."
	NOT_FOUND      = "Does Not Exist."
)

const (
	AccountRequested  string = "We have sent you a link to your email for completing signup."
	AccountExists     string = "Account with this email already exists."
	AccountNotExist   string = "Account with this email does not exist."
	WrongPassword            = "The password you entered is incorrect, Please check and try again."
	InvalidHash              = "You are using an invalid or expired link."
	PasswordResetSent string = "We have sent you a link with instructions to reset your password."
	AdminPrivilegeError string = "Admin Privilege Error Occurred"
	ErrorBinding	string	=	"Error Binding Params"
	DbError	string	= "Error Opening Database"
	InvalidUser string = "Invalid User"
	ErrorCreatingAccount string = "Internal Error Creating User, try again later"
	JwtAccountNotExist   string = "Could Not Get Email From Jwt Token."

)
