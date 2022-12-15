package database

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Server struct {
	DB *sqlx.DB
}

func InitializeDBConnection() Server {
	hostname := os.Getenv("DB_HOSTNAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")
	db, err := sqlx.Open("mysql", username+":"+password+"@("+hostname+":"+port+")/"+dbname+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	log.Println("Successfully connected to the database")

	return Server{
		DB: db,
	}
}

func CreateTables(server Server) error {
	_, err := server.DB.Exec("CREATE TABLE IF NOT EXISTS `users` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,`name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,`email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,`email_verified_at` timestamp NULL DEFAULT NULL,`password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,`scope` enum('user', 'moderator', 'admin') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'user',`remember_token` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,`created_at` timestamp NOT NULL DEFAULT current_timestamp(),`updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),PRIMARY KEY (`id`),UNIQUE KEY `users_email_unique` (`email`)) ENGINE = InnoDB AUTO_INCREMENT = 57 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci")
	if err != nil {
		log.Println("Error creating users table")
		return err
	}

	_, err = server.DB.Exec("CREATE TABLE IF NOT EXISTS `auth_tokens` (`id` int(11) NOT NULL AUTO_INCREMENT,`user_id` bigint(20) unsigned NOT NULL,`scope` varchar(255) NOT NULL,`token` text NOT NULL,`createdAt` timestamp NOT NULL DEFAULT current_timestamp(),`updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),`expiresAt` timestamp NOT NULL DEFAULT (current_timestamp() + interval 90 day),`active` tinyint(4) NOT NULL DEFAULT 0,PRIMARY KEY (`id`),KEY `token` (`token`(768)),KEY `user_id` (`user_id`),CONSTRAINT `auth_tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE) ENGINE = InnoDB AUTO_INCREMENT = 51 DEFAULT CHARSET = utf8mb4")
	return err
}
