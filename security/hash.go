/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
