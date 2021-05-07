package send

import (
	"fmt"
	"os"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/notification"
)

func SendComment(postOwner, newUsername string, postID int) {
	notif := notification.Notification{
		Username: postOwner,
		Title:    "New comment!",
		Message:  fmt.Sprintf("%s commented on your post", newUsername),
		Link:     fmt.Sprintf("%spost/%d", os.Getenv("EMOTIVES_URL"), postID),
	}

	err := notification.SendEvent(os.Getenv("NOTIFICATION_URL"), os.Getenv("NOTIFICATION_AUTH"), notif)
	if err != nil {
		logger.Error(err)
	}
}

func SendLike(postOwner, newUsername string, postID int) {
	notif := notification.Notification{
		Username: postOwner,
		Title:    "New like!",
		Message:  fmt.Sprintf("%s commented on your post", newUsername),
		Link:     fmt.Sprintf("%spost/%d", os.Getenv("EMOTIVES_URL"), postID),
	}

	err := notification.SendEvent(os.Getenv("NOTIFICATION_URL"), os.Getenv("NOTIFICATION_AUTH"), notif)
	if err != nil {
		logger.Error(err)
	}
}
