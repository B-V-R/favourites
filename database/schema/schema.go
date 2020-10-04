/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package schema

// Schema for models
type Schema struct {
	User      User
	Favourite Favourite
}

// User user data
type User struct {
	User     string `gorm:"primaryKey", json:"user"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Favourite user's favourites data
type Favourite struct {
	User      string `json:"user"`
	Favourite string `json:"favourite"`
}
