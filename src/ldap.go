package main

import (
	"net/http"

	ldap "imuslab.com/arozos/mod/auth/ldap"
	prout "imuslab.com/arozos/mod/prouter"
)

func ldapInit() {
	//ldap
	ldapHandler := ldap.NewLdapHandler(authAgent, registerHandler, sysdb, permissionHandler, userHandler, nightlyManager, iconSystem)

	//add a entry to the system settings
	adminRouter := prout.NewModuleRouter(prout.RouterOption{
		ModuleName:  "System Setting",
		AdminOnly:   true,
		UserHandler: userHandler,
		DeniedHandler: func(w http.ResponseWriter, r *http.Request) {
			errorHandlePermissionDenied(w, r)
		},
	})
	registerSetting(settingModule{
		Name:         "LDAP",
		Desc:         "Allows external account access to system",
		IconPath:     "SystemAO/advance/img/small_icon.png",
		Group:        "Security",
		StartDir:     "SystemAO/advance/ldap.html",
		RequireAdmin: true,
	})

	adminRouter.HandleFunc("/system/auth/ldap/config/read", ldapHandler.ReadConfig)
	adminRouter.HandleFunc("/system/auth/ldap/config/write", ldapHandler.WriteConfig)
	adminRouter.HandleFunc("/system/auth/ldap/config/testConnection", ldapHandler.TestConnection)
	adminRouter.HandleFunc("/system/auth/ldap/config/syncorizeUser", ldapHandler.SynchronizeUser)

	//login interface and login handler
	http.HandleFunc("/system/auth/ldap/login", ldapHandler.HandleLogin)
	http.HandleFunc("/system/auth/ldap/setPassword", ldapHandler.HandleSetPassword)
	http.HandleFunc("/system/auth/ldap/newPassword", ldapHandler.HandleNewPasswordPage)
	http.HandleFunc("/ldapLogin.system", ldapHandler.HandleLoginPage)
	http.HandleFunc("/system/auth/ldap/checkldap", ldapHandler.HandleCheckLDAP)
}
