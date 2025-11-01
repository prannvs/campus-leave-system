# Campus Leave & Attendance Management System
## Features

- **Role-Based Access Control**: Admin, Faculty, Warden, and Student roles
- **JWT Authentication**: Secure token-based authentication
- **Leave Management**: Apply, approve/reject leave requests with validation
- **Attendance Tracking**: Daily attendance marking and statistics
- **Analytics Dashboard**: Leave patterns, attendance trends, and reports
- **Notifications**: Automated notifications for leave status changes
- **RESTful API**: Clean, documented API endpoints

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/prannvs/campus-leave-system.git
cd campus-leave-system
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment Variables

Create a `.env` file:

```env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=campus_leave_db
DB_SSLMODE=disable

JWT_SECRET="(fill later)"
JWT_EXPIRY=24h
```

### 4. Run with Docker

```bash
docker-compose up -d
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication

#### Register
POST /api/v1/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "student",
  "dept": "Computer Science",
  "hostel": "Block A"
}

#### Login
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "student"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}

### Leave Management

#### Apply for Leave (Student)
POST /api/v1/leaves/apply
Authorization: Bearer <token>

{
  "leave_type": "Medical",
  "reason": "High fever and need rest",
  "start_date": "2025-11-01",
  "end_date": "2025-11-03"
}

#### Get My Leaves (Student)
GET /api/v1/leaves/my
Authorization: Bearer <token>

#### Get Pending Leaves (Faculty/Warden)
GET /api/v1/leaves/pending
Authorization: Bearer <token>

#### Approve/Reject Leave (Faculty/Warden)
PUT /api/v1/leaves/{id}/approve
Authorization: Bearer <token>

{
  "status": "approved",
  "remarks": "Approved as per medical certificate"
}

### Attendance

#### Mark Attendance (Faculty/Warden)
POST /api/v1/attendance/mark
Authorization: Bearer <token>

{
  "student_id": 1,
  "date": "2025-10-29",
  "present": true
}

#### Get Attendance Stats
GET /api/v1/attendance/stats?student_id=1&start_date=2025-10-01&end_date=2025-10-31
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "student_id": 1,
    "present_days": 22,
    "total_days": 25,
    "attendance_percentage": 88.0
  }
}

### Analytics (Admin Only)

#### Get Analytics Summary
GET /api/v1/analytics/summary?start_date=2025-10-01&end_date=2025-10-31
Authorization: Bearer <token>

#### Get Leave Type Breakdown
GET /api/v1/analytics/leave-breakdown
Authorization: Bearer <token>

## Security

- Passwords are hashed using bcrypt
- JWT tokens for stateless authentication
- Role-based access control on all protected routes
- Input validation on all endpoints
- SQL injection prevention via GORM ORM

## Database Schema

### Users Table
- id (Primary Key)
- name, email, password
- role (admin/faculty/warden/student)
- dept, hostel
- timestamps

### Leave Requests Table
- id (Primary Key)
- student_id (Foreign Key → users.id)
- leave_type, reason
- start_date, end_date
- status (pending/approved/rejected)
- approved_by (Foreign Key → users.id)
- remarks
- timestamps

### Attendance Table
- id (Primary Key)
- student_id (Foreign Key → users.id)
- date, present
- marked_by (Foreign Key → users.id)
- timestamps