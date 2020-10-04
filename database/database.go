/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"vighnesh.org/favourate/database/schema"
)

// Database database with Database operations
type Database interface {
	CreateTables() error
	Save(data interface{})
	User(user string) schema.User
	Favourite(favourite *schema.Favourite)
	Favourites(user string) []schema.Favourite
	Delete(favourite *schema.Favourite)
}

type database struct {
	*gorm.DB
	Schema schema.Schema
	log    *log.Logger
}

// Favourite add favourite to user favourites
func (database *database) Favourite(favourite *schema.Favourite) {
	database.DB.Create(favourite)
}

// Favourites retrieve user favourites
func (database *database) Favourites(user string) []schema.Favourite {
	u := schema.Favourite{User: user}
	favourites := make([]schema.Favourite, 0)
	database.DB.Model(&u).Find(&favourites)
	return favourites
}

// User retrieves user
func (database *database) User(user string) schema.User {
	u := schema.User{User: user}
	database.DB.First(&u)
	return u
}

// Save saves data in database
func (database *database) Save(data interface{}) {
	database.DB.Create(data)
	database.DB.Save(data)
}

// Delete deletes favourites from database
func (database *database) Delete(favourite *schema.Favourite) {
	database.DB.Exec("DELETE FROM \"favourites\" WHERE \"user\"='" + favourite.User + "'").Commit()
}

// CreateTables creates tables in database
func (database *database) CreateTables() error {
	err := database.DB.AutoMigrate(
		database.Schema.User,
		database.Schema.Favourite,
	)
	database.log.Println("creating tables:", err)
	return err
}

// New creates a new Database with connection and logger
func New() Database {
	database := &database{}
	database.log = log.New(os.Stdout, "favourite-app-database:\t", log.Ldate|log.Ltime|log.Lshortfile)
	conn, err := connect()
	database.log.Println("creating database connection:", err)
	database.DB = conn
	database.Schema = schema.Schema{
		User:      schema.User{},
		Favourite: schema.Favourite{},
	}
	return database
}

func connect() (*gorm.DB, error) {
	dsn := "user=postgres password=password dbname=favourites_db host=database port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
