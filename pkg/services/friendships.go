package services

import (
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/identity/pkg/database"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"gorm.io/gorm"
)

func ListFriend(anyside models.Account, status models.FriendshipStatus) ([]models.AccountFriendship, error) {
	var relationships []models.AccountFriendship
	if err := database.C.
		Where("(account_id = ? OR related_id = ?) AND status = ?", anyside.ID, anyside.ID, status).
		Preload("Account").
		Preload("Related").
		Find(&relationships).Error; err != nil {
		return relationships, err
	}

	return relationships, nil
}

func GetFriend(anysideId uint) (models.AccountFriendship, error) {
	var relationship models.AccountFriendship
	if err := database.C.
		Where(&models.AccountFriendship{AccountID: anysideId}).
		Or(&models.AccountFriendship{RelatedID: anysideId}).
		Preload("Account").
		Preload("Related").
		First(&relationship).Error; err != nil {
		return relationship, err
	}

	return relationship, nil
}

func GetFriendWithTwoSides(userId, relatedId uint, noPreload ...bool) (models.AccountFriendship, error) {
	var tx *gorm.DB
	if len(noPreload) > 0 && noPreload[0] {
		tx = database.C
	} else {
		tx = database.C.Preload("Account").Preload("Related")
	}

	var relationship models.AccountFriendship
	if err := tx.
		Where(&models.AccountFriendship{AccountID: userId, RelatedID: relatedId}).
		Or(&models.AccountFriendship{RelatedID: userId, AccountID: relatedId}).
		First(&relationship).Error; err != nil {
		return relationship, err
	}

	return relationship, nil
}

func NewFriend(user models.Account, related models.Account, status models.FriendshipStatus) (models.AccountFriendship, error) {
	relationship := models.AccountFriendship{
		AccountID: user.ID,
		RelatedID: related.ID,
		Status:    status,
	}

	if user.ID == related.ID {
		return relationship, fmt.Errorf("you cannot make friendship with yourself")
	} else if _, err := GetFriendWithTwoSides(user.ID, related.ID, true); err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return relationship, fmt.Errorf("you already have a friendship with him or her")
	}

	if err := database.C.Save(&relationship).Error; err != nil {
		return relationship, err
	} else {
		_ = NewNotification(models.Notification{
			Subject:     fmt.Sprintf("New friend request from %s", user.Name),
			Content:     fmt.Sprintf("You got a new friend request from %s. Go to your settings and decide how to deal it.", user.Nick),
			RecipientID: related.ID,
		})
	}

	return relationship, nil
}

func EditFriendWithCheck(relationship models.AccountFriendship, user models.Account, originalStatus models.FriendshipStatus) (models.AccountFriendship, error) {
	if relationship.Status != originalStatus {
		if originalStatus == models.FriendshipBlocked && relationship.BlockedBy != nil && user.ID != *relationship.BlockedBy {
			return relationship, fmt.Errorf("the friendship has been blocked by the otherside, you cannot modify it status")
		}
		if relationship.Status == models.FriendshipPending && relationship.RelatedID != user.ID {
			return relationship, fmt.Errorf("only related person can accept friendship")
		}
	}
	if originalStatus != models.FriendshipBlocked && relationship.Status == models.FriendshipBlocked {
		relationship.BlockedBy = &user.ID
	}

	return EditFriend(relationship)
}

func EditFriend(relationship models.AccountFriendship) (models.AccountFriendship, error) {
	if err := database.C.Save(&relationship).Error; err != nil {
		return relationship, err
	}
	return relationship, nil
}

func DeleteFriend(relationship models.AccountFriendship) error {
	if err := database.C.Delete(&relationship).Error; err != nil {
		return err
	}
	return nil
}
