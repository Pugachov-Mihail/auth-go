package base_validate

import authServer "auth/protos/gen/dota_traker.auth.v1"

func ValidateLoginRequest(value *authServer.AuthLoginRequest) bool {
	if value.GetLogin() != "" && value.GetPassword() != "" {
		return true
	}
	return false
}

func ValidatePassword(value *authServer.AuthRegistrationRequest) bool {
	if value.GetPassword() == value.GetPassword2() {
		return true
	}
	return false
}
