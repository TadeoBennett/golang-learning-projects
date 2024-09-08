package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"grademgmt.com/final/initializers"
	"grademgmt.com/final/pkg/models/postgresql"
)

type AppConfig struct {
	students *postgresql.StudentModel
	grades   *postgresql.GradeModel
	teachers *postgresql.TeacherModel
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
}

func main() {
	initializers.LoadEnv() // load environment variables

	//provide flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	defaultDSN := os.Getenv("DSN")
	dsn := flag.String("dsn", defaultDSN, "PostgreSQL DSN (Data Source Name)")
	secret := flag.String("secret", "8693b89c15217db6a4a90aa41cf0e6d5f31752aaf318b4e184f7c5a93a9a90c2", "Secret Key")
	flag.Parse()

	db, err := initializers.ConnectToDatabase(*dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	//setup logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//manage session
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour //expire after 12 hours
	session.Secure = true             //encrypted session keys

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &AppConfig{
		students: &postgresql.StudentModel{
			DB: db,
		},
		grades: &postgresql.GradeModel{
			DB: db,
		},
		teachers: &postgresql.TeacherModel{
			DB: db,
		},
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
	}

	//create our own http server
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second, //max time taken to read a request
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem") //use the certifate values
	srv.ErrorLog.Fatal(err)
}
