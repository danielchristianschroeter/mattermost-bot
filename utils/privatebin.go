package utils

import (
	"fmt"
	"mattermost-bot/confighandler"
	"net/url"
	"strconv"

	"github.com/gearnode/privatebin"
)

func PrivateBinPaste(msg string) (string, error) {
	if confighandler.App.Config.MB_PRIVATEBIN_ENABLE == "true" {
		uri, err := url.Parse(confighandler.App.Config.MB_PRIVATEBIN_URL)
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "PrivateBinPaste").Str("type", "response").Msg("Cannot parse host in MB_PRIVATEBIN_URL")
		}
		// Convert string value to bool
		opendiscussion, err := strconv.ParseBool(confighandler.App.Config.MB_PRIVATEBIN_OPENDISCUSSION)
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "PrivateBinPaste_ParseBool").Str("type", "response").Msg("Cannot parse bool in MB_PRIVATEBIN_OPENDISCUSSION")
		}
		burnafterreading, err := strconv.ParseBool(confighandler.App.Config.MB_PRIVATEBIN_BURNAFTERREADING)
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "PrivateBinPaste_ParseBool").Str("type", "response").Msg("Cannot parse bool in MB_PRIVATEBIN_BURNAFTERREADING")
		}
		client := privatebin.NewClient(uri, "", "")
		resp, err := client.CreatePaste(
			msg,
			confighandler.App.Config.MB_PRIVATEBIN_EXPIRE,
			confighandler.App.Config.MB_PRIVATEBIN_FORMATTER,
			opendiscussion,
			burnafterreading,
			confighandler.App.Config.MB_PRIVATEBIN_PASSWORD)
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "PrivateBinPaste_CreatePaste").Str("type", "response").Msg("Cannot create the paste")
			return "", err
		}
		fmt.Printf("%s\n", resp.URL)
		return resp.URL, nil
	}
	return "", nil
}
