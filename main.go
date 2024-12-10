package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

func main() {
	// Define the DSN
	dsn := "web66_65011212075:65011212075@csmsu@tcp(202.28.34.197:3306)/web66_65011212075"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	app := fiber.New()
	app.Post("/user/Loging", LoginUser)
	app.Post("/Log", PostLog)

	log.Fatal(app.Listen(":3000"))
}

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func LoginUser(c *fiber.Ctx) error {
	p := new(User)
	if err := c.BodyParser(p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	// Query the database for the user by email
	row := db.QueryRow(`SELECT Username,Password FROM MACUser WHERE Username = ?`, p.Username)
	// Create a User instance to hold the queried data
	P := new(User)
	err := row.Scan(&P.Username, &P.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return an error if no user is found
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		// Handle other errors (e.g., database issues)
		return fiber.NewError(fiber.StatusInternalServerError, "Database query error")
	}
	// Verify the provided password against the hashed password
	if P.Password == p.Password {
		U := new(User)
		U.Username = P.Username
		U.Password = P.Password

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Login successful",
			"user":    U,
		})
	} else {
		return c.JSON("Invalid email or password")
	}
}

type Log struct {
	Mac       string `json:"Mac"`
	Hostname  string `json:"Hostname"`
	IP        string `json:"IP"`
	Time      string `json:"Time"`
	formatime string `json:"formatime"`
}

func PostLog(c *fiber.Ctx) error {
	// Parse the request body into the `Log` struct
	p := new(Log)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	// Ensure the database connection is valid
	if db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection not initialized",
		})
	}

	// Prepare the SQL query
	query := `INSERT INTO Log (Mac, Hostname, IP, Time,formatime) VALUES (?, ?, ?, ?,?)`

	// Execute the query
	_, err := db.Exec(query, p.Mac, p.Hostname, p.IP, p.Time, p.formatime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to insert log into database",
			"error":   err.Error(),
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Log inserted successfully",
	})
}
