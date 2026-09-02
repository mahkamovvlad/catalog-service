package section

import "time"

type (
	Repository struct {
		Postgres       RepositoryPostgres
		MigrationTable string `split_words:"true" default:"schema_migrations"`
	}

	RepositoryPostgres struct {
		Address        string        `required:"true"`
		Username       string        `required:"true"`
		Password       string        `required:"true"`
		Name           string        `required:"true"`
		ConnTimeout    time.Duration `split_words:"true" default:"10s"` // Добавили таймаут подключения!
		ReadTimeout    time.Duration `split_words:"true" default:"30s"`
		WriteTimeout   time.Duration `split_words:"true" default:"30s"`
		MigrationTable string        `split_words:"true" default:"schema_migrations"` // Добавили таблицу миграций для TASK-1003!
	}
)
