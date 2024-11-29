package app

type TwitchAPI interface {
	BroadcastExtensionData(data []byte) error
}
