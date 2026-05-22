package reverb_test

import (
	"github.com/bedrock/packages/websockets"
)

func signedChannelAuth(app *websockets.App, socketID, channel, channelData string) string {
	return app.Key() + ":" + websockets.SignChannel(app.Secret(), socketID, channel, channelData)
}
