package psqldb

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/om00/assig-web/models"
	seeders "github.com/om00/assig-web/psqldb/seeders"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

var Dbpath string

type DbIns struct {
	mainDB *sql.DB
}

func RunMigrations(cmd string) {
	migrationsDir := "file://psqldb/migrations"

	m, err := migrate.New(
		migrationsDir,
		Dbpath,
	)
	if err != nil {
		log.Fatalf("Could not initialize migration: %v", err)
	}

	defer m.Close()

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migration UP done successfully.")

	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Migration DOWN done successfully.")

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
		fmt.Println("Database dropped successfully.")

	default:
		fmt.Println("Unknown migration command:", cmd)
	}
}

func CallSeederFunction(db *sql.DB, funcName string) {
	seedersMap := getSeederFunctions()

	if seedFunc, exists := seedersMap[funcName]; exists {

		reflect.ValueOf(seedFunc).Call([]reflect.Value{reflect.ValueOf(db)})

	} else {
		log.Printf("Seeder function %s not found.", funcName)
	}
}

func getSeederFunctions() map[string]interface{} {
	seedersMap := make(map[string]interface{})

	// Add the seeder functions from the seeders package to the map
	seedersMap["seedUsers"] = seeders.SeedUsers

	return seedersMap
}

func NewDB(db *sql.DB) *DbIns {
	return &DbIns{mainDB: db}
}

func (db *DbIns) GetAllUsers(req models.UserRequest) ([]models.User, error) {

	var rows *sql.Rows
	var err error
	if req.Name == "" && req.Email == "" && len(req.Phone) == 0 && req.StatusInt == nil && req.BlockReasonCodeInt == nil {
		rows, err = db.mainDB.Query("SELECT * FROM select_users()") // Call the select_user() function

	} else {

		name, email, phones, _ := prepareStringFields(req)
		rows, err = db.mainDB.Query("SELECT * FROM select_users($1, $2, $3, $4, $5)", name, phones, email, req.BlockReasonCodeInt, req.StatusInt)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store the data
	var data []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Age, &u.Phone, &u.Email, &u.Status,
			&u.BlockReason, &u.BlockReasonCode, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, u)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {

		return nil, err
	}

	return data, nil
}

func (db *DbIns) BlockUser(updateRqeuest models.UserRequest) error {

	if updateRqeuest.UserId > 0 {
		var reason interface{}
		if updateRqeuest.Reason != "" {
			reason = updateRqeuest.Reason
		}
		_, err := db.mainDB.Exec("CALL block_user($1, $2, $3, $4, $5, $6, $7, $8)", updateRqeuest.UserId, nil, nil, nil, updateRqeuest.ReasonCodeInt, reason, nil, nil)
		if err != nil {
			return err
		}
	} else {

		name, email, phones, reason := prepareStringFields(updateRqeuest)

		tx, err := db.mainDB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		_, err = tx.Exec("CALL block_user($1, $2, $3, $4, $5, $6, $7, $8)", nil, name, phones, email, updateRqeuest.ReasonCodeInt, reason, updateRqeuest.StatusInt, updateRqeuest.BlockReasonCodeInt)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Commit the transaction after the procedure call
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}
	}

	return nil
}

func (db *DbIns) UnblockUser(request models.UserRequest) error {

	if request.UserId > 0 {
		_, err := db.mainDB.Exec("CALL unblock_user($1, $2, $3, $4, $5, $6)", request.UserId, nil, nil, nil, nil, nil)
		if err != nil {
			return err
		}
	} else {
		name, email, phones, _ := prepareStringFields(request)

		tx, err := db.mainDB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}
		_, err = tx.Exec("CALL unblock_user($1, $2, $3, $4, $5, $6)", nil, name, phones, email, request.StatusInt, request.BlockReasonCodeInt)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Commit the transaction after the procedure call
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}
	}

	return nil
}

func prepareStringFields(req models.UserRequest) (interface{}, interface{}, interface{}, interface{}) {
	var name, email, reason interface{}
	var phones interface{}

	if req.Name != "" {
		name = req.Name
	}
	if req.Email != "" {
		email = req.Email
	}
	if len(req.Phone) > 0 {
		phones = pq.Array(req.Phone)
	}
	if req.Reason != "" {
		reason = req.Reason
	}

	return name, email, phones, reason
}
