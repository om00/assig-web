package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/om00/assig-web/handler"
	"github.com/om00/assig-web/psqldb"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	psqldb.Dbpath = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening connection to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	fmt.Println("Successfully connected to the database")
}
func main() {

	migrateCmd := flag.String("m", "", "Run migration command (up, down, force N, drop)")
	seederName := flag.String("s", "", "Name of the seeder to run")
	flag.Parse()

	if *migrateCmd != "" {
		psqldb.RunMigrations(*migrateCmd)
		return
	}

	if *seederName != "" {
		psqldb.CallSeederFunction(db, *seederName)
	}

	app := handler.App{Db: psqldb.NewDB(db)}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("app/static-file"))))

	// Route handlers
	http.HandleFunc("/", app.ShowAllUser)
	http.HandleFunc("/login", app.HandleLogin)
	http.HandleFunc("/block-user", app.BlockUser)
	http.HandleFunc("/unblock-user", app.UnblockUser)

	fmt.Println("Server started on :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
