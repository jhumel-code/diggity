package stream

import "github.com/vmware/transport-go/bus"

func Subscribe(channel string) (bus.MessageHandler, error) {
	handler, err := tr.ListenStream(channel)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
