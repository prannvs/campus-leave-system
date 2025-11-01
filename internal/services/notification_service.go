package services

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/models"
)

type NotificationService struct {
	cfg core.SMTPConfig
}

func NewNotificationService(cfg core.SMTPConfig) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) SendLeaveStatusNotification(leave *models.LeaveRequest) {
	go func() {
		subject := fmt.Sprintf("Leave Request %s", leave.Status)
		body := fmt.Sprintf(
			"Your leave request from %s to %s has been %s.",
			leave.StartDate.Format("2006-01-02"),
			leave.EndDate.Format("2006-01-02"),
			leave.Status,
		)

		to := leave.Student.Email
		if err := s.sendEmail(to, subject, body); err != nil {
			log.Printf("Failed to send leave status email: %v", err)
			return
		}
		log.Printf("Email sent to %s about leave status: %s", to, leave.Status)
	}()
}

func (s *NotificationService) ScheduleLeaveReminder(leave *models.LeaveRequest) {
	reminderTime := leave.StartDate.Add(-24 * time.Hour)
	delay := time.Until(reminderTime)

	if delay <= 0 {
		return
	}

	go func() {
		time.Sleep(delay)
		subject := "Upcoming Leave Reminder"
		body := fmt.Sprintf("Your leave starts tomorrow (%s). Please make necessary arrangements.",
			leave.StartDate.Format("2006-01-02"))

		if err := s.sendEmail(leave.Student.Email, subject, body); err != nil {
			log.Printf("Failed to send leave reminder: %v", err)
			return
		}
		log.Printf("Reminder email sent to %s", leave.Student.Email)
	}()
}

func (s *NotificationService) sendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)

	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	return smtp.SendMail(addr, auth, s.cfg.User, []string{to}, msg)
}
