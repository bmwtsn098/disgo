package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake/v2"
)

var (
	_ Interaction = (*ComponentInteraction)(nil)
)

type ComponentInteraction struct {
	BaseInteraction
	Data    ComponentInteractionData `json:"data"`
	Message Message                  `json:"message"`
}

func (i *ComponentInteraction) UnmarshalJSON(data []byte) error {
	var baseInteraction baseInteractionImpl
	if err := json.Unmarshal(data, &baseInteraction); err != nil {
		return err
	}

	var interaction struct {
		Data    json.RawMessage `json:"data"`
		Message Message         `json:"message"`
	}
	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}

	var cType struct {
		Type ComponentType `json:"component_type"`
	}

	if err := json.Unmarshal(interaction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ComponentInteractionData
		err             error
	)
	switch cType.Type {
	case ComponentTypeButton:
		v := ButtonInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ComponentTypeStringSelectMenu:
		v := StringSelectMenuInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ComponentTypeUserSelectMenu:
		v := UserSelectMenuInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ComponentTypeRoleSelectMenu:
		v := RoleSelectMenuInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ComponentTypeMentionableSelectMenu:
		v := MentionableSelectMenuInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ComponentTypeChannelSelectMenu:
		v := ChannelSelectMenuInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	default:
		return fmt.Errorf("unknown component interaction data with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	i.BaseInteraction = baseInteraction

	i.Data = interactionData
	i.Message = interaction.Message
	i.Message.GuildID = baseInteraction.guildID
	return nil
}

func (ComponentInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (i ComponentInteraction) ButtonInteractionData() ButtonInteractionData {
	return i.Data.(ButtonInteractionData)
}

func (i ComponentInteraction) SelectMenuInteractionData() SelectMenuInteractionData {
	return i.Data.(SelectMenuInteractionData)
}

func (ComponentInteraction) interaction() {}

type ComponentInteractionData interface {
	Type() ComponentType
	CustomID() string

	componentInteractionData()
}

type rawButtonInteractionData struct {
	ComponentType ComponentType `json:"component_type"`
	Custom        string        `json:"custom_id"`
}

type ButtonInteractionData struct {
	customID string
}

func (d *ButtonInteractionData) UnmarshalJSON(data []byte) error {
	var v rawButtonInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.Custom
	return nil
}

func (d *ButtonInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawButtonInteractionData{
		Custom: d.customID,
	})
}

func (ButtonInteractionData) Type() ComponentType {
	return ComponentTypeButton
}

func (d ButtonInteractionData) CustomID() string {
	return d.customID
}

func (ButtonInteractionData) componentInteractionData() {}

type stringSelectMenuInteractionData struct {
	ComponentType ComponentType `json:"component_type"`
	CustomID      string        `json:"custom_id"`
	Values        []string      `json:"values"`
}

type snowflakeSelectMenuInteractionData struct {
	ComponentType ComponentType      `json:"component_type"`
	CustomID      string             `json:"custom_id"`
	Resolved      selectMenuResolved `json:"resolved"`
	Values        []snowflake.ID     `json:"values"`
}

type selectMenuResolved struct {
	Users    map[snowflake.ID]User            `json:"users"`
	Members  map[snowflake.ID]ResolvedMember  `json:"members"`
	Roles    map[snowflake.ID]Role            `json:"roles"`
	Channels map[snowflake.ID]ResolvedChannel `json:"channels"`
}

type SelectMenuInteractionData interface {
	ComponentInteractionData
	selectMenuInteractionData()
}

type StringSelectMenuInteractionData struct {
	customID string
	Values   []string
}

func (d *StringSelectMenuInteractionData) UnmarshalJSON(data []byte) error {
	var v stringSelectMenuInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.CustomID
	d.Values = v.Values
	return nil
}

func (d StringSelectMenuInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(stringSelectMenuInteractionData{
		ComponentType: d.Type(),
		CustomID:      d.customID,
		Values:        d.Values,
	})
}

func (StringSelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeStringSelectMenu
}

func (d StringSelectMenuInteractionData) CustomID() string {
	return d.customID
}

func (StringSelectMenuInteractionData) componentInteractionData()  {}
func (StringSelectMenuInteractionData) selectMenuInteractionData() {}

type UserSelectMenuInteractionData struct {
	customID string
	Resolved UserSelectMenuResolved `json:"resolved"`
	Values   []snowflake.ID         `json:"values"`
}

type UserSelectMenuResolved struct {
	Users   map[snowflake.ID]User           `json:"users"`
	Members map[snowflake.ID]ResolvedMember `json:"members"`
}

func (d *UserSelectMenuInteractionData) UnmarshalJSON(data []byte) error {
	var v snowflakeSelectMenuInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.CustomID
	d.Resolved = UserSelectMenuResolved{
		Users:   v.Resolved.Users,
		Members: v.Resolved.Members,
	}
	d.Values = v.Values
	return nil
}

func (d UserSelectMenuInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(snowflakeSelectMenuInteractionData{
		ComponentType: d.Type(),
		CustomID:      d.customID,
		Resolved: selectMenuResolved{
			Users:   d.Resolved.Users,
			Members: d.Resolved.Members,
		},
		Values: d.Values,
	})
}

func (d UserSelectMenuInteractionData) Users() []User {
	users := make([]User, 0, len(d.Resolved.Users))
	for _, userID := range d.Values {
		if user, ok := d.Resolved.Users[userID]; ok {
			users = append(users, user)
		}
	}
	return users
}

func (d UserSelectMenuInteractionData) Members() []ResolvedMember {
	members := make([]ResolvedMember, 0, len(d.Resolved.Members))
	for _, userID := range d.Values {
		if member, ok := d.Resolved.Members[userID]; ok {
			members = append(members, member)
		}
	}
	return members
}

func (UserSelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeUserSelectMenu
}

func (d UserSelectMenuInteractionData) CustomID() string {
	return d.customID
}

func (UserSelectMenuInteractionData) componentInteractionData()  {}
func (UserSelectMenuInteractionData) selectMenuInteractionData() {}

type RoleSelectMenuInteractionData struct {
	customID string
	Resolved RoleSelectMenuResolved `json:"resolved"`
	Values   []snowflake.ID         `json:"values"`
}

type RoleSelectMenuResolved struct {
	Roles map[snowflake.ID]Role `json:"roles"`
}

func (d *RoleSelectMenuInteractionData) UnmarshalJSON(data []byte) error {
	var v snowflakeSelectMenuInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.CustomID
	d.Resolved = RoleSelectMenuResolved{
		Roles: v.Resolved.Roles,
	}
	d.Values = v.Values
	return nil
}

func (d RoleSelectMenuInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(snowflakeSelectMenuInteractionData{
		ComponentType: d.Type(),
		CustomID:      d.customID,
		Resolved: selectMenuResolved{
			Roles: d.Resolved.Roles,
		},
		Values: d.Values,
	})
}

func (d RoleSelectMenuInteractionData) Roles() []Role {
	roles := make([]Role, 0, len(d.Values))
	for _, roleID := range d.Values {
		if role, ok := d.Resolved.Roles[roleID]; ok {
			roles = append(roles, role)
		}
	}
	return roles
}

func (RoleSelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeRoleSelectMenu
}

func (d RoleSelectMenuInteractionData) CustomID() string {
	return d.customID
}

func (RoleSelectMenuInteractionData) componentInteractionData()  {}
func (RoleSelectMenuInteractionData) selectMenuInteractionData() {}

type MentionableSelectMenuInteractionData struct {
	customID string
	Resolved MentionableSelectMenuResolved `json:"resolved"`
	Values   []snowflake.ID                `json:"values"`
}

type MentionableSelectMenuResolved struct {
	Users   map[snowflake.ID]User           `json:"users"`
	Members map[snowflake.ID]ResolvedMember `json:"members"`
	Roles   map[snowflake.ID]Role           `json:"roles"`
}

func (d *MentionableSelectMenuInteractionData) UnmarshalJSON(data []byte) error {
	var v snowflakeSelectMenuInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.CustomID
	d.Resolved = MentionableSelectMenuResolved{
		Users:   v.Resolved.Users,
		Members: v.Resolved.Members,
		Roles:   v.Resolved.Roles,
	}
	d.Values = v.Values
	return nil
}

func (d MentionableSelectMenuInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(snowflakeSelectMenuInteractionData{
		ComponentType: d.Type(),
		CustomID:      d.customID,
		Resolved: selectMenuResolved{
			Users:   d.Resolved.Users,
			Members: d.Resolved.Members,
			Roles:   d.Resolved.Roles,
		},
		Values: d.Values,
	})
}

func (d MentionableSelectMenuInteractionData) Users() []User {
	users := make([]User, 0, len(d.Resolved.Users))
	for _, userID := range d.Values {
		if user, ok := d.Resolved.Users[userID]; ok {
			users = append(users, user)
		}
	}
	return users
}

func (d MentionableSelectMenuInteractionData) Members() []ResolvedMember {
	members := make([]ResolvedMember, 0, len(d.Resolved.Members))
	for _, userID := range d.Values {
		if member, ok := d.Resolved.Members[userID]; ok {
			members = append(members, member)
		}
	}
	return members
}

func (d MentionableSelectMenuInteractionData) Roles() []Role {
	roles := make([]Role, 0, len(d.Resolved.Roles))
	for _, roleID := range d.Values {
		if role, ok := d.Resolved.Roles[roleID]; ok {
			roles = append(roles, role)
		}
	}
	return roles
}

func (MentionableSelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeMentionableSelectMenu
}

func (d MentionableSelectMenuInteractionData) CustomID() string {
	return d.customID
}

func (MentionableSelectMenuInteractionData) componentInteractionData()  {}
func (MentionableSelectMenuInteractionData) selectMenuInteractionData() {}

type ChannelSelectMenuInteractionData struct {
	customID string
	Resolved ChannelSelectMenuResolved `json:"resolved"`
	Values   []snowflake.ID            `json:"values"`
}

type ChannelSelectMenuResolved struct {
	Channels map[snowflake.ID]ResolvedChannel `json:"channels"`
}

func (d *ChannelSelectMenuInteractionData) UnmarshalJSON(data []byte) error {
	var v snowflakeSelectMenuInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.CustomID
	d.Resolved.Channels = v.Resolved.Channels
	d.Values = v.Values
	return nil
}

func (d ChannelSelectMenuInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(snowflakeSelectMenuInteractionData{
		ComponentType: d.Type(),
		CustomID:      d.customID,
		Resolved: selectMenuResolved{
			Channels: d.Resolved.Channels,
		},
		Values: d.Values,
	})
}

func (d ChannelSelectMenuInteractionData) Channels() []ResolvedChannel {
	channels := make([]ResolvedChannel, 0, len(d.Values))
	for _, channelID := range d.Values {
		if channel, ok := d.Resolved.Channels[channelID]; ok {
			channels = append(channels, channel)
		}
	}
	return channels
}

func (ChannelSelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeChannelSelectMenu
}

func (d ChannelSelectMenuInteractionData) CustomID() string {
	return d.customID
}

func (ChannelSelectMenuInteractionData) componentInteractionData()  {}
func (ChannelSelectMenuInteractionData) selectMenuInteractionData() {}
