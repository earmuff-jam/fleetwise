package config

const (
	// Test and Token Validity
	CTO_USER                 = "community_test"
	CEO_USER                 = "ceo_test"
	DefaultTokenValidityTime = "15"
)

const (
	// Email Service Notification
	ResetPasswordTokenStringURI     = "/api/v1/reset"
	EmailVerificationTokenStringURI = "/api/v1/verify"
	ResetPasswordSubjectLine        = "Reset your password for FleetWise Application"
	ResetPasswordPlainTextMessage   = "Enter the OTP to reset your password"
	VerifyEmailSubjectLine          = "Verify your email address for FleetWise Application"
	VerifyEmailPlainTextMessage     = "Click on the following link to verify your email address"
)

const (
	// Errors
	ErrorGeneratedOTPFailure     = "invalid otp detected"
	ErrorInvalidUserID           = "invalid user id"
	ErrorInvalidApiKey           = "invalid api key"
	ErrorUnableToSendEmail       = "system down. unable to send email"
	ErrorUserCredentialsNotFound = "user credentials are not configured"
	ErrorWebApplicationEndpoint  = "email endpoint not found"
	ErrorTokenValidation         = "unable to validate token"
	ErrorTokenSubjectValidation  = "unable to validate token subject"
	ErrorFetchingCurrentUser     = "unable to retrieve system user"
	ErrorUserIsAlreadyVerified   = "unable to validate user. user is already verified"
)
