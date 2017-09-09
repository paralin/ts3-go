package serverquery

import (
	"context"
)

// ServerGroupSummary is the summary of a server group.
type ServerGroupSummary struct {
	// ID is the server group ID.
	ID int `serverquery:"sgid"`
	// Name is the name of the server group.
	Name string `serverquery:"name"`
	// Type is the type of the server group.
	// 1 = server query, 2 = user
	Type int `serverquery:"type"`
	// IconId is the ID of the group icon.
	IconId int `serverquery:"iconid"`
	// SaveDb marks if the server group can save to the database.
	SaveDb bool `serverquery:"savedb"`
	// SortId ???
	SortId int `serverquery:"sortid"`
	// NameMode ???
	NameMode int `serverquery:"namemode"`
	// MemberModifyPower is the modify power of the group.
	MemberModifyPower int `serverquery:"n_member_modifyp"`
	// MemberAddPower is the add power of the group.
	MemberAddPower int `serverquery:"n_member_addp"`
	// MemberRemovePower is the member remove power of the group.
	MemberRemovePower int `serverquery:"n_member_removep"`
}

// GetServerGroupListCommand requests the server group list.
type GetServerGroupListCommand struct{}

// GetResponseType returns an instance of the response type.
func (c *GetServerGroupListCommand) GetResponseType() interface{} {
	return make([]*ServerGroupSummary, 0)
}

// GetCommandName returns the name of the command.
func (c *GetServerGroupListCommand) GetCommandName() string {
	return "servergrouplist"
}

// GetServerGroupList returns the list of server groups.
func (c *ServerQueryAPI) GetServerGroupList(ctx context.Context) ([]*ServerGroupSummary, error) {
	i, err := c.ExecuteCommand(ctx, &GetServerGroupListCommand{})
	if err != nil {
		return nil, err
	}
	return i.([]*ServerGroupSummary), nil
}

// ServerGroupAddClientCommand is the command to add a user to a server group.
type ServerGroupAddClientCommand struct {
	// ServerGroupID is the group to add
	ServerGroupID int `servequery:"sgid"`
	// ClientDBId is the client ID in the database.
	ClientDBId int `serverquery:"cldbid"`
}

// GetResponseType returns an instance of the response type.
func (c *ServerGroupAddClientCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *ServerGroupAddClientCommand) GetCommandName() string {
	return "servergroupaddclient"
}

// GetServerGroupList returns the list of server groups.
func (c *ServerQueryAPI) ServerGroupAddClient(ctx context.Context, clientDbId int, serverGroupId int) error {
	_, err := c.ExecuteCommand(ctx, &ServerGroupAddClientCommand{ClientDBId: clientDbId, ServerGroupID: serverGroupId})
	return err
}

// ServerGroupDelClientCommand is the command to delete a user from a server group.
type ServerGroupDelClientCommand struct {
	// ServerGroupID is the group to add
	ServerGroupID int `servequery:"sgid"`
	// ClientDBId is the client ID in the database.
	ClientDBId int `serverquery:"cldbid"`
}

// GetResponseType returns an instance of the response type.
func (c *ServerGroupDelClientCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *ServerGroupDelClientCommand) GetCommandName() string {
	return "servergroupdelclient"
}

// GetServerGroupList returns the list of server groups.
func (c *ServerQueryAPI) ServerGroupDelClient(ctx context.Context, clientDbId int, serverGroupId int) error {
	_, err := c.ExecuteCommand(ctx, &ServerGroupDelClientCommand{ClientDBId: clientDbId, ServerGroupID: serverGroupId})
	return err
}
