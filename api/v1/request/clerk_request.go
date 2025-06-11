package request

type EmailAddress struct {
	EmailAddress string `json:"email_address"`
	ID           string `json:"id"`
	LinkedTo     []any  `json:"linked_to"`
	Object       string `json:"object"`
	Verification struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
	} `json:"verification"`
}

type ClerkUserData struct {
	Birthday              string         `json:"birthday"`
	CreatedAt             int64          `json:"created_at"`
	EmailAddresses        []EmailAddress `json:"email_addresses"`
	ExternalAccounts      []any          `json:"external_accounts"`
	ExternalID            string         `json:"external_id"`
	FirstName             string         `json:"first_name"`
	Gender                string         `json:"gender"`
	ID                    string         `json:"id"`
	ImageURL              string         `json:"image_url"`
	LastName              string         `json:"last_name"`
	LastSignInAt          int64          `json:"last_sign_in_at"`
	Object                string         `json:"object"`
	PasswordEnabled       bool           `json:"password_enabled"`
	PhoneNumbers          []string       `json:"phone_numbers"`
	PrimaryEmailAddressID string         `json:"primary_email_address_id"`
	PrimaryPhoneNumberID  any            `json:"primary_phone_number_id"`
	PrimaryWeb3WalletID   any            `json:"primary_web3_wallet_id"`
	PrivateMetadata       map[string]any `json:"private_metadata"`
	ProfileImageURL       string         `json:"profile_image_url"`
	PublicMetadata        map[string]any `json:"public_metadata"`
	TwoFactorEnabled      bool           `json:"two_factor_enabled"`
	UnsafeMetadata        map[string]any `json:"unsafe_metadata"`
	UpdatedAt             int64          `json:"updated_at"`
	Username              any            `json:"username"`
	Web3Wallets           []any          `json:"web3_wallets"`
}

type ClerkWebhookRequest struct {
	Data       ClerkUserData `json:"data"`
	InstanceID string        `json:"instance_id"`
	Object     string        `json:"object"`
	Timestamp  int64         `json:"timestamp"`
	Type       string        `json:"type"`
}
