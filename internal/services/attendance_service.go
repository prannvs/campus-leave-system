package services

import (
	"time"

	"github.com/prannvs/campus-leave-system/internal/models"
	"github.com/prannvs/campus-leave-system/internal/repositories"
)

type AttendanceService struct {
	repo *repositories.AttendanceRepository
}

func NewAttendanceService(repo *repositories.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) MarkAttendance(studentID uint, date time.Time, present bool, markedBy uint) error {
	// Check if attendance already exists
	existing, _ := s.repo.FindByStudentAndDate(studentID, date)
	if existing != nil {
		return models.ErrAttendanceExists
	}

	attendance := &models.Attendance{
		StudentID: studentID,
		Date:      date,
		Present:   present,
		MarkedBy:  markedBy,
	}

	return s.repo.Create(attendance)
}

func (s *AttendanceService) GetStats(studentID uint, startDate, endDate time.Time) (*models.AttendanceStats, error) {
	return s.repo.GetStats(studentID, startDate, endDate)
}

func (s *AttendanceService) GetLowAttendanceStudents(threshold float64) ([]map[string]interface{}, error) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0) // Last month
	return s.repo.GetLowAttendanceStudents(threshold, startDate, now)
}
