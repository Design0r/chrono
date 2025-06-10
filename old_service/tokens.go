package service

import (
	"context"
	"log"
	"time"

	"chrono/db/repo"
)

const TokenMonthLifetime = 15

func createToken(
	r *repo.Queries,
	userId int64,
	startDate time.Time,
	endDate time.Time,
	value float64,
) (repo.VacationToken, error) {
	params := repo.CreateTokenParams{
		UserID:    userId,
		StartDate: startDate,
		EndDate:   endDate,
		Value:     value,
	}
	token, err := r.CreateToken(context.Background(), params)
	if err != nil {
		log.Printf("Failed to create vacation token: %v", err)
		return repo.VacationToken{}, err
	}

	return token, nil
}

func CreateToken(
	r *repo.Queries,
	userId int64,
	year int,
	value float64,
) (repo.VacationToken, error) {
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Now().Location())
	endDate := time.Date(year+1, time.March, 1, 0, 0, 0, 0, time.Now().Location())

	return createToken(r, userId, startDate, endDate, value)
}

func DeleteToken(r *repo.Queries, id int64) error {
	err := r.DeleteToken(context.Background(), id)
	if err != nil {
		log.Printf("Failed to delete vacation token: %v", err)
		return err
	}

	return nil
}

func GetValidUserTokenSum(r *repo.Queries, userId int64,
	startDate time.Time,
) (float64, error) {
	params := repo.GetValidUserTokenSumParams{
		UserID:    userId,
		StartDate: startDate,
		EndDate:   startDate,
	}
	count, err := r.GetValidUserTokenSum(context.Background(), params)
	if err != nil {
		log.Printf("Failed to get valid vacation tokens: %v", err)
		return 0, err
	}
	if count == nil {
		return 0.0, nil
	}

	return *count, nil
}

func InitYearlyTokens(r *repo.Queries, user repo.User, year int) error {
	if user.VacationDays <= 0 {
		return nil
	}

	if year < user.CreatedAt.Year() {
		return nil
	}

	params := repo.GetTokenRefreshParams{UserID: user.ID, Year: int64(year)}
	count, _ := r.GetTokenRefresh(context.Background(), params)
	if count > 0 {
		return nil
	}
	_, err := CreateToken(r, user.ID, year, float64(user.VacationDays))
	if err != nil {
		return err
	}

	return AddTokenRefresh(r, user.ID, year)
}

func AddTokenRefresh(r *repo.Queries, userId int64, year int) error {
	params := repo.CreateTokenRefreshParams{UserID: userId, Year: int64(year)}
	_, err := r.CreateTokenRefresh(context.Background(), params)
	if err != nil {
		log.Printf("Failed creating token refresh: %v", err)
		return err
	}

	return nil
}

func UpdateYearlyTokens(r *repo.Queries, user repo.User, year int, value int) error {
	params := repo.GetTokenRefreshParams{UserID: user.ID, Year: int64(year)}
	count, _ := r.GetTokenRefresh(context.Background(), params)
	if count == 0 {
		err := AddTokenRefresh(r, user.ID, year)
		if err != nil {
			return err
		}
	}
	_, err := CreateToken(r, user.ID, year, float64(value))
	if err != nil {
		return err
	}

	return nil
}

func DebugResetTokens(r *repo.Queries) error {
	return r.DebugResetTokens(context.Background())
}

func DebugResetTokenRefresh(r *repo.Queries) error {
	err := r.ResetTokenRefresh(context.Background())
	if err != nil {
		log.Printf("Failed to reset token refresh table: %v", err)
		return err
	}

	return nil
}

func DebugCreateTokenForAcceptedEvents(r *repo.Queries) error {
	users, err := GetAllUsers(r)
	if err != nil {
		return err
	}

	years, err := GetAPICacheYears(r)
	if err != nil {
		return err
	}

	for _, user := range users {
		for _, year := range years {
			InitYearlyTokens(r, user, int(year))
			count, _ := GetVacationCountForUser(r, user.ID, int(year))
			CreateToken(r, user.ID, int(year), float64(-count))
		}
	}

	return nil
}
