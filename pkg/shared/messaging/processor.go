package messaging

import "project-structure/internal/user"

type EventProcessor struct {
	UserService *user.Service
}

func NewEventProcessor(userService *user.Service) *EventProcessor {
	return &EventProcessor{
		UserService: userService,
	}
}
