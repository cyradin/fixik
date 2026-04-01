//nolint:mnd
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/user"
	"github.com/cyradin/fixik/pkg/logger"
)

var GitCommit string = "dev"

const (
	teamFrontend = "frontend"
	teamBackend  = "backend"
	teamInfra    = "infra"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	container := container.New(GitCommit, cfg)

	if err := run(cfg, container); err != nil {
		container.Logger().Error("application error", logger.Error(err))
	}
}

func run(_ *config.Config, container *container.Container) error {
	ctx := context.Background()

	container.Logger().Info("truncating tables...")

	_, err := container.PgPool().Exec(ctx, `
		TRUNCATE TABLE users, teams, statuses, priorities
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("truncate tables: %w", err)
	}

	container.Logger().Info("tables truncated...")

	return seed(ctx, container)
}

func seed(ctx context.Context, c *container.Container) error {
	if _, err := createStatuses(ctx, c); err != nil {
		return err
	}

	if _, err := createPriorities(ctx, c); err != nil {
		return err
	}

	teamIDs, err := createTeams(ctx, c)
	if err != nil {
		return err
	}

	if err := createAdmin(ctx, c); err != nil {
		return err
	}

	if err := createUsers(ctx, c, teamIDs); err != nil {
		return err
	}

	if err := createIncidents(ctx, c, teamIDs); err != nil {
		return err
	}

	return nil
}

func createStatuses(ctx context.Context, c *container.Container) (map[string]int64, error) {
	data := []status.Status{
		{
			Name:        "TODO",
			Code:        "todo",
			Description: "Инцидент создан, но работа ещё не начата",
			Sort:        10,
		},
		{
			Name:        "В процессе",
			Code:        "in_progress",
			Description: "Инцидент находится в работе",
			Sort:        20,
		},
		{
			Name:        "Готово",
			Code:        "done",
			Description: "Инцидент успешно решён",
			Sort:        30,
			IsFinal:     true,
		},
		{
			Name:        "Отменено",
			Code:        "cancelled",
			Description: "Работа по инциденту отменена",
			Sort:        40,
			IsFinal:     true,
		},
	}

	ids := make(map[string]int64)

	for _, d := range data {
		s, err := c.StatusManager().Create(ctx, d)
		if err != nil {
			return nil, fmt.Errorf("create status %s: %w", d.Code, err)
		}

		ids[d.Code] = s.ID
	}

	return ids, nil
}

func createPriorities(ctx context.Context, c *container.Container) (map[string]int64, error) {
	data := []dict.Entity{
		{
			Name:        "P1",
			Code:        "p1",
			Description: "Критический инцидент, система не работает. Время реакции: немедленно",
			Sort:        10,
		},
		{
			Name:        "P2",
			Code:        "p2",
			Description: "Высокий приоритет, серьёзная деградация. Время реакции: 30мин",
			Sort:        20,
		},
		{
			Name:        "P3",
			Code:        "p3",
			Description: "Средний приоритет, частичная проблема. Время реакции: 4ч",
			Sort:        30,
		},
		{
			Name:        "P4",
			Code:        "p4",
			Description: "Низкий приоритет, незначительная ошибка. Время реакции: 1 нед",
			Sort:        40,
		},
	}

	ids := make(map[string]int64)

	for _, d := range data {
		p, err := c.PriorityManager().Create(ctx, d)
		if err != nil {
			return nil, fmt.Errorf("create priority %s: %w", d.Code, err)
		}

		ids[d.Code] = p.ID
	}

	return ids, nil
}

func createTeams(ctx context.Context, c *container.Container) (map[string]int64, error) {
	data := []dict.Entity{
		{
			Name:        teamFrontend,
			Code:        teamFrontend,
			Description: "Команда разработки пользовательского интерфейса",
			Sort:        10,
		},
		{
			Name:        teamBackend,
			Code:        teamBackend,
			Description: "Команда разработки серверной логики",
			Sort:        20,
		},
		{
			Name:        teamInfra,
			Code:        teamInfra,
			Description: "Команда инфраструктуры и DevOps",
			Sort:        30,
		},
	}

	ids := make(map[string]int64)

	for _, d := range data {
		t, err := c.TeamManager().Create(ctx, d)
		if err != nil {
			return nil, fmt.Errorf("create team %s: %w", d.Code, err)
		}

		ids[d.Code] = t.ID
	}

	return ids, nil
}

func createAdmin(ctx context.Context, c *container.Container) error {
	_, err := c.UserManager().Create(ctx, user.CreateUser{
		Name:     "Администратор",
		Username: "admin",
		Email:    "admin@example.com",
		Password: "1234",
		Role:     "admin",
	})
	if err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	return nil
}

func createUsers(ctx context.Context, c *container.Container, teamIDs map[string]int64) error {
	users := []user.CreateUser{
		{
			Name:     "Алексей Смирнов",
			Username: "alexey.smirnov",
			Email:    "alexey.smirnov@example.com",
			Password: "1234",
			Role:     "manager",
			TeamID:   new(teamIDs[teamFrontend]),
		},
		{
			Name:     "Мария Иванова",
			Username: "maria.ivanova",
			Email:    "maria.ivanova@example.com",
			Password: "1234",
			Role:     "user",
			TeamID:   new(teamIDs[teamFrontend]),
		},
		{
			Name:     "Илья Кузнецов",
			Username: "ilya.kuznetsov",
			Email:    "ilya.kuznetsov@example.com",
			Password: "1234",
			Role:     "user",
			TeamID:   new(teamIDs[teamFrontend]),
		},
		{
			Name:     "Дмитрий Петров",
			Username: "dmitry.petrov",
			Email:    "dmitry.petrov@example.com",
			Password: "1234",
			Role:     "manager",
			TeamID:   new(teamIDs[teamBackend]),
		},
		{
			Name:     "Анна Соколова",
			Username: "anna.sokolova",
			Email:    "anna.sokolova@example.com",
			Password: "1234",
			Role:     "user",
			TeamID:   new(teamIDs[teamBackend]),
		},
		{
			Name:     "Сергей Волков",
			Username: "sergey.volkov",
			Email:    "sergey.volkov@example.com",
			Password: "1234",
			Role:     "user",
			TeamID:   new(teamIDs[teamBackend]),
		},
		{
			Name:     "Андрей Попов",
			Username: "andrey.popov",
			Email:    "andrey.popov@example.com",
			Password: "1234",
			Role:     "manager",
			TeamID:   new(teamIDs[teamInfra]),
		},
		{
			Name:     "Екатерина Морозова",
			Username: "ekaterina.morozova",
			Email:    "ekaterina.morozova@example.com",
			Password: "1234",
			Role:     "user",
			TeamID:   new(teamIDs[teamInfra]),
		},
		{
			Name:     "Никита Лебедев",
			Username: "nikita.lebedev",
			Email:    "nikita.lebedev@example.com",
			Password: "1234",
			Role:     "user",
			TeamID:   new(teamIDs[teamInfra]),
		},
	}

	for _, u := range users {
		if _, err := c.UserManager().Create(ctx, u); err != nil {
			return fmt.Errorf("create user %s: %w", u.Name, err)
		}
	}

	return nil
}

func createIncidents(ctx context.Context, c *container.Container, teamIDs map[string]int64) error {
	users, err := c.UserManager().List(ctx, 100, 0) //nolint:mnd
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}

	var (
		frontendUsers []int64
		backendUsers  []int64
		infraUsers    []int64
	)

	for _, u := range users {
		if u.TeamID == nil {
			continue
		}

		switch *u.TeamID {
		case teamIDs[teamFrontend]:
			frontendUsers = append(frontendUsers, u.ID)
		case teamIDs[teamBackend]:
			backendUsers = append(backendUsers, u.ID)
		case teamIDs[teamInfra]:
			infraUsers = append(infraUsers, u.ID)
		}
	}

	statuses, err := c.StatusManager().List(ctx)
	if err != nil {
		return fmt.Errorf("list statuses: %w", err)
	}

	priorities, err := c.PriorityManager().List(ctx)
	if err != nil {
		return fmt.Errorf("list priorities: %w", err)
	}

	statusMap := map[string]int64{}
	for _, s := range statuses {
		statusMap[s.Code] = s.ID
	}

	priorityMap := map[string]int64{}
	for _, p := range priorities {
		priorityMap[p.Code] = p.ID
	}

	type incidentSeed struct {
		Title       string
		Description string
		StatusID    int64
		PriorityID  int64
		TeamID      *int64
		UserID      *int64
	}

	data := []incidentSeed{
		{
			Title:       "Не открывается главная страница",
			Description: "HTTP 500 при заходе на сайт",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p1"],
			TeamID:      nil,
			UserID:      nil,
		},
		{
			Title:       "Ошибка авторизации",
			Description: "Пользователь не может войти в систему",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p2"],
			TeamID:      nil,
			UserID:      nil,
		},
		{
			Title:       "Медленная загрузка каталога",
			Description: "Каталог товаров открывается более 10 секунд",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p3"],
			TeamID:      nil,
			UserID:      nil,
		},
		{
			Title:       "Проблема с CSS",
			Description: "Некорректное отображение кнопок",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p4"],
			TeamID:      nil,
			UserID:      nil,
		},
		{
			Title:       "Кнопка оплаты не работает",
			Description: "Клик не отправляет запрос",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamFrontend]),
			UserID:      nil,
		},
		{
			Title:       "Неверная верстка Safari",
			Description: "Сетка ломается на Safari",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p3"],
			TeamID:      new(teamIDs[teamFrontend]),
			UserID:      nil,
		},
		{
			Title:       "Ошибка API заказов",
			Description: "При создании нового заказа API возвращает 500. Нужно проверить обработку данных на backend, отладить сервис заказов, исключить nil pointer и убедиться, что новые заказы создаются корректно, включая все необходимые проверки и валидацию.",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p1"],
			TeamID:      new(teamIDs[teamBackend]),
			UserID:      nil,
		},
		{
			Title:       "Медленный SQL запрос",
			Description: "Запрос выполняется более 15 секунд",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamBackend]),
			UserID:      nil,
		},
		{
			Title:       "Ошибка Redis кэша",
			Description: "Кэш возвращает старые данные",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p3"],
			TeamID:      new(teamIDs[teamBackend]),
			UserID:      nil,
		},
		{
			Title:       "Падает CI pipeline",
			Description: "GitLab runner падает при сборке проекта. Нужно проверить конфигурацию CI, docker-образа, доступность зависимостей и корректность скриптов, чтобы pipeline стабильно выполнялся и не блокировал релизы.",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p1"],
			TeamID:      new(teamIDs[teamInfra]),
			UserID:      nil,
		},
		{
			Title:       "Недостаточно места на диске",
			Description: "Disk usage достиг 95%",
			StatusID:    statusMap["todo"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamInfra]),
			UserID:      nil,
		},

		{
			Title:       "Фронтенд билд падает",
			Description: "Сборка проекта через webpack падает с ошибками компиляции. Необходимо проверить конфигурацию сборки, зависимости npm, версию Node.js и исправить проблемы, чтобы фронтенд собирался стабильно на всех окружениях.",
			StatusID:    statusMap["in_progress"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamFrontend]),
			UserID:      new(frontendUsers[0]),
		},
		{
			Title:       "React exception в production",
			Description: "Unhandled promise rejection",
			StatusID:    statusMap["in_progress"],
			PriorityID:  priorityMap["p1"],
			TeamID:      new(teamIDs[teamFrontend]),
			UserID:      new(frontendUsers[1]),
		},
		{
			Title:       "Форма регистрации ломается",
			Description: "Ошибка валидации email",
			StatusID:    statusMap["done"],
			PriorityID:  priorityMap["p3"],
			TeamID:      new(teamIDs[teamFrontend]),
			UserID:      new(frontendUsers[2]),
		},

		{
			Title:       "Падает сервис заказов",
			Description: "Сервис заказов падает с panic: nil pointer dereference. Нужно найти участок кода, который вызывает ошибку при обработке данных, добавить проверки на nil, протестировать создание заказов и восстановить стабильную работу сервиса.",
			StatusID:    statusMap["in_progress"],
			PriorityID:  priorityMap["p1"],
			TeamID:      new(teamIDs[teamBackend]),
			UserID:      new(backendUsers[0]),
		},
		{
			Title:       "Некорректная сериализация JSON",
			Description: "marshal error",
			StatusID:    statusMap["done"],
			PriorityID:  priorityMap["p3"],
			TeamID:      new(teamIDs[teamBackend]),
			UserID:      new(backendUsers[1]),
		},
		{
			Title:       "Проблема с поиском",
			Description: "Elasticsearch не возвращает результаты для некоторых запросов. Необходимо проверить индексирование данных, корректность mapping и запросов, а также нагрузку на кластер, чтобы поиск работал корректно для всех пользователей.",
			StatusID:    statusMap["in_progress"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamBackend]),
			UserID:      new(backendUsers[2]),
		},

		{
			Title:       "Kubernetes pod crashloop",
			Description: "Контейнер в k8s pod постоянно перезапускается (CrashLoopBackOff). Нужно проверить логи контейнера, ошибки конфигурации, healthchecks, зависимости и ресурсы, чтобы pod стабильно запускался и сервис оставался доступным.",
			StatusID:    statusMap["in_progress"],
			PriorityID:  priorityMap["p1"],
			TeamID:      new(teamIDs[teamInfra]),
			UserID:      new(infraUsers[0]),
		},
		{
			Title:       "Ошибка деплоя Helm",
			Description: "Helm upgrade failed",
			StatusID:    statusMap["done"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamInfra]),
			UserID:      new(infraUsers[1]),
		},
		{
			Title:       "Проблема с балансировщиком",
			Description: "nginx не проксирует трафик",
			StatusID:    statusMap["in_progress"],
			PriorityID:  priorityMap["p2"],
			TeamID:      new(teamIDs[teamInfra]),
			UserID:      new(infraUsers[2]),
		},
	}

	for _, d := range data {
		_, err := c.IncidentManager().Create(ctx, incident.CreateIncident{
			Title:       d.Title,
			Description: d.Description,
			StatusID:    d.StatusID,
			PriorityID:  d.PriorityID,
			TeamID:      d.TeamID,
			UserID:      d.UserID,
		})
		if err != nil {
			return fmt.Errorf("create incident %s: %w", d.Title, err)
		}
	}

	return nil
}
