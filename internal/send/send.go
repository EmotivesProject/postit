package send

import (
	"fmt"
	"net/http"
	"os"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/notification"
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

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		logger.Error(err)

		return
	}

	req.Header.Set("Authorization", os.Getenv("NOTIFICATION_AUTH"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
	} else {
		resp.Body.Close()
	}
}
