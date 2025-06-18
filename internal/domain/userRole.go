package domain

type UserRole string

const (
	UserEmployee UserRole = "employee"
	UserManager  UserRole = "manager"
	UserHR       UserRole = "hr"
	userDirector UserRole = "director"
)
