package redisurlparser

import (
	"net/url"
	"strings"
)

type RedisURL struct {
	Username string
	Password string
	Host     string
	Port     string
}

func Parse(redis_url string) (RedisURL, error) {
	u, err := url.Parse(redis_url)
	if err != nil {
		return RedisURL{}, err
	}

	result := strings.Split(u.Host, ":")
	if err != nil {
		return RedisURL{}, err
	}

	username, password := getUsernameAndPassword(u.User)
	host := result[0]
	port := result[1]

	ru := RedisURL{username, password, host, port}
	if err != nil {
		return RedisURL{}, err
	}

	return ru, nil
}

func getUsernameAndPassword(user *url.Userinfo) (string, string) {
	var username string
	var password string

	if user != nil {
		username = user.Username()
		password, _ = user.Password()
	}

	return username, password
}
