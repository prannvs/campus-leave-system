package services

import (
	"time"

	"github.com/prannvs/campus-leave-system/internal/models"
	"github.com/prannvs/campus-leave-system/internal/repositories"
)

type LeaveService struct {
	leaveRepo       *repositories.LeaveRepository
	attendanceRepo  *repositories.AttendanceRepository
	notificationSvc *NotificationService
}

func NewLeaveService(
	leaveRepo *repositories.LeaveRepository,
	attendanceRepo *repositories.AttendanceRepository,
	notificationSvc *NotificationService,
) *LeaveService {
	return &LeaveService{
		leaveRepo:       leaveRepo,
		attendanceRepo:  attendanceRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *LeaveService) ApplyLeave(studentID uint, req models.ApplyLeaveRequest) (*models.LeaveRequest, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	leave := &models.LeaveRequest{
		StudentID: studentID,
		LeaveType: models.LeaveType(req.LeaveType),
		Reason:    req.Reason,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    models.LeaveStatusPending,
	}

	if err := leave.Validate(); err != nil {
		return nil, err
	}

	overlaps, err := s.leaveRepo.CheckOverlapping(studentID, startDate, endDate, 0)
	if err != nil {
		return nil, err
	}
	if overlaps {
		return nil, models.ErrOverlappingLeave
	}

	if err := s.leaveRepo.Create(leave); err != nil {
		return nil, err
	}

	return leave, nil
}

func (s *LeaveService) ApproveLeave(leaveID, approverID uint, status models.LeaveStatus, remarks *string) error {
	leave, err := s.leaveRepo.FindByID(leaveID)
	if err != nil {
		return models.ErrLeaveNotFound
	}

	if leave.Status != models.LeaveStatusPending {
		return models.ErrInvalidRole
	}

	leave.Status = status
	leave.ApprovedBy = &approverID
	leave.Remarks = remarks

	if err := s.leaveRepo.Update(leave); err != nil {
		return err
	}

	if status == models.LeaveStatusApproved {
		go s.markLeaveAttendance(leave, approverID)
	}

	go s.notificationSvc.SendLeaveStatusNotification(leave)

	return nil
}

func (s *LeaveService) markLeaveAttendance(leave *models.LeaveRequest, markerID uint) {
	currentDate := leave.StartDate
	for !currentDate.After(leave.EndDate) {
		attendance := &models.Attendance{
			StudentID: leave.StudentID,
			Date:      currentDate,
			Present:   false,
			MarkedBy:  markerID,
		}
		s.attendanceRepo.Create(attendance)
		currentDate = currentDate.AddDate(0, 0, 1)
	}
}

func (s *LeaveService) GetMyLeaves(studentID uint) ([]models.LeaveRequest, error) {
	return s.leaveRepo.FindByStudentID(studentID)
}

func (s *LeaveService) GetPendingLeaves() ([]models.LeaveRequest, error) {
	return s.leaveRepo.FindPending()
}

func (s *LeaveService) GetByStatus(status models.LeaveStatus, page, pageSize int) ([]models.LeaveRequest, int64, error) {
	return s.leaveRepo.FindByStatus(status, page, pageSize)
}

func (s *LeaveService) GetByID(id uint) (*models.LeaveRequest, error) {
	return s.leaveRepo.FindByID(id)
}

func (s *LeaveService) Delete(id uint) error {
	return s.leaveRepo.Delete(id)
}

func (s *LeaveService) GetLeaveStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.leaveRepo.GetLeaveStats(startDate, endDate)
}
