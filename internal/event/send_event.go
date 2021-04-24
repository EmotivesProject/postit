package event

import (
	"postit/model"

	commonKafka "github.com/TomBowyerResearchProject/common/kafka"
)

func SendPostEvent(username, status string, post *model.Post) {
	eventData := model.EventData{
		Data:   post,
		Status: status,
	}
	event := model.Event{
		Username:      username,
		CustomerEvent: model.EventPost,
		Data:          eventData,
	}
	commonKafka.ProduceEvent(event)
}

func SendLikeEvent(username, status string, like *model.Like) {
	eventData := model.EventData{
		Data:   like,
		Status: status,
	}
	event := model.Event{
		Username:      username,
		CustomerEvent: model.EventLike,
		Data:          eventData,
	}
	commonKafka.ProduceEvent(event)
}

func SendCommentEvent(username, status string, comment *model.Comment) {
	eventData := model.EventData{
		Data:   comment,
		Status: status,
	}
	event := model.Event{
		Username:      username,
		CustomerEvent: model.EventComment,
		Data:          eventData,
	}
	commonKafka.ProduceEvent(event)
}
