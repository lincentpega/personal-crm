package services

import (
	"context"
	"fmt"
	"time"

	"github.com/lincentpega/personal-crm/internal/config"
	"github.com/lincentpega/personal-crm/internal/log"
	"github.com/lincentpega/personal-crm/internal/models/notifications"
	"github.com/lincentpega/personal-crm/internal/models/person"
	"gopkg.in/telebot.v3"
)

type NotificationService struct {
	bot               *telebot.Bot
	notificationsRepo *notifications.NotificationRepository
	personRepo        *person.PersonRepository
	log               *log.Logger
	config            *config.AppConfig
}

func NewNotificationService(bot *telebot.Bot, notificationsRepo *notifications.NotificationRepository,
	personRepo *person.PersonRepository, log *log.Logger, config *config.AppConfig) *NotificationService {
	return &NotificationService{
		bot:               bot,
		notificationsRepo: notificationsRepo,
		personRepo:        personRepo,
		log:               log,
		config:            config,
	}
}

func (s *NotificationService) ProcessNotifications(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.InfoLog.Print("Stopping schedulled notifications processing")
			return
		case <-ticker.C:
			s.execProcessNotifications(ctx)
		}
	}
}

func (s *NotificationService) execProcessNotifications(ctx context.Context) error {
	ns, err := s.notificationsRepo.GetAwaitingSend(ctx)
	if err != nil {
		return err
	}

	for _, n := range ns {
		switch n.Type {
		case notifications.KeepInTouch:
			go s.processKeepInTouch(ctx, &n)
		default:
			s.log.ErrorLog.Print("Notification type is not defined yet")
		}
	}

	return nil
}

func (s *NotificationService) processKeepInTouch(ctx context.Context, n *notifications.Notification) {
	person, err := s.personRepo.Get(ctx, n.PersonID)
	if err != nil {
		s.log.ErrorLog.Printf("Failed to load person: personID %d. Notification moved to failed", n.PersonID)
		s.failNotification(ctx, n)
		return
	}

	msg := fmt.Sprintf("It's time to contact with %s %s", person.FirstName, *person.LastName)
	_, err = s.bot.Send(telebot.ChatID(s.config.UserID), msg)
	if err != nil {
		s.failNotification(ctx, n)
	}

	s.markNotificationsRaised(ctx, n)
}

func (s *NotificationService) failNotification(ctx context.Context, n *notifications.Notification) {
	err := s.notificationsRepo.UpdateNotificationStatus(ctx, n.ID, notifications.Failed)
	if err != nil {
		s.log.ErrorLog.Print("Failed to mark notification as failed")
	}
}

func (s *NotificationService) markNotificationsRaised(ctx context.Context, n *notifications.Notification) {
	err := s.notificationsRepo.UpdateNotificationStatus(ctx, n.ID, notifications.Raised)
	if err != nil {
		s.log.ErrorLog.Print("Failed to mark notification as failed")
	}
}
