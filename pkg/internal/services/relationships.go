package services

import (
	"errors"
	"fmt"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"gorm.io/gorm"
)

func ListAllRelationship(user models.Account) ([]models.AccountRelationship, error) {
	var relationships []models.AccountRelationship
	if err := database.C.
		Where("account_id = ?", user.ID).
		Preload("Account").
		Preload("Related").
		Find(&relationships).Error; err != nil {
		return relationships, err
	}

	return relationships, nil
}

func ListRelationshipWithFilter(user models.Account, status models.RelationshipStatus) ([]models.AccountRelationship, error) {
	var relationships []models.AccountRelationship
	if err := database.C.
		Where("account_id = ? AND status = ?", user.ID, status).
		Preload("Account").
		Preload("Related").
		Find(&relationships).Error; err != nil {
		return relationships, err
	}

	return relationships, nil
}

func GetRelationship(otherId uint) (models.AccountRelationship, error) {
	var relationship models.AccountRelationship
	if err := database.C.
		Where(&models.AccountRelationship{AccountID: otherId}).
		Preload("Account").
		Preload("Related").
		First(&relationship).Error; err != nil {
		return relationship, err
	}

	return relationship, nil
}

func GetRelationWithTwoNode(userId, relatedId uint, noPreload ...bool) (models.AccountRelationship, error) {
	var tx *gorm.DB
	if len(noPreload) > 0 && noPreload[0] {
		tx = database.C
	} else {
		tx = database.C.Preload("Account").Preload("Related")
	}

	var relationship models.AccountRelationship
	if err := tx.
		Where(&models.AccountRelationship{AccountID: userId, RelatedID: relatedId}).
		First(&relationship).Error; err != nil {
		return relationship, err
	}

	return relationship, nil
}

func NewFriend(userA models.Account, userB models.Account, skipPending ...bool) (models.AccountRelationship, error) {
	relA := models.AccountRelationship{
		AccountID: userA.ID,
		RelatedID: userB.ID,
		Status:    models.RelationshipFriend,
	}
	relB := models.AccountRelationship{
		AccountID: userB.ID,
		RelatedID: userA.ID,
		Status:    models.RelationshipPending,
	}

	if len(skipPending) > 0 && skipPending[0] {
		relB.Status = models.RelationshipFriend
	}

	if userA.ID == userB.ID {
		return relA, fmt.Errorf("you cannot make friendship with yourself")
	} else if _, err := GetRelationWithTwoNode(userA.ID, userB.ID, true); err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return relA, fmt.Errorf("you already have a friendship with him or her")
	}

	if err := database.C.Save(&relA).Error; err != nil {
		return relA, err
	} else if err = database.C.Save(&relB).Error; err != nil {
		return relA, err
	} else {
		_ = NewNotification(models.Notification{
			Title:     fmt.Sprintf("New friend request from %s", userA.Name),
			Body:      fmt.Sprintf("You got a new friend request from %s. Go to your settings and decide how to deal it.", userA.Nick),
			AccountID: userB.ID,
		})
	}

	return relA, nil
}

func EditRelationship(relationship models.AccountRelationship) (models.AccountRelationship, error) {
	if err := database.C.Save(&relationship).Error; err != nil {
		return relationship, err
	}
	return relationship, nil
}

func DeleteRelationship(relationship models.AccountRelationship) error {
	if err := database.C.Delete(&relationship).Error; err != nil {
		return err
	}
	return nil
}
