package reverb_test

import (
	"github.com/bedrock/packages/reverb"
)

func signedChannelAuth(app *reverb.App, socketID, channel, channelData string) string {
	return app.Key() + ":" + reverb.SignChannel(app.Secret(), socketID, channel, channelData)
}
