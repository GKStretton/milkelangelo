package twitchextension

type ListeningData struct{}

func (e *extensionSession) SubscribeData() (chan *ListeningData, error) {
	ch := make(chan *ListeningData)
	// todo: implement open http sse with ebs
	// todo: implement closing
	return ch, nil
}

func (e *extensionSession) UnsubscribeData(ch chan *ListeningData) {
}
