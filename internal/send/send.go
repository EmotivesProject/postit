package send

import (
	"fmt"
	"os"

	"github.com/EmotivesProject/common/logger"
	"github.com/EmotivesProject/common/notification"
)

func Comment(postOwner, newUsername string, postID int) {
	notif := notification.Notification{
		Username:   postOwner,
		Type:       notification.Comment,
		Title:      "New comment!",
		Message:    fmt.Sprintf("%s commented on your post", newUsername),
		Link:       fmt.Sprintf("%spost/%d", os.Getenv("EMOTIVES_URL"), postID),
		PostID:     &postID,
		UsernameTo: &newUsername,
	}

	_, err := notification.SendEvent(os.Getenv("NOTIFICATION_URL"), os.Getenv("NOTIFICATION_AUTH"), notif)
	if err != nil {
		logger.Error(err)
	}
}

func Like(postOwner, newUsername string, postID int) {
	notif := notification.Notification{
		Username:   postOwner,
		Type:       notification.Like,
		Title:      "New like!",
		Message:    fmt.Sprintf("%s liked your post", newUsername),
		Link:       fmt.Sprintf("%spost/%d", os.Getenv("EMOTIVES_URL"), postID),
		PostID:     &postID,
		UsernameTo: &newUsername,
	}

	_, err := notification.SendEvent(os.Getenv("NOTIFICATION_URL"), os.Getenv("NOTIFICATION_AUTH"), notif)
	if err != nil {
		logger.Error(err)
	}
}

func DeletePost(postID int) {
	deleteURL := fmt.Sprintf("%s/post/%d", os.Getenv("NOTIFICATION_URL"), postID)

	_, err := notification.SendDelete(deleteURL, os.Getenv("NOTIFICATION_AUTH"))
	if err != nil {
		logger.Error(err)
	}
}

func DeleteLike(postID int, username string) {
	deleteURL := fmt.Sprintf("%s/like/post/%d/user/%s", os.Getenv("NOTIFICATION_URL"), postID, username)

	_, err := notification.SendDelete(deleteURL, os.Getenv("NOTIFICATION_AUTH"))
	if err != nil {
		logger.Error(err)
	}
}
