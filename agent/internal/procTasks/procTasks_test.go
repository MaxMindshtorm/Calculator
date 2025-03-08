package procTasks

import (
	"github.com/MaxMindshtorm/calculator/agent/internal/models"
	"testing"
	"time"
)

func TestProcessTask(t *testing.T) {
	tests := []struct {
		name          string
		task          models.Task
		expected      float64
		expectedError bool
	}{
		{
			name: "Addition",
			task: models.Task{
				Operation:     "+",
				Arg1:          "5.2",
				Arg2:          "3.8",
				OperationTime: 10,
			},
			expected:      9.0,
			expectedError: false,
		},
		{
			name: "Subtraction",
			task: models.Task{
				Operation:     "-",
				Arg1:          "10",
				Arg2:          "4.5",
				OperationTime: 10,
			},
			expected:      5.5,
			expectedError: false,
		},
		{
			name: "Multiplication",
			task: models.Task{
				Operation:     "*",
				Arg1:          "6",
				Arg2:          "7",
				OperationTime: 10,
			},
			expected:      42.0,
			expectedError: false,
		},
		{
			name: "Division",
			task: models.Task{
				Operation:     "/",
				Arg1:          "20",
				Arg2:          "4",
				OperationTime: 10,
			},
			expected:      5.0,
			expectedError: false,
		},
		{
			name: "Division by Zero",
			task: models.Task{
				Operation:     "/",
				Arg1:          "20",
				Arg2:          "0",
				OperationTime: 10,
			},
			expected:      0.0,
			expectedError: true,
		},
		{
			name: "Invalid Operation",
			task: models.Task{
				Operation:     "%",
				Arg1:          "20",
				Arg2:          "5",
				OperationTime: 10,
			},
			expected:      0.0,
			expectedError: true,
		},
		{
			name: "Invalid Arg1",
			task: models.Task{
				Operation:     "+",
				Arg1:          "abc",
				Arg2:          "5",
				OperationTime: 10,
			},
			expected:      0.0,
			expectedError: true,
		},
		{
			name: "Invalid Arg2",
			task: models.Task{
				Operation:     "+",
				Arg1:          "10",
				Arg2:          "xyz",
				OperationTime: 10,
			},
			expected:      0.0,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()
			result, err := ProcessTask(&tc.task)
			duration := time.Since(startTime)

			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectedError && result != tc.expected {
				t.Errorf("Expected %f but got %f", tc.expected, result)
			}

			if !tc.expectedError {
				minDuration := time.Duration(tc.task.OperationTime) * time.Millisecond
				if duration < minDuration {
					t.Errorf("Operation completed too quickly. Expected at least %v, got %v", minDuration, duration)
				}
			}
		})
	}
}

func TestProcessTaskFractionalValues(t *testing.T) {
	tests := []struct {
		name     string
		task     models.Task
		expected float64
	}{
		{
			name: "Addition with fractional numbers",
			task: models.Task{
				Operation:     "+",
				Arg1:          "3.14159",
				Arg2:          "2.71828",
				OperationTime: 5,
			},
			expected: 5.85987,
		},
		{
			name: "Subtraction with small numbers",
			task: models.Task{
				Operation:     "-",
				Arg1:          "0.0001",
				Arg2:          "0.0002",
				OperationTime: 5,
			},
			expected: -0.0001,
		},
		{
			name: "Multiplication with large numbers",
			task: models.Task{
				Operation:     "*",
				Arg1:          "99999.9",
				Arg2:          "0.00001",
				OperationTime: 5,
			},
			expected: 0.999999,
		},
		{
			name: "Division with repeating decimal",
			task: models.Task{
				Operation:     "/",
				Arg1:          "1",
				Arg2:          "3",
				OperationTime: 5,
			},
			expected: 0.3333333333333333,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ProcessTask(&tc.task)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			epsilon := 0.0000001
			if diff := result - tc.expected; diff > epsilon || diff < -epsilon {
				t.Errorf("Expected about %f but got %f", tc.expected, result)
			}
		})
	}
}

func BenchmarkProcessTask(b *testing.B) {
	benchmarks := []struct {
		name string
		task models.Task
	}{
		{
			name: "Addition",
			task: models.Task{
				Operation:     "+",
				Arg1:          "5.2",
				Arg2:          "3.8",
				OperationTime: 1,
			},
		},
		{
			name: "Division",
			task: models.Task{
				Operation:     "/",
				Arg1:          "100",
				Arg2:          "2",
				OperationTime: 1,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ProcessTask(&bm.task)
			}
		})
	}
}
