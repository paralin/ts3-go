package serverquery

import (
	"context"
)

// ClientBasicInfo is the basic information given out in the client list.
type ClientBasicInfo struct {
	// Id is the ID of the client.
	Id int `serverquery:"clid"`
	// DatabaseId is the ID of the client in the database.
	DatabaseId int `serverquery:"client_database_id"`
	// Nickname is the nickname of the client.
	Nickname string `serverquery:"client_nickname"`
	// Type is the type of the client.
	Type int `serverquery:"client_type"`
	// UniqueIdentifier contains the client unique id.
	UniqueIdentifier string `serverquery:"client_unique_identifier"`
}

// GetClientListCommand requests the client list.
type GetClientListCommand struct{}

// GetResponseType returns an instance of the response type.
func (c *GetClientListCommand) GetResponseType() interface{} {
	return make([]*ClientBasicInfo, 0)
}

// GetCommandName returns the name of the command.
func (c *GetClientListCommand) GetCommandName() string {
	return "clientlist -uid"
}

// GetClientList returns the list of the clients.
func (c *ServerQueryAPI) GetClientList(ctx context.Context) ([]*ClientBasicInfo, error) {
	i, err := c.ExecuteCommand(ctx, &GetClientListCommand{})
	if err != nil {
		return nil, err
	}
	return i.([]*ClientBasicInfo), nil
}

// ClientInfo contains client information.
type ClientInfo struct {
	ClientBasicInfo

	// ClientIdleTime is the time the client has been idle.
	ClientIdleTime int `serverquery:"client_idle_time"`
	// Version is the client version.
	Version string `serverquery:"client_version"`
	// Platform is the client platform.
	Platform string `serverquery:"client_platform"`
	// InputMuted indicates if the client mic is muted.
	InputMuted bool `serverquery:"client_input_muted"`
	// OutputMuted indicates if the client speakers are muted.
	OutputMuted bool `serverquery:"client_output_muted"`
	// OutputOnlyMuted - ?
	OutputOnlyMuted bool `serverquery:"client_outputonly_muted"`
	// HasInputHardware indicates if the client has a mic.
	HasInputHardware bool `serverquery:"client_input_hardware"`
	// HasOutputHardware indicates if the client has a mic.
	HasOutputHardware bool `serverquery:"client_output_hardware"`
	// DefaultChannelName is the name of the client's default channel.
	DefaultChannelName string `serverquery:"client_default_channel"`
	// IsRecording indicates if the client is recording.
	IsRecording bool `serverquery:"client_is_recording"`
	// ChannelGroupId is the ID of the channel group the client is in.
	ChannelGroupId int `serverquery:"client_channel_group_id"`
	// ServerGroups is the list of server groups on the client.
	ServerGroups []int `serverquery:"client_servergroups"`
	// TotalConnections is the total number of times the client has connected.
	TotalConnections int `serverquery:"client_totalconnections"`
	// Away is set if the client is marked as away.
	Away bool `serverquery:"client_away"`
	// AwayMessage is the client away message.
	AwayMessage string `serverquery:"client_away_message"`
	// TalkPower is the talk power of the client.
	TalkPower int `serverquery:"client_talk_power"`
	// TalkRequest indicates if the client has requested to talk.
	TalkRequest bool `serverquery:"client_talk_request"`
	// TalkRequestMessage is the message the client gave when requesting to talk.
	TalkRequestMessage string `serverquery:"client_talk_request_msg"`
	// Description is the description of the client.
	Description string `serverquery:"client_description"`
	// IsTalker indicates if the client is a talker.
	IsTalker bool `serverquery:"client_is_talker"`
	// MonthBytesUploaded is the number of bytes uploaded this month.
	MonthBytesUploaded int `serverquery:"client_month_bytes_uploaded"`
	// MonthBytesDownloaded is the number of bytes downloaded this month.
	MonthBytesDownloaded int `serverquery:"client_month_bytes_downloaded"`
	// TotalBytesUploaded is the number of bytes uploaded total.
	TotalBytesUploaded int `serverquery:"client_total_bytes_uploaded"`
	// TotalBytesDownloaded is the number of bytes downloaded total.
	TotalBytesDownloaded int `serverquery:"client_total_bytes_downloaded"`
	// IsPrioritySpeaker indicates if the client is a priority speaker
	IsPrioritySpeaker bool `serverquery:"client_is_priority_speaker"`
	// PhoneticNickname is the phonetic nickname if given.
	PhoneticNickname string `serverquery:"client_nickname_phonetic"`
	// NeededServerQueryViewPower is the view power necessary to view the client.
	NeededServerQueryViewPower int `serverquery:"client_needed_serverquery_view_power"`
	// IconId is the icon ID of the client.
	IconId int `serverquery:"client_icon_id"`
	// IsChannelCommander indicates if the client is a channel commander.
	IsChannelCommander bool `serverquery:"is_channel_commander"`
}

// GetClientInfoCommand requests the info about a specific client.
type GetClientInfoCommand struct {
	// ClientId is the ID of the client.
	ClientId int `serverquery:"clid"`
}

// GetResponseType returns an instance of the response type.
func (c *GetClientInfoCommand) GetResponseType() interface{} {
	return &ClientInfo{}
}

// GetCommandName returns the name of the command.
func (c *GetClientInfoCommand) GetCommandName() string {
	return "clientinfo"
}

// GetClientInfo returns the info of a client.
func (c *ServerQueryAPI) GetClientInfo(ctx context.Context, clid int) (*ClientInfo, error) {
	i, err := c.ExecuteCommand(ctx, &GetClientInfoCommand{ClientId: clid})
	if err != nil {
		return nil, err
	}
	return i.(*ClientInfo), nil
}

// SendTextMessageCommand sends a text message.
type SendTextMessageCommand struct {
	// TargetMode is the mode of the target.
	TargetMode int `serverquery:"targetmode"`
	// Target is the target of the message
	Target int `serverquery:"target"`
	// Message is the message to send.
	Message string `serverquery:"msg"`
}

// GetResponseType returns an instance of the response type.
func (c *SendTextMessageCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *SendTextMessageCommand) GetCommandName() string {
	return "sendtextmessage"
}

// SendTextMessage sends a text message to a target.
func (c *ServerQueryAPI) SendTextMessage(
	ctx context.Context,
	targetType int,
	targetId int,
	message string,
) error {
	_, err := c.ExecuteCommand(ctx, &SendTextMessageCommand{
		TargetMode: targetType,
		Target:     targetId,
		Message:    message,
	})
	return err
}
