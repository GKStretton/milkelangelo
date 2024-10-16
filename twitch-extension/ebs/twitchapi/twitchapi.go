package twitchapi

type TwitchAPI interface {
	BroadcastExtensionData(d *BroadcastData) error
}
