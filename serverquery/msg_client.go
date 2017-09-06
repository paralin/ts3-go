package serverquery

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
}

// ClientInfo contains client information.
type ClientInfo struct {
	ClientBasicInfo

	// UniqueIdentifier contains the client unique id.
	UniqueIdentifier string `serverquery:"client_unique_identifier`
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
	// ClientServergroups is the list of server groups on the client.
	ClientServergroups []int `serverquery:"client_servergroups"`
}
