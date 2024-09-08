//Cahlil Tillett
//Tadeo Bennett

package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"advancedweb.com/test2/pkg/models/postgresql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectToDatabase(dsn string) (*sql.DB, error) {
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err //return nil and the error
	}

	err = db.Ping() //Test our connection
	if err != nil {
		return nil, err
	}
	fmt.Println("Database Connection established")
	return db, nil //return the connection and nil(no errors)
}

// Dependencies (things/variables)
// DEPENDENCY INJECTION
type application struct {
	student_details *postgresql.StudentModel //references the QuoteModel which has the db connection
	teachers        *postgresql.TeacherModel //references the TeacherModel which has the db connection
	errorLog        *log.Logger
	infoLog         *log.Logger
	session         *sessions.Session
}

func main() {
	//load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//server port connection flags
	addr := flag.String("addr", ":4000", "HTTP network address") //set a flag to use a custom port or port :4000 by default
	//saving the default dsn
	dsnDefault := os.Getenv("DB_CONNECTION")
	//using a custom dsn
	dsn := flag.String("dsn", dsnDefault, "PostgreSQL DSN (Data Source Name)")
	secret := flag.String("secret", "8693b89c15217db6a4a90aa41cf0e6d5f31752aaf318b4e184f7c5a93a9a90c2", "Secret Key")//secret key is not important to me so you can't hack me :')
	flag.Parse()

	//create logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	///connect to database
	db, err := ConnectToDatabase(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour //session will expire after 12 hours
	//encrypted session keys
	session.Secure = true //makes cookies become encrypted

	//configure TLS
	//ECDHE - Elliptic curve Diffie-Hellman
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		session:  session,
		student_details: &postgresql.StudentModel{
			DB: db,
		},
		teachers: &postgresql.TeacherModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,     //connection time on the server
		ReadTimeout:  5 * time.Second, //how long should the server take when reading a request, helps stop DOS, DDOS attacks
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem") //use the certifate values
	// err = srv.ListenAndServe() //use the certifate values
	srv.ErrorLog.Fatal(err)
}
