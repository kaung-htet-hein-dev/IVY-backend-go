package seeder

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seeder interface {
	Seed() error
}

type DBSeeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) Seeder {
	return &DBSeeder{db: db}
}

func (s *DBSeeder) Seed() error {
	// Clear all existing data
	if err := s.clearTables(); err != nil {
		return err
	}

	// Seed users
	users := make([]entity.User, 0, 10)

	// Seed categories
	categories, err := s.seedCategories()
	if err != nil {
		return err
	}

	// Seed branches
	branches, err := s.seedBranches()
	if err != nil {
		return err
	}

	// Seed services with branches
	services, err := s.seedServices(categories, branches)
	if err != nil {
		return err
	}

	// Seed bookings
	if err := s.seedBookings(users, services, branches); err != nil {
		return err
	}

	log.Println("Seeding completed successfully!")
	return nil
}

func (s *DBSeeder) clearTables() error {
	// Delete in reverse order of dependencies
	if err := s.db.Exec("DELETE FROM bookings").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM branch_service").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM services").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM branches").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM categories").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	return nil
}

func (s *DBSeeder) seedCategories() ([]entity.Category, error) {
	// Initialize categories slice with capacity for all categories
	categories := make([]entity.Category, 0, 11)

	// Create 11 categories
	for i := 1; i <= 11; i++ {
		category := entity.Category{
			ID:   uuid.New(),
			Name: fmt.Sprintf("Category %d", i),
		}
		categories = append(categories, category)
	}

	if err := s.db.Create(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *DBSeeder) seedBranches() ([]entity.Branch, error) {
	// Initialize branches slice with capacity for all branches
	branches := make([]entity.Branch, 0, 11)

	// Create 11 branches
	for i := 1; i <= 11; i++ {
		branch := entity.Branch{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("Branch %d", i),
			Location:    fmt.Sprintf("%d Main Street", i*100),
			Longitude:   fmt.Sprintf("%f", 40.7128+float64(i)*0.01),
			Latitude:    fmt.Sprintf("%f", -74.0060+float64(i)*0.01),
			PhoneNumber: fmt.Sprintf("555-%04d", i),
		}
		branches = append(branches, branch)
	}

	if err := s.db.Create(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

func (s *DBSeeder) seedServices(categories []entity.Category, branches []entity.Branch) ([]entity.Service, error) {
	// Initialize services slice with capacity for all services
	services := make([]entity.Service, 0, 20)

	// Create 20 services
	for i := 1; i <= 20; i++ {
		// Assign random price between 25 and 200
		price := 25 + (i*8)%175
		// Assign to random category
		categoryIndex := (i - 1) % len(categories)
		// Assign to 3 random branches
		branchStart := (i - 1) % (len(branches) - 2)
		selectedBranches := branches[branchStart : branchStart+3]

		service := entity.Service{
			ID:             uuid.New(),
			Name:           fmt.Sprintf("Service %d", i),
			Description:    fmt.Sprintf("This is service number %d", i),
			DurationMinute: 30,
			Price:          price,
			CategoryID:     categories[categoryIndex].ID,
			Image:          fmt.Sprintf("service_%d.jpg", i),
			IsActive:       true,
			Branches:       selectedBranches,
		}
		services = append(services, service)
	}

	if err := s.db.Create(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (s *DBSeeder) seedBookings(users []entity.User, services []entity.Service, branches []entity.Branch) error {
	// Create some bookings with different statuses
	statuses := []string{"PENDING", "CONFIRMED", "COMPLETED", "CANCELLED"}

	// Available time slots
	timeSlots := []string{
		"09:00 AM", "10:00 AM", "11:00 AM", "01:00 PM", "02:00 PM",
		"03:00 PM", "04:00 PM", "05:00 PM", "06:00 PM",
	}

	// Create bookings for the next 7 days
	bookings := make([]entity.Booking, 0, 30)
	for i := 0; i < 30; i++ {
		// Distribute bookings across next 7 days
		bookingDate := time.Now().AddDate(0, 0, i%7)

		// Distribute users, services, branches, times and statuses
		booking := entity.Booking{
			ID:         uuid.New(),
			UserID:     users[i%len(users)].ID,
			ServiceID:  services[i%len(services)].ID,
			BranchID:   branches[i%len(branches)].ID,
			BookedDate: bookingDate.Format("2006-01-02"),
			BookedTime: timeSlots[i%len(timeSlots)],
			Status:     statuses[i%len(statuses)],
		}
		bookings = append(bookings, booking)
	}

	return s.db.Create(&bookings).Error
}
