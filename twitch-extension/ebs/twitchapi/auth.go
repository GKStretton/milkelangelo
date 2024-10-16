package twitchapi

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// returns token used to send to extensions pubsub through Helix api.
func (t *connectedTwitchAPI) getBroadcastToken(dur time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(dur).UTC().Unix(),
		"user_id":    t.channelID,
		"role":       "external",
		"channel_id": t.channelID,
		"pubsub_perms": map[string]interface{}{
			"send": []string{"broadcast"},
		},
	})
	signedToken, err := token.SignedString(t.sharedSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
