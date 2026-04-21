package broadcasting

import contractsauth "github.com/bedrock/packages/contracts/auth"

func broadcastingIdentifier(user contractsauth.Authenticatable) string {
	if user == nil {
		return ""
	}

	if broadcastingUser, ok := user.(contractsauth.BroadcastingAuthenticatable); ok {
		return broadcastingUser.GetAuthIdentifierForBroadcasting()
	}

	return user.GetAuthIdentifier()
}
