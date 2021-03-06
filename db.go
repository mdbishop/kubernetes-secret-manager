/*
Copyright (c) 2016, UPMC Enterprises
All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name UPMC Enterprises nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL UPMC ENTERPRISES BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
*/

package main

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
)

func getSecretLocal(name string, db *bolt.DB) (*CustomSecretSpec, error) {
	var secret *CustomSecretSpec
	err := db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("Secrets")).Get([]byte(name))
		if data == nil {
			return nil
		}
		decoder := gob.NewDecoder(bytes.NewReader(data))
		err := decoder.Decode(&secret)
		if err != nil {
			return err
		}
		return nil
	})
	return secret, err
}

func persistSecretLocal(name string, customSecret CustomSecretSpec, db *bolt.DB) error {

	data := new(bytes.Buffer)
	enc := gob.NewEncoder(data)
	err := enc.Encode(customSecret)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if err != nil {
			return err
		}
		bucket := tx.Bucket([]byte("Secrets"))
		err = bucket.Put([]byte(name), data.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func deleteSecretLocal(secretName string, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("Secrets")).Delete([]byte(secretName))
	})
	return err
}
