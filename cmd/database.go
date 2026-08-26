package main

import (
	"context"
	"fmt"

	"net"
	"os"
	"time"

	sqlmysql "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDb() (*gorm.DB, error) {
	// Memuat file .env untuk variabel lingkungan
	err := godotenv.Load()
	if err != nil {
		// log.Fatalf("Error loading .env file")
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	// Membaca variabel lingkungan dari file .env
	sshHost := os.Getenv("SSH_HOST")
	sshPort := os.Getenv("SSH_PORT")
	sshUser := os.Getenv("SSH_USER")
	sshPassword := os.Getenv("SSH_PASSWORD")

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Konfigurasi SSH
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10, // Timeout untuk koneksi SSH
	}

	// Fungsi Dial untuk membuka koneksi SSH setiap kali dibutuhkan
	sqlmysql.RegisterDialContext("mysql+tcp", func(ctx context.Context, addr string) (net.Conn, error) {
		// Membuat koneksi SSH setiap kali fungsi ini dipanggil
		sshConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), sshConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to dial SSH: %v", err)
		}

		// Membuat koneksi MySQL melalui SSH tunnel
		mysqlConn, err := sshConn.Dial("tcp", fmt.Sprintf("%s:%s", dbHost, dbPort))
		if err != nil {
			sshConn.Close() // Tutup koneksi SSH jika gagal ke MySQL
			return nil, fmt.Errorf("failed to connect to MySQL over SSH: %v", err)
		}

		// Pastikan koneksi SSH tertutup setelah koneksi MySQL ditutup
		go func() {
			<-ctx.Done() // Tunggu hingga koneksi selesai
			mysqlConn.Close()
			sshConn.Close()
		}()

		return mysqlConn, nil
	})

	// Menggunakan dsn untuk menghubungkan MySQL melalui SSH
	dsn := fmt.Sprintf(
		"%s:%s@mysql+tcp(127.0.0.1:%s)/%s",
		dbUser,
		dbPassword,
		dbPort,
		dbName,
	)

	// Open using GORM
	db, err := gorm.Open(
		gormmysql.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi koneksi ke MySQL: %v", err)
	}

	// Configure the underlying database/sql connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mendapatkan SQL DB: %w",
			err,
		)
	}

	// Set konfigurasi pooling koneksi
	sqlDB.SetMaxOpenConns(10)                  // Maksimal 10 koneksi terbuka
	sqlDB.SetMaxIdleConns(5)                   // Maksimal 5 koneksi idle
	sqlDB.SetConnMaxLifetime(time.Minute * 10) // Koneksi aktif maksimum 10 menit

	// Cek koneksi dengan ping
	if err = sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("gagal terhubung ke MySQL: %v", err)
	}

	fmt.Println("Koneksi ke MySQL melalui SSH berhasil!")
	return db, nil
}
