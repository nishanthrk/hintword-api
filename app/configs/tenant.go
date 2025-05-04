package configs

var (
	TenantConfig *TenantConfiguration
)

type TenantConfiguration struct {
	TenantCode           string `json:"TENANT_CODE"`
	TenantCodePrefix     string `json:"TENANT_CODE_PREFIX"`
	TenantFavicon        string `json:"TENANT_FAVICON"`
	TenantIcon           string `json:"TENANT_ICON"`
	TenantLogo           string `json:"TENANT_LOGO"`
	TenantName           string `json:"TENANT_NAME"`
	TenantPrimaryColor   string `json:"TENANT_PRIMARY_COLOR"`
	TenantSecondaryColor string `json:"TENANT_SECONDARY_COLOR"`
	TenantType           string `json:"TENANT_TYPE"`
}

//func LoadTenantConfiguration() {
//
//	TenantConfig = &TenantConfiguration{
//		TenantCode:           data["TENANT_CODE"],
//		TenantCodePrefix:     data["TENANT_CODE_PREFIX"],
//		TenantFavicon:        data["TENANT_FAVICON"],
//		TenantIcon:           data["TENANT_ICON"],
//		TenantLogo:           data["TENANT_LOGO"],
//		TenantName:           data["TENANT_NAME"],
//		TenantPrimaryColor:   data["TENANT_PRIMARY_COLOR"],
//		TenantSecondaryColor: data["TENANT_SECONDARY_COLOR"],
//		TenantType:           data["TENANT_TYPE"],
//	}
//}
