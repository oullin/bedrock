// Package socialauth provides OAuth1 and OAuth2 social authentication,
//
// Quick start:
//
//	manager := socialauth.NewManager(r, session, map[string]socialauth.ProviderConfig{
//		"github": {
//			ClientID:     "your-client-id",
//			ClientSecret: "your-client-secret",
//			RedirectURL:  "https://example.com/callback/github",
//		},
//	})
//
//	// Redirect
//	url, err := manager.Driver("github").Redirect(ctx)
//	http.Redirect(w, r, url, http.StatusFound)
//
//	// Callback
//	user, err := manager.Driver("github").User(ctx)
//	fmt.Println(user.Email, user.Token)
//
// Testing:
//
//	manager.Fake("github", &socialauth.User{ID: "123", Email: "test@example.com"})
//	user, _ := manager.Driver("github").User(ctx)
package socialauth
