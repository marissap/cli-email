package main

import (
	"encoding/json"
	"io/ioutil"
)

type Item struct {
	pwd   []byte
	email string
	smtp  string
}

type Cache struct {
	items map[string]Item
}

func (c *Cache) Get(pin string) (string, error) {
	encrypedPw := c.items[pin].pwd
	decrypedPw, err := Decrypt(pin, encrypedPw)
	if err != nil {
		return "", err
	}
	return string(decrypedPw), nil
}

func (c *Cache) Add(email, pwd, pin string) error {
	encryptedPw, err := Encrypt(pwd, pin)
	if err != nil {
		return err
	}

	smtp := GetSmtp(email)

	item := Item{
		pwd:   encryptedPw,
		email: email,
		smtp:  smtp,
	}

	c.items[pin] = item

	return nil
}

func (c *Cache) LoadCacheFromFile() (*Cache, error) {
	file, err := ioutil.ReadFile("emailCache.json")
	if err != nil {
		return &Cache{}, err
	}

	data := &Cache{}

	_ = json.Unmarshal([]byte(file), &data)

	return data, nil
}

func (c *Cache) SaveCacheToFile() {
	file, _ := json.MarshalIndent(c.items, "", " ")
	_ = ioutil.WriteFile("emailCache.json", file, 0644)
}
