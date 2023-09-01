/*
Copyright © 2023 Daniel Wu <wxc@wxccs.org>
*/
package main

import "github.com/iceking2nd/rustdesk-api-server/cmd"

// @title rustdesk-api-server
// @version 0.0.1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Retrieving access token by login api

func main() {
	cmd.Execute()
}
