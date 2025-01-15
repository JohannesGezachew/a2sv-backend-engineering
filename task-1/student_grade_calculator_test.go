package main

import (
	"fmt"
	"testing"
)

func TestNewStudent(t *testing.T) {
	student := NewStudent("John Doe")

	if student.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", student.Name)
	}

	if student.Subjects == nil {
		t.Error("Expected subjects map to be initialized")
	}

	if len(student.Subjects) != 0 {
		t.Errorf("Expected empty subjects map, got %d subjects", len(student.Subjects))
	}
}

func TestAddSubject(t *testing.T) {
	student := NewStudent("Jane Smith")

	// Test valid grade
	err := student.AddSubject("Math", 85.5)
	if err != nil {
		t.Errorf("Expected no error for valid grade, got: %v", err)
	}

	if student.Subjects["Math"] != 85.5 {
		t.Errorf("Expected Math grade to be 85.5, got %f", student.Subjects["Math"])
	}

	// Test invalid grade - too low
	err = student.AddSubject("Science", -10)
	if err == nil {
		t.Error("Expected error for negative grade")
	}

	// Test invalid grade - too high
	err = student.AddSubject("History", 150)
	if err == nil {
		t.Error("Expected error for grade over 100")
	}

	// Test boundary values
	err = student.AddSubject("English", 0)
	if err != nil {
		t.Errorf("Expected no error for grade 0, got: %v", err)
	}

	err = student.AddSubject("Art", 100)
	if err != nil {
		t.Errorf("Expected no error for grade 100, got: %v", err)
	}
}

func TestCalculateAverage(t *testing.T) {
	student := NewStudent("Test Student")

	// Test with no subjects
	average := student.CalculateAverage()
	if average != 0 {
		t.Errorf("Expected average 0 for no subjects, got %f", average)
	}

	// Test with one subject
	student.AddSubject("Math", 80)
	average = student.CalculateAverage()
	if average != 80 {
		t.Errorf("Expected average 80 for single subject, got %f", average)
	}

	// Test with multiple subjects
	student.AddSubject("Science", 90)
	student.AddSubject("English", 70)
	average = student.CalculateAverage()
	expected := (80.0 + 90.0 + 70.0) / 3.0
	if average != expected {
		t.Errorf("Expected average %f, got %f", expected, average)
	}
}

func TestCalculateAverageWithDecimals(t *testing.T) {
	student := NewStudent("Decimal Test")

	student.AddSubject("Math", 85.5)
	student.AddSubject("Science", 92.3)
	student.AddSubject("English", 78.7)

	average := student.CalculateAverage()
	expected := (85.5 + 92.3 + 78.7) / 3.0

	if average != expected {
		t.Errorf("Expected average %f, got %f", expected, average)
	}
}

func TestMultipleSubjectsWithSameName(t *testing.T) {
	student := NewStudent("Override Test")

	// Add subject
	student.AddSubject("Math", 80)

	// Override with new grade
	student.AddSubject("Math", 90)

	if student.Subjects["Math"] != 90 {
		t.Errorf("Expected Math grade to be overridden to 90, got %f", student.Subjects["Math"])
	}

	if len(student.Subjects) != 1 {
		t.Errorf("Expected only 1 subject after override, got %d", len(student.Subjects))
	}
}

// Benchmark tests
func BenchmarkAddSubject(b *testing.B) {
	student := NewStudent("Benchmark Student")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		student.AddSubject("Subject", 85.0)
	}
}

func BenchmarkCalculateAverage(b *testing.B) {
	student := NewStudent("Benchmark Student")

	// Add some subjects
	for i := 0; i < 100; i++ {
		student.AddSubject(fmt.Sprintf("Subject%d", i), float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		student.CalculateAverage()
	}
}

// Example test demonstrating usage
func ExampleStudent_CalculateAverage() {
	student := NewStudent("Example Student")
	student.AddSubject("Math", 85)
	student.AddSubject("Science", 90)
	student.AddSubject("English", 80)

	average := student.CalculateAverage()
	fmt.Printf("Average: %.1f", average)
	// Output: Average: 85.0
}