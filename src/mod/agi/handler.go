package agi

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

//Handle AGI Exectuion Request with token, design for letting other web scripting language like php to interface with AGI
func (g *Gateway) HandleAgiExecutionRequestWithToken(w http.ResponseWriter, r *http.Request) {
	token, err := mv(r, "token", false)
	if err != nil {
		//Username not defined
		sendErrorResponse(w, "Token not defined or empty.")
		return
	}

	script, err := mv(r, "script", false)
	if err != nil {
		//Username not defined
		sendErrorResponse(w, "Script path not defined or empty.")
		return
	}

	//Try to get the username from token
	username, err := g.Option.UserHandler.GetAuthAgent().GetUsernameFromToken(token)
	if err != nil {
		//This token is not valid
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized (Token not valid)"))
		return
	}

	//Check if user exists and have access to the script
	targetUser, err := g.Option.UserHandler.GetUserInfoFromUsername(username)
	if err != nil {
		//This user not exists
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized (User not exists)"))
		return
	}

	scriptScope := ""
	allowAccess := checkUserAccessToScript(targetUser, script, scriptScope)
	if !allowAccess {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized (Permission Denied)"))
		return
	}

	//Get the content of the script
	scriptContentByte, err := ioutil.ReadFile(filepath.Join("./web/", script))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Script Not Found"))
		return
	}
	scriptContent := string(scriptContentByte)

	g.ExecuteAGIScript(scriptContent, script, scriptScope, w, r, targetUser)
}
