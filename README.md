# Academix

A comprehensive Learning Management System (LMS) backend API built with Go, designed to manage courses, assignments, and submissions with role-based access control.

## 🚀 Features

### User Management
- **User Authentication**: JWT-based authentication with secure login/logout
- **Role-Based Access Control**: Three distinct roles (Student, Teacher, Admin)
- **Profile Management**: View and edit user profiles
- **Password Security**: Bcrypt password hashing

### Course Management
- **Browse Courses**: View all available courses
- **Course Enrollment**: Students can enroll in courses
- **Course Creation**: Teachers and admins can create new courses
- **Course Modification**: Edit course details including title, description, and assignments
- **Instructor Assignment**: Admins can assign teachers to courses

### Assignment System
- **Create Assignments**: Teachers can create assignments with questions, instructions, and deadlines
- **Assignment Scheduling**: Set publish time and deadline for assignments
- **View Assignments**: Students and teachers can view all assignments for a course
- **Update/Delete**:  Teachers can modify or remove assignments
- **File Submissions**: Support for PDF, DOC, and DOCX file formats

### Submission & Grading
- **Submit Assignments**: Students can submit their work before deadlines
- **View Submissions**: Teachers can see all submissions for an assignment
- **Grading System**: Teachers can provide marks and feedback
- **Submission Tracking**: Students can view their submission status and grades

### Admin Functions
- **User Management**: View lists of students and teachers
- **Course Assignment**: Assign or remove users from courses
- **System Administration**: Full access to modify courses and manage users

## 🛠️ Tech Stack

- **Language**: Go 1.23.4
- **Web Framework**: Gin (v1.10.0)
- **Database**: PostgreSQL
- **ORM**: GORM (v1.25.12)
- **Authentication**: JWT (golang-jwt/jwt v5.2.2)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **CORS**: gin-contrib/cors
- **Environment Variables**: godotenv

## 📋 Prerequisites

- Go 1.23.4 or higher
- PostgreSQL database
- Git

## 🔧 Installation

1. **Clone the repository**
```bash
git clone https://github.com/Riasat38/Academix.git
cd Academix/api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**

Create a `.env` file in the `api` directory with the following variables:

```env
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=academix
DB_PORT=5432
JWT_SECRET=your_secret_key
```

4. **Run database migrations**

The application automatically migrates the database schema on startup.

5. **Start the server**
```bash
go run server.go
```

The server will start on `http://localhost:8080`

## 📚 API Endpoints

### Authentication
- `POST /academix/signup` - Register a new user
- `POST /academix/login` - Login and receive JWT token
- `POST /academix/logout` - Logout user

### User Profile
- `GET /academix/profile` - Get user profile (Protected)
- `PUT /academix/profile` - Update user profile (Protected)

### Courses
- `GET /academix/course` - Browse all courses (Protected)
- `GET /academix/own-course` - View enrolled/teaching courses (Protected)
- `GET /academix/course/: courseCode` - View specific course details (Protected)
- `POST /academix/enroll-course/: courseCode` - Enroll in a course (Protected)
- `POST /academix/create-course` - Create a new course (Protected)
- `PUT /academix/course/:courseCode` - Edit course details (Protected)

### Assignments
- `POST /academix/: courseCode/assignment` - Create assignment (Protected)
- `GET /academix/:courseCode/assignments` - Get all assignments for a course (Protected)
- `GET /academix/:courseCode/assignments/:assignment_id` - Get specific assignment (Protected)
- `PUT /academix/:courseCode/assignments/:assignment_id` - Update assignment (Protected)
- `DELETE /academix/:courseCode/assignments/: assignment_id` - Delete assignment (Protected)

### Submissions
- `POST /academix/:courseCode/assignment/:assignment_id` - Submit assignment (Protected)
- `GET /academix/:courseCode/assignment/:assignment_id/submissions` - View all submissions (Protected)
- `GET /academix/submission/:submission_id` - View specific submission (Protected)
- `PUT /academix/submissions/:submission_id` - Update marks and feedback (Protected)

### Admin
- `GET /academix/admin/student-list` - Get all students (Admin only)
- `GET /academix/admin/teacher-list` - Get all teachers (Admin only)
- `POST /academix/admin/assign-user/: courseCode` - Assign user to course (Admin only)
- `DELETE /academix/admin/remove-user/:courseCode` - Remove user from course (Admin only)

## 🔐 Authentication

All protected routes require a JWT token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

Tokens are valid for 7 days after issuance.

## 👥 User Roles & Permissions

### Student
- View and enroll in courses
- View assignments
- Submit assignments
- View their own submissions and grades
- Edit their profile

### Teacher
- View all courses
- Create, edit, and delete assignments
- View all submissions for their courses
- Grade submissions (provide marks and feedback)
- Edit their profile

### Admin
- Full access to all features
- Manage users (view student/teacher lists)
- Assign/remove users to/from courses
- Create and modify courses
- Delete users and courses

## 📁 Project Structure

```
Academix/
├── api/
│   ├── auth/              # JWT authentication logic
│   ├── config/            # Database configuration
│   ├── controllers/       # API request handlers
│   │   ├── adminHandler.go
│   │   ├── assignmentHandler.go
│   │   ├── courseHandler.go
│   │   ├── submissionHandler.go
│   │   └── userHandler.go
│   ├── middleware/        # Authentication middleware
│   ├── models/            # Database models (User, Course, Assignment, Submission)
│   ├── permissions/       # Role-based permission validation
│   ├── go.mod            # Go module dependencies
│   └── server.go         # Main application entry point
```

## 🗄️ Database Schema

### UserModel
- ID, Name, Username, Email, Password (hashed), Role
- Relationships:  Courses (enrolled), TaughtCourses, Submissions

### CourseModel
- Code (Primary Key), Title, Description
- Relationships: Students, Instructors, Assignments

### Assignment
- ID, Serial, CourseCode, Question, Instructions, PublishTime, Deadline
- Relationships: Course, Submissions

### AssignmentSubmission
- ID, AssignmentID, StudentID, Submission (file path), Marks, Feedback
- Relationships: Assignment, Student

## 🔒 Security Features

- Password hashing with bcrypt
- JWT-based authentication
- Role-based access control (RBAC)
- CORS configuration
- SQL injection protection via GORM
- Secure cookie handling

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is part of CSE471 Lab Project (Group 3).

## 📧 Contact

Riasat38 - [@Riasat38](https://github.com/Riasat38)

Project Link: [https://github.com/Riasat38/Academix](https://github.com/Riasat38/Academix)

## 🙏 Acknowledgments

- CSE471 Lab Project - Group 3
- Built with [Gin Web Framework](https://gin-gonic.com/)
- Database management with [GORM](https://gorm.io/)
- Authentication with [JWT](https://jwt.io/)
```
