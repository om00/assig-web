package seeders

import (
	"database/sql"
	"fmt"
	"log"
)

func SeedUsers(db *sql.DB) {
	_, err := db.Exec(`
        INSERT INTO users (name, age, phone, email, status, blockReason, created_at, updated_at)
VALUES
    ('John Doe', 30, '123-456-7890', 'john.doe@example.com', 1, NULL, now(), now()),
    ('Jane Smith', 25, '234-567-8901', 'jane.smith@example.com', 1, NULL, now(), now()),
    ('Emily Johnson', 22, '345-678-9012', 'emily.johnson@example.com', 1, NULL, now(), now()),
    ('Michael Brown', 35, '456-789-0123', 'michael.brown@example.com', 1, NULL, now(), now()),
    ('Linda Williams', 28, '567-890-1234', 'linda.williams@example.com', 1, NULL, now(), now()),
    ('David Miller', 40, '678-901-2345', 'david.miller@example.com', 1, NULL, now(), now());
    `)
	if err != nil {
		log.Fatalf("Error seeding products table: %v", err)
	}
	fmt.Println("data seeded successfully")
}
