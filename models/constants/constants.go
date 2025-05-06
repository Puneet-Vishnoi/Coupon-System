package constants

const (
	APIVersion = "v1"
	MAX_DB_ATTEMPTS = 1

	BadRequestMessage = "request not fulfilled"

	// schedular constants
	HealthCheckRoute = "/health"
	MDBUri           = "mongodb://localhost:27017"
	Database         = "ecommerce"
	Sender           = "puneetvishnoiias@gmail.com"
	RedisAddr = "localhost:6379"

	RedisPassword =""
	RedisDB=0

	// email verification routes
	VerifyEmailRoute = "/verify-email"
	VerifyOtpRoute   = "/verify-otp"
	ResendEmailRoute = "/resend-email"

	// user related routes
	UserRegisterRoute = "/user-register"
	UserLoginRoute    = "/login"

	// product routes

)

const (
	NormalUser = "user"
	AdminUser  = "admin"
)

const (
	// time slot for otp validation
	OtpValidation = 60
)

// collections
const (
	VerificationsCollection = "verifications"
	PersonnelCollection = "personnels"
	AssignOrderCollection = "assigned_ordered"
)

// messages
const (
	AlreadyRegisterWithThisEmail = "already register with this email"
	EmailIsNotVerified           = "your email is not verified please verify your email"
	EmailValidationError         = "wrong email passed"
	OtpValidationError           = "wrong otp passed"
	OtpExpiredValidationError    = "otp expired"
	AlreadyVerifiedError         = "already verified"
	OptAlreadySentError          = "otp already sent to email"
	NotRegisteredUser            = "you are not register user"
	PasswordNotMatchedError      = "password doesn't match"
	NotAuthorizedUserError       = "you are not authorized to do this"
	NoProductAvaliable           = "no product avaliable"
	UserDoesNotExists            = "user not exists"
	AddressNotExists             = "address not exists. please add one address"
)
