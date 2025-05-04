package constants

const (
	USER_TYPE_TENANT   = "TENANT"
	USER_TYPE_ADVISOR  = "advisor"
	USER_TYPE_EMPLOYEE = "EMPLOYEE"
	USER_TYPE_CHANNEL  = "CHANNEL"
	USER_TYPE_CUSTOMER = "CUSTOMER"
	UserTypeGuest      = "GUEST"

	GRAPHY_MEHTOD_E                  = "encrypt"
	WORKFLOW_TYPE_LENDER_INTEGRATION = "LENDER_INTEGRATION"

	APPLY_CAPACITY_ENTITY       = "ENTITY"
	APPLY_CAPACITY_PERSON       = "PERSON"
	APPLICANT_TYPE_PRIMARY      = "PRIMARY"
	APPLICANT_TYPE_CO_APPLICANT = "CO_APPLICANT"
	APPLICANT_TYPE_GUARANTOR    = "GUARANTOR"
	CloudStorageProviderS3      = "s3"

	OtpPartnerRegistration = "OTP_PARTNER_REGISTRATION"

	LoginWithOtp = "LOGIN_WITH_OTP"

	ResetPasswordLink = "RESET_PASSWORD_LINK"

	Msg91BaseUrl = "https://control.msg91.com/api/v5"

	ApplicationParticipantTypeSourcedBy = "SOURCED_BY"

	WorkflowStatusUpdateNone = 0
	LmsWorkflowStatusUpdate  = 1
	LosWorkflowStatusUpdate  = 2

	DirectIntegrateLeadPool string = "LEAD_POOL"

	DirectIntegrateLoanDetails string = "LOAN_DETAILS"
	DirectIntegrateDashboard   string = "DASHBOARD"

	EnvLoadMethodLocal = "LOCAL"
	EnvLoadMethodSSM   = "SSM"

	TaskTypeVerification = "VERIFICATION"

	AddressTypeOfficeAddress      = "OFFICE_ADDRESS"
	AddressTypePermanentResidence = "PERMANENT_RESIDENCE"
	AddressTypeCurrentResidence   = "CURRENT_RESIDENCE"
	AddressTypeRegisteredAddress  = "REGISTERED_ADDRESS"
	AddressTypeBusinessAddress    = "BUSINESS_ADDRESS"
	AddressTypeWarehouseAddress   = "WAREHOUSE_ADDRESS"

	VerificationTypeProfileVerification    = "PROFILE_VERIFICATION"
	VerificationTypePersonalDiscussion     = "PERSONAL_DISCUSSION"
	VerificationTypeResidenceVerification  = "RESIDENCE_VERIFICATION"
	VerificationTypeEmploymentVerification = "EMPLOYMENT_VERIFICATION"
	PartnerCategoryDsa                     = "DSA"
	PartnerCategoryDigital                 = "Digital Alliance"
	PartnerCategoryPos                     = "POS"
	PartnerCategoryConnector               = "Digital Connector"
	PartnerCategoryEcom                    = "Ecom"

	SystemPartnerPortal  = "PARTNER_PORTAL"
	SystemPartnerMobile  = "PARTNER_MOBILE"
	SystemEmployeePortal = "EMPLOYEE_PORTAL"
	SystemEmployeeMobile = "EMPLOYEE_MOBILE"

	LocationTypeCountry  = "COUNTRY"
	LocationTypeState    = "STATE"
	LocationTypeDistrict = "DISTRICT"

	DataAccessAll              = "ALL"
	DataAccessSubordinatesOnly = "SUBORDINATES_ONLY"

	ApplicationTypeFullFledged = "FULL_FLEDGED_APPLICATION"
	ApplicationTypeShort       = "SHORT_APPLICATION"

	WorkflowTypeLeadCreation      = "LEAD_CREATION"
	WorkflowTypePartnerOnboarding = "PARTNER_ONBOARDING"

	WorkflowStepTypeManual      = "MANUAL"
	WorkflowStepTypeConditional = "CONDITIONAL"
	WorkflowStepTypeAutomatic   = "AUTOMATIC"

	TaskStatusPending    = 0
	TaskStatusInProgress = 1
	TaskStatusCompleted  = 2
	TaskStatusFailed     = 3

	TaskStatusSkipped = 4
)
