package service

import (
	"encoding/json"
	"log"
	"os"

	"chrono/config"
	"chrono/db/repo"
	"chrono/schemas"
)

func LoadDebugUsers(r *repo.Queries, cfg *config.Config) (schemas.DebugUsers, error) {
	if !cfg.Debug {
		return schemas.DebugUsers{}, nil
	}

	file, err := os.ReadFile(cfg.DebugUsers)
	if err != nil {
		log.Printf("Failed loading debug user file: %v", err)
		return schemas.DebugUsers{}, err
	}

	users := schemas.DebugUsers{}
	err = json.Unmarshal(file, &users)
	if err != nil {
		log.Printf("Failed unmarshalling debug users: %v", err)
		return schemas.DebugUsers{}, err
	}

	for _, user := range users.Users {
		pw, err := HashPassword(user.Password)
		if err != nil {
			continue
		}
		user.Password = pw
		CreateUser(r, user)
	}

	log.Println("Created debug users")
	return users, nil
}
