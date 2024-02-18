package services

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/database"
	"code.smartsheep.studio/hydrogen/identity/pkg/models"
	"fmt"
)

func GetThirdClient(id string) (models.ThirdClient, error) {
	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{
		Alias: id,
	}).First(&client).Error; err != nil {
		return client, err
	}

	return client, nil
}

func GetThirdClientWithSecret(id, secret string) (models.ThirdClient, error) {
	client, err := GetThirdClient(id)
	if err != nil {
		return client, err
	}

	if client.Secret != secret {
		return client, fmt.Errorf("invalid client secret")
	}

	return client, nil
}
