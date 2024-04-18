package validator

import adminServer "auth/protos/gen/dota_traker.admin.v1"

func ValidatorAdminPermission(value *adminServer.AdminPermissionRequest) bool {
	if value.GetIsAdmin() {
		return true
	}
	if value.GetEmail() != "" {
		return true
	}
	return false
}

func ValidatorAdminSetting(value *adminServer.AdminSettiongsPanelRequest) bool {
	if value.GetEmail() != "" {
		return true
	}
	return false
}
