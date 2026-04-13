package main

import (
	"fmt"
	"slices"
)

// Employee represents a worker in the system.
type Employee struct {
	ID     int
	Name   string
	Age    int
	Salary float64
}

// Manager manages a list of employees.
type Manager struct {
	Employees []Employee
}

// AddEmployee adds a new employee.
func (m *Manager) AddEmployee(e Employee) {
	m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID safely.
func (m *Manager) RemoveEmployee(id int) bool {
	for i, e := range m.Employees {
		if e.ID == id {
			m.Employees = slices.Delete(m.Employees, i, i+1)
			return true
		}
	}
	return false
}

// GetAverageSalary calculates average salary safely.
func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) == 0 {
		return 0
	}

	var sum float64
	for _, e := range m.Employees {
		sum += e.Salary
	}

	return sum / float64(len(m.Employees))
}

// FindEmployeeByID finds employee safely and returns pointer to real slice item.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	for i := range m.Employees {
		if m.Employees[i].ID == id {
			return &m.Employees[i]
		}
	}
	return nil
}

// ListEmployees prints all employees (helper function).
func (m *Manager) ListEmployees() {
	if len(m.Employees) == 0 {
		fmt.Println("No employees found.")
		return
	}

	for _, e := range m.Employees {
		fmt.Printf("ID: %d | Name: %s | Age: %d | Salary: %.2f\n",
			e.ID, e.Name, e.Age, e.Salary)
	}
}

func main() {
	manager := &Manager{}

	// Add employees
	manager.AddEmployee(Employee{ID: 1, Name: "Alice", Age: 30, Salary: 70000})
	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})

	// Remove employee safely
	removed := manager.RemoveEmployee(1)
	fmt.Println("Employee removed:", removed)

	// Average salary
	avg := manager.GetAverageSalary()
	fmt.Printf("Average Salary: %.2f\n", avg)

	// Find employee
	emp := manager.FindEmployeeByID(2)
	if emp != nil {
		fmt.Printf("Found employee: %+v\n", *emp)
	} else {
		fmt.Println("Employee not found.")
	}

	// List all
	manager.ListEmployees()
}
