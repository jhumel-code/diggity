package stream

import "github.com/vmware/transport-go/bus"

const (
	PackageChannel = "package-channel"
	StatusChannel  = "status-channel"
	ErrorChannel   = "error-channel"
)

var (
	tr             = bus.GetBus()
	packageHandler bus.MessageHandler
	statusHandler  bus.MessageHandler
	errorHandler   bus.MessageHandler
)

func init() {
	tr.GetChannelManager().CreateChannel(PackageChannel)
	tr.GetChannelManager().CreateChannel(StatusChannel)
	tr.GetChannelManager().CreateChannel(ErrorChannel)
}