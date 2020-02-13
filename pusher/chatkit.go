package pusher

import (
	"context"

	chatkit "github.com/pusher/chatkit-server-go"
	log "github.com/sirupsen/logrus"
)

const (
	// todo move to config
	key     = "ebb1c9bd-127d-4626-9688-d80e7b4da191:1wSLxTDl771gNyjripyqEq3Op6tUeSRWzg6F+aBWvwQ="
	locator = "v1:us1:1ef279a7-c68d-4475-8139-4b142d2d21d3"
)

// ChatClient wrapps up chatkit sdk client
type ChatClient struct {
	sdk *chatkit.Client
}

// NewChatClient inits a new chat client
func NewChatClient() *ChatClient {
	c, err := chatkit.NewClient(locator, key)
	if err != nil {
		log.Error("cannot create new chatkit client")
		return nil
	}
	return &ChatClient{
		sdk: c,
	}
}

// CreaterUser creates a new chitkit user
func (c *ChatClient) CreaterUser(userId, username string, avatar *string) (err error) {
	data := chatkit.CreateUserOptions{
		ID:        userId,
		Name:      username,
		AvatarURL: avatar,
	}

	err = c.sdk.CreateUser(context.Background(), data)
	if err != nil {
		log.Error("Failed to create chatkit user:" + userId)
	}
	return
}

//CreateRoom creates chatkit chat room
func (c *ChatClient) CreateRoom(roomID string, isPrivate bool, users []string, creator string) (room string, err error) {
	data := chatkit.CreateRoomOptions{
		Name:      roomID,
		Private:   isPrivate,
		UserIDs:   users,
		CreatorID: creator,
	}

	r, err := c.sdk.CreateRoom(context.Background(), data)

	if err != nil {
		log.Error("failed to create chat room for:" + roomID)
		return
	}
	room = r.ID
	return
}

//AddUserToRoom adds a list of users to a specified room
func (c *ChatClient) AddUserToRoom(roomID string, users []string) (err error) {
	err = c.sdk.AddUsersToRoom(context.Background(), roomID, users)
	if err != nil {
		log.Error("failed to create chat room for:" + roomID)
	}
	return
}

func (c *ChatClient) SendSimpleMessage(roomID, senderID, content string) (ui uint, err error) {
	ui, err = c.sdk.SendSimpleMessage(context.Background(), chatkit.SendSimpleMessageOptions{
		RoomID:   roomID,
		SenderID: senderID,
		Text:     content,
	})
	if err != nil {
		log.Error(err.Error())
	}
	return
}
