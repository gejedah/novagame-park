package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Car represents a car in the parking lot
type Car struct {
	registrationNumber string
}

// ParkingSlot represents a single parking slot
type ParkingSlot struct {
	slotNumber int
	isOccupied bool
	car        *Car
	parkedTime int // In hours, for simplicity. In a real app, use time.Time
}

// ParkingLot represents the entire parking lot
type ParkingLot struct {
	capacity      int
	slots         []*ParkingSlot
	occupiedSlots int
}

// NewParkingLot creates a new parking lot with a given capacity
func NewParkingLot(capacity int) *ParkingLot {
	slots := make([]*ParkingSlot, capacity)
	for i := 0; i < capacity; i++ {
		slots[i] = &ParkingSlot{
			slotNumber: i + 1,
			isOccupied: false,
			car:        nil,
			parkedTime: 0,
		}
	}
	return &ParkingLot{
		capacity: capacity,
		slots:    slots,
	}
}

// create_parking_lot {capacity}
func (pl *ParkingLot) CreateParkingLot(capacity int) {
	if pl.capacity > 0 {
		fmt.Println("Parking lot already created.")
		return
	}
	*pl = *NewParkingLot(capacity)
	fmt.Printf("Created a parking lot with %d slots\n", capacity)
}

// park {car_number}
func (pl *ParkingLot) ParkCar(registrationNumber string) {
	if pl.capacity == 0 {
		fmt.Println("Parking lot not created yet.")
		return
	}

	if pl.occupiedSlots >= pl.capacity {
		fmt.Println("Sorry, parking lot is full")
		return
	}

	// Find the nearest available slot (smallest slot number)
	for _, slot := range pl.slots {
		if !slot.isOccupied {
			slot.car = &Car{registrationNumber: registrationNumber}
			slot.parkedTime = 0 // Reset parked time on new park
			slot.isOccupied = true
			pl.occupiedSlots++
			fmt.Printf("Allocated slot number: %d\n", slot.slotNumber)
			return
		}
	}
}

// leave {car_number} {hours}
func (pl *ParkingLot) LeaveCar(registrationNumber string, hours int) {
	if pl.capacity == 0 {
		fmt.Println("Parking lot not created yet.")
		return
	}

	found := false
	for _, slot := range pl.slots {
		if slot.isOccupied && slot.car.registrationNumber == registrationNumber {
			charge := calculateCharge(hours)
			fmt.Printf("Registration number %s with Slot Number %d is free with Charge $%d\n", registrationNumber, slot.slotNumber, charge)
			slot.isOccupied = false
			slot.car = nil
			slot.parkedTime = 0
			pl.occupiedSlots--
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("Registration number %s not found\n", registrationNumber)
	}
}

// calculateCharge calculates the parking charge based on hours
func calculateCharge(hours int) int {
	if hours <= 2 {
		return 10
	}
	return 10 + (hours-2)*10
}

// status
func (pl *ParkingLot) PrintStatus() {
	if pl.capacity == 0 {
		fmt.Println("Parking lot not created yet.")
		return
	}
	if pl.occupiedSlots == 0 {
		fmt.Println("Parking lot is empty.")
		return
	}

	fmt.Println("Slot No. Registration No.")
	for _, slot := range pl.slots {
		if slot.isOccupied {
			fmt.Printf("%d %s\n", slot.slotNumber, slot.car.registrationNumber)
		}
	}
}

// processCommand processes a single command
func processCommand(command string, parkingLot *ParkingLot) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	switch cmd {
	case "create_parking_lot":
		if len(parts) > 1 {
			capacity, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid capacity")
				return
			}
			parkingLot.CreateParkingLot(capacity)
		} else {
			fmt.Println("Usage: create_parking_lot {capacity}")
		}
	case "park":
		if len(parts) > 1 {
			parkingLot.ParkCar(parts[1])
		} else {
			fmt.Println("Usage: park {car_number}")
		}
	case "leave":
		if len(parts) > 2 {
			hours, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid hours")
				return
			}
			parkingLot.LeaveCar(parts[1], hours)
		} else {
			fmt.Println("Usage: leave {car_number} {hours}")
		}
	case "status":
		parkingLot.PrintStatus()
	default:
		fmt.Println("Unknown command:", command)
	}
}

func main() {
	parkingLot := &ParkingLot{} // Initialize an empty parking lot

	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as a parameter.")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := scanner.Text()
		processCommand(command, parkingLot)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
