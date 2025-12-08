package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chrono/internal/domain"
)

type KrankheitsExport struct {
	event *EventService
	user  *UserService
}

func NewKrankheitsExportService(e *EventService, u *UserService) KrankheitsExport {
	return KrankheitsExport{
		event: e,
		user:  u,
	}
}

func (svc *KrankheitsExport) ExportForUser(ctx context.Context, userId int64) (string, error) {
	user, err := svc.user.GetById(ctx, userId)
	if err != nil {
		return "", err
	}

	events, err := svc.event.GetAllByUserId(ctx, userId)
	if err != nil {
		return "", err
	}

	// Mitarbeiter | Krankheitstage Gesamt | Daten...	|
	// Alex		   | 5                     | 12.5.2025  | 18.7.2025

	csv := []string{}
	csv = append(csv, "Mitarbeiter,Krankheitstage Gesamt,Krankheitstage")

	processed, err := svc.processUser(user.Username, events)
	if err != nil {
		return "", err
	}

	csv = append(csv, processed)

	return strings.Join(csv, "\n"), nil
}

func (svc *KrankheitsExport) processUser(userName string, events []domain.Event) (string, error) {
	krank := []domain.Event{}
	for _, e := range events {
		if e.Name != "krank" {
			continue
		}
		krank = append(krank, e)
	}

	line := []string{}
	line = append(line, userName)
	line = append(line, fmt.Sprint(len(krank)))
	for _, k := range krank {
		line = append(line, k.ScheduledAt.Format(time.DateOnly))
	}
	//line = append(line, "")

	return strings.Join(line, ","), nil
}

func (svc *KrankheitsExport) ExportAll(ctx context.Context, year int) (string, error) {
	events, err := svc.event.GetForYear(ctx, year)
	if err != nil {
		return "", err
	}

	// Mitarbeiter | Krankheitstage Gesamt | Daten...	|
	// Alex		   | 5                     | 12.5.2025  | 18.7.2025

	userEvents := map[string][]domain.Event{}
	for _, e := range events {
		userEvents[e.User.Username] = append(userEvents[e.User.Username], e.Event)
	}

	csv := []string{}
	csv = append(csv, "Mitarbeiter,Krankheitstage Gesamt,Krankheitstage")

	for u, e := range userEvents {
		processed, err := svc.processUser(u, e)
		if err != nil {
			continue
		}

		csv = append(csv, processed)
	}

	fmt.Println(csv)

	return strings.Join(csv, "\n"), nil
}
