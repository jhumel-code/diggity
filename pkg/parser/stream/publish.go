package stream

func Publish(payload interface{}, channel string) error {
	err := tr.SendBroadcastMessage(channel, payload)
	if err != nil {
		return err
	}
	return nil
}
