package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Student represents a student with their subjects and grades
type Student struct {
	Name     string
	Subjects map[string]float64
}

// NewStudent creates a new student instance
func NewStudent(name string) *Student {
	return &Student{
		Name:     name,
		Subjects: make(map[string]float64),
	}
}

// AddSubject adds a subject and grade to the student
func (s *Student) AddSubject(subject string, grade float64) error {
	if grade < 0 || grade > 100 {
		return fmt.Errorf("grade must be between 0 and 100")
	}
	s.Subjects[subject] = grade
	return nil
}

// CalculateAverage calculates the average grade for all subjects
func (s *Student) CalculateAverage() float64 {
	if len(s.Subjects) == 0 {
		return 0
	}

	total := 0.0
	for _, grade := range s.Subjects {
		total += grade
	}
	return total / float64(len(s.Subjects))
}

// DisplayResults shows the student's information and grades
func (s *Student) DisplayResults() {
	fmt.Printf("\n=== Grade Report for %s ===\n", s.Name)
	fmt.Println("Individual Subject Grades:")

	for subject, grade := range s.Subjects {
		fmt.Printf("  %s: %.2f\n", subject, grade)
	}

	average := s.CalculateAverage()
	fmt.Printf("\nAverage Grade: %.2f\n", average)

	// Grade classification
	var classification string
	switch {
	case average >= 90:
		classification = "Excellent (A)"
	case average >= 80:
		classification = "Good (B)"
	case average >= 70:
		classification = "Satisfactory (C)"
	case average >= 60:
		classification = "Needs Improvement (D)"
	default:
		classification = "Failing (F)"
	}

	fmt.Printf("Grade Classification: %s\n", classification)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get student name
	fmt.Print("Enter student name: ")
	studentName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	studentName = strings.TrimSpace(studentName)

	if studentName == "" {
		fmt.Println("Error: Student name cannot be empty")
		return
	}

	student := NewStudent(studentName)

	// Get number of subjects
	fmt.Print("Enter number of subjects: ")
	numSubjectsStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	numSubjectsStr = strings.TrimSpace(numSubjectsStr)

	numSubjects, err := strconv.Atoi(numSubjectsStr)
	if err != nil || numSubjects <= 0 {
		fmt.Println("Error: Please enter a valid positive number for subjects")
		return
	}

	// Input subjects and grades
	fmt.Printf("\nEnter details for %d subjects:\n", numSubjects)

	for i := 0; i < numSubjects; i++ {
		fmt.Printf("\nSubject %d:\n", i+1)

		// Get subject name
		fmt.Print("  Subject name: ")
		subjectName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("  Error reading input: %v\n", err)
			i-- // Retry this iteration
			continue
		}
		subjectName = strings.TrimSpace(subjectName)

		if subjectName == "" {
			fmt.Println("  Error: Subject name cannot be empty. Please try again.")
			i-- // Retry this iteration
			continue
		}

		// Get grade with validation loop
		var grade float64
		for {
			fmt.Print("  Grade (0-100): ")
			gradeStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("  Error reading input: %v\n", err)
				continue
			}
			gradeStr = strings.TrimSpace(gradeStr)

			parsedGrade, err := strconv.ParseFloat(gradeStr, 64)
			if err != nil {
				fmt.Println("  Error: Please enter a valid number")
				continue
			}

			if parsedGrade < 0 || parsedGrade > 100 {
				fmt.Println("  Error: Grade must be between 0 and 100")
				continue
			}

			grade = parsedGrade
			break
		}

		// Add subject to student
		err = student.AddSubject(subjectName, grade)
		if err != nil {
			fmt.Printf("  Error adding subject: %v\n", err)
			i-- // Retry this iteration
			continue
		}

		fmt.Printf("  ✓ Added %s with grade %.2f\n", subjectName, grade)
	}

	// Display results
	student.DisplayResults()

	// Wait for user to press Enter before closing
	fmt.Print("\nPress Enter to exit...")
	reader.ReadString('\n')
}