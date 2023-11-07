package msg

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"gorm.io/gorm"

	"github.com/vincadrn/capstone/models"
)

func SetFCMClientToken(db *gorm.DB, deviceName string, token string, willBeNotified bool) error {
	var clientToken models.FCMClientToken

	tx := db.Begin()
	if tx.Where("device_name = ?", deviceName).First(&clientToken); tx.Error != nil {
		tx.Rollback()
		return errors.New("read query failed")
	}
	
	if clientToken.DeviceName == "" {
		clientToken.DeviceName = deviceName	
	}
	clientToken.Token = token
	clientToken.WillBeNotified = willBeNotified
	
	if tx.Save(&clientToken); tx.Error != nil {
		tx.Rollback()
		return errors.New("update query failed")
	}
	tx.Commit()
	return nil
}

func ListFCMClientTokens(db *gorm.DB) []string {
	var clientTokens []models.FCMClientToken
	
	db.Where("will_be_notified = ?", true).Find(&clientTokens)
	
	var tokens []string
	for _, clientToken := range clientTokens {
		tokens = append(tokens, clientToken.Token)
	}
	return tokens
}

func NotifyArrivingBus(db *gorm.DB, client *messaging.Client, ctx context.Context) {
	tokens := ListFCMClientTokens(db)
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Bus Stop",
			Body: "Bus has arrived!",
		},
		Tokens: tokens,
	}

	br, err := client.SendEachForMulticast(ctx, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully sent arriving bus notification:", br)
}
