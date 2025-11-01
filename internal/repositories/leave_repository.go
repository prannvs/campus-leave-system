package repositories

import (
	"time"

	"github.com/prannvs/campus-leave-system/internal/models"
	"gorm.io/gorm"
)

type LeaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// internal/repositories/leave_repository.go

func (r *LeaveRepository) Create(leave *models.LeaveRequest) error {
	if err := r.db.Create(leave).Error; err != nil {
		return err
	}

	if err := r.db.Preload("Student").First(leave, leave.ID).Error; err != nil {
		return err
	}

	return nil
}

func (r *LeaveRepository) FindByID(id uint) (*models.LeaveRequest, error) {
	var leave models.LeaveRequest
	err := r.db.Preload("Student").Preload("Approver").First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *LeaveRepository) FindByStudentID(studentID uint) ([]models.LeaveRequest, error) {
	var leaves []models.LeaveRequest
	err := r.db.Where("student_id = ?", studentID).
		Preload("Student"). // ← Add this line!
		Preload("Approver").
		Order("created_at DESC").
		Find(&leaves).Error
	return leaves, err
}

func (r *LeaveRepository) FindPending() ([]models.LeaveRequest, error) {
	var leaves []models.LeaveRequest
	err := r.db.Where("status = ?", models.LeaveStatusPending).
		Preload("Student").
		Order("created_at ASC").
		Find(&leaves).Error
	return leaves, err
}

func (r *LeaveRepository) FindByStatus(status models.LeaveStatus, page, pageSize int) ([]models.LeaveRequest, int64, error) {
	var leaves []models.LeaveRequest
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.LeaveRequest{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Student").Preload("Approver").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&leaves).Error

	return leaves, total, err
}

func (r *LeaveRepository) Update(leave *models.LeaveRequest) error {
	return r.db.Save(leave).Error
}

func (r *LeaveRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeaveRequest{}, id).Error
}

func (r *LeaveRepository) CheckOverlapping(studentID uint, startDate, endDate time.Time, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.LeaveRequest{}).
		Where("student_id = ?", studentID).
		Where("status != ?", models.LeaveStatusRejected).
		Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?)",
			endDate, startDate, startDate, endDate)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *LeaveRepository) GetLeaveStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	var stats []struct {
		LeaveType string
		Count     int64
	}

	err := r.db.Model(&models.LeaveRequest{}).
		Select("leave_type, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("leave_type").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, stat := range stats {
		result[stat.LeaveType] = stat.Count
	}

	return result, nil
}
