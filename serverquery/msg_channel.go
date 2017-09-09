package serverquery

import (
	"context"
)

// ChannelState contains transient information about a channel.
type ChannelState struct {
	// TotalClients lists the total number of clients in this channel.
	TotalClients int `serverquery:"total_clients"`
	// TotalClientsFamily TODO: ??
	TotalClientsFamily int `serverquery:"total_clients_family"`
}

// ChannelBasicInfo contains basic information about a channel.
type ChannelBasicInfo struct {
	// Id is the identifier of the channel.
	Id int `serverquery:"cid"`
	// ParentId is the identifier of the parent.
	ParentId int `serverquery:"pid"`
	// Order is the relative ordering of the channel.
	Order int `serverquery:"channel_order"`
	// Name is the name of the channel.
	Name string `serverquery:"channel_name"`
	// Topic is the topic of the channel.
	Topic string `serverquery:"channel_topic"`
	// IsDefault is set if the channel is the default channel.
	IsDefault bool `serverquery:"channel_flag_default"`
	// IsPassworded is set if the channel has a password.
	IsPassworded bool `serverquery:"channel_flag_password"`
	// IsPermanent is set if the channel is not temporary and is permanent.
	IsPermanent bool `serverquery:"channel_flag_permanent"`
	// IsSemiPermanent is set if the channel is semi-permanent.
	IsSemiPermanent bool `serverquery:"channel_flag_semi_permanent"`
	// MaxClients sets the maximum client count, or -1 for infinite.
	MaxClients int `serverquery:"channel_maxclients"`
	// MaxFamilyClients sets the maximum family client count, or -1 for infinite.
	// Note: the default is infinite.
	MaxFamilyClients int `serverquery:"channel_maxfamilyclients"`
}

// ChannelFullInfo contains all information about a channel.
type ChannelFullInfo struct {
	ChannelBasicInfo

	// Password is the channel password, encoded.
	Password string `serverquery:"channel_password"`
	// Codec is the ID of the codec in use.
	Codec int `serverquery:"channel_codec"`
	// CodecQuality is the quality between 1-10 of the codec.
	CodecQuery int `serverquery:"channel_codec_quality"`
	// CodecLatencyFactor is the latency factor of the codec.
	CodecLatencyFactor int `serverquery:"channel_codec_latency_factor"`
	// CodecIsUnencrypted is set if the codec is not encrypted.
	CodecIsUnencrypted bool `serverquery:"channel_codec_is_unencrypted"`
	// SecuritySalt is the channel security salt.
	SecuritySalt string `serverquery:"channel_security_salt"`
	// DeleteDelay is the delay deleting the channel.
	DeleteDelay int `serverquery:"channel_delete_delay"`
	// MaxClientsUnlimited is set when the channel has unlimited client count.
	MaxClientsUnlimited bool `serverquery:"channel_flag_maxclients_unlimited"`
	// MaxFamilyClientsUnlimited is set when the channel has unlimited family client count.
	MaxFamilyClientsUnlimited bool `serverquery:"channel_flag_maxfamilyclients_unlimited"`
	// MaxFamilyClientsInherited is set when the channel has inherited its parent max family clients count.
	MaxFamilyClientsInherited bool `serverquery:"channel_flag_maxfamilyclients_inherited"`
	// FilePath is the path to the channel files.
	FilePath string `serverquery:"channel_filepath"`
	// NeededTalkPower is the needed channel talk power.
	NeededTalkPower int `serverquery:"channel_needed_talk_power"`
	// ForcedSilence indicates the channel must stay silent.
	ForcedSilence bool `serverquery:"channel_forced_silence"`
	// PhoneticName is the phonetic name of the channel if set.
	PhoneticName string `serverquery:"channel_name_phonetic"`
	// IconId is the id of the icon for the channel.
	IconId int `serverquery:"channel_icon_id"`
	// IsPrivate is set if the channel is private.
	IsPrivate bool `serverquery:"channel_flag_private"`
	// SecondsEmpty is the number of seconds nobody has been in the channel, since the channel was last refreshed.
	SecondsEmpty int `serverquery:"seconds_empty"`
}

// ChannelListEntry is an entry in the channel list.
type ChannelListEntry struct {
	ChannelBasicInfo
	ChannelState
}

// GetChannelListCommand lists the channels in the server.
type GetChannelListCommand struct{}

// GetResponseType returns an instance of the response type.
func (c *GetChannelListCommand) GetResponseType() interface{} {
	return make([]*ChannelListEntry, 0)
}

// GetCommandName returns the name of the command.
func (c *GetChannelListCommand) GetCommandName() string {
	return "channellist -topic -flags -limits"
}

// GetChannelList returns the list of channels.
func (c *ServerQueryAPI) GetChannelList(ctx context.Context) ([]*ChannelListEntry, error) {
	i, err := c.ExecuteCommand(ctx, &GetChannelListCommand{})
	if err != nil {
		return nil, err
	}
	return i.([]*ChannelListEntry), nil
}

// GetChannelInfoCommand gets info about a specific channel.
type GetChannelInfoCommand struct {
	// Id is the id of the channel.
	Id int `serverquery:"cid"`
}

// GetChannelInfoResponse contains the data returned by channelinfo.
type GetChannelInfoResponse struct {
	ChannelFullInfo
	ChannelState
}

// GetResponseType returns an instance of the response type.
func (c *GetChannelInfoCommand) GetResponseType() interface{} {
	return &GetChannelInfoResponse{}
}

// GetCommandName returns the name of the command.
func (c *GetChannelInfoCommand) GetCommandName() string {
	return "channelinfo"
}

// GetChannelInfo returns information about a channel.
func (c *ServerQueryAPI) GetChannelInfo(ctx context.Context, channelID int) (*GetChannelInfoResponse, error) {
	i, err := c.ExecuteCommand(ctx, &GetChannelInfoCommand{Id: channelID})
	if err != nil {
		return nil, err
	}
	r := i.(*GetChannelInfoResponse)
	r.Id = channelID
	return r, nil
}
