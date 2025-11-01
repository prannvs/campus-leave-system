package repositories

import (
	"time"

	"github.com/prannvs/campus-leave-system/internal/models"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) FindByStudentAndDate(studentID uint, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Where("student_id = ? AND DATE(date) = DATE(?)", studentID, date).
		First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *AttendanceRepository) GetStats(studentID uint, startDate, endDate time.Time) (*models.AttendanceStats, error) {
	var presentDays, totalDays int64

	// Count total days
	r.db.Model(&models.Attendance{}).
		Where("student_id = ? AND date BETWEEN ? AND ?", studentID, startDate, endDate).
		Count(&totalDays)

	// Count present days
	r.db.Model(&models.Attendance{}).
		Where("student_id = ? AND date BETWEEN ? AND ? AND present = ?", studentID, startDate, endDate, true).
		Count(&presentDays)

	percentage := 0.0
	if totalDays > 0 {
		percentage = (float64(presentDays) / float64(totalDays)) * 100
	}

	return &models.AttendanceStats{
		StudentID:            studentID,
		PresentDays:          presentDays,
		TotalDays:            totalDays,
		AttendancePercentage: percentage,
	}, nil
}

func (r *AttendanceRepository) GetLowAttendanceStudents(threshold float64, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			u.id as student_id,
			u.name as student_name,
			u.dept,
			COUNT(*) as total_days,
			SUM(CASE WHEN a.present THEN 1 ELSE 0 END) as present_days,
			(SUM(CASE WHEN a.present THEN 1 ELSE 0 END)::float / COUNT(*)::float * 100) as attendance_percentage
		FROM users u
		INNER JOIN attendances a ON u.id = a.student_id
		WHERE u.role = 'student' 
			AND a.date BETWEEN ? AND ?
		GROUP BY u.id, u.name, u.dept
		HAVING (SUM(CASE WHEN a.present THEN 1 ELSE 0 END)::float / COUNT(*)::float * 100) < ?
		ORDER BY attendance_percentage ASC
		LIMIT 10
	`

	err := r.db.Raw(query, startDate, endDate, threshold).Scan(&results).Error
	return results, err
}

func (r *AttendanceRepository) BulkCreate(attendances []models.Attendance) error {
	return r.db.Create(&attendances).Error
}
