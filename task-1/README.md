# Student Grade Calculator

A simple Go application for calculating and managing student grades with timestamp tracking.

## Features

- Add students with multiple subjects and grades
- Calculate average grades automatically
- Grade classification (A, B, C, D, F)
- Timestamp tracking for creation and updates
- Input validation and error handling
- Interactive command-line interface

## Usage

```bash
# Run the application
go run student_grade_calculator.go

# Run tests
go test -v

# Build executable
go build -o student_calculator.exe
```

## Functions

- `NewStudent()` - Creates a new student with timestamp tracking
- `AddSubject()` - Adds subjects with grade validation and timestamp updates
- `CalculateAverage()` - Computes average grade across all subjects
- `DisplayResults()` - Shows detailed grade report with timestamps
- Time utility functions for better timestamp handling

## Grade Classification

- 90-100: Excellent (A)
- 80-89: Good (B)
- 70-79: Satisfactory (C)
- 60-69: Needs Improvement (D)
- Below 60: Failing (F)