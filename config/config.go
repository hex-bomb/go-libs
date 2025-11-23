package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hex-bomb/go-libs/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Имя и версия приложения (заполняются при билде)
var (
	appName    = "unknown"
	appVersion = "development"
)

// Получить имя приложения
func AppName() string { return appName }

// Получить версию приложения
func AppVersion() string { return appVersion }

// Хуки для viper, применяются при анмаршале
var hooks []viper.DecoderConfigOption

// Регистрирует хук для декодирования конфигов
func RegisterHook(hook mapstructure.DecodeHookFuncType) {
	hooks = append(hooks, viper.DecodeHook(hook))
}

// Загружает конфиг приложения в структуру и проводит валидацию
func GetSettings[AppConf any](defaults map[string]any) (*AppConf, error) {
	var settings AppConf
	var err error

	// Инициализация viper: пути, дефолты, env
	if err = configViper(defaults); err != nil {
		return nil, err
	}

	// Чтение файла или env, анмаршал в структуру, валидация
	if err = readConfigs(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// Настраивает viper: дефолты, пути конфигов, env-переменные
func configViper(defaults map[string]any, paths ...string) error {
	// Устанавливаем значения по умолчанию для каждого ключа
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	// Заменяем точки на нижние подчеркивания для поиска в env-переменных
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Добавляем пути для поиска конфиг-файла
	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigName("config")

	// Если указан путь к конфиг-файлу — проверяем его существование
	if ConfigFileFlag != "" {
		if _, err := os.Stat(ConfigFileFlag); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("указанный конфиг-файл '%s' не найден", ConfigFileFlag)
		}
		viper.SetConfigFile(ConfigFileFlag)
	}

	// Автоматически подхватываем значения из env-переменных
	viper.AutomaticEnv()
	return nil
}

// Читает конфиг, маппит в структуру, валидирует
func readConfigs(settings any) error {
	// Пробуем прочитать конфиг-файл; если не найден — ошибку не возвращаем
	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return err
	}

	// Копируем значения конфига/env в целевую структуру
	if err := viper.Unmarshal(settings, hooks...); err != nil {
		return fmt.Errorf("ошибка анмаршала конфига: %w", err)
	}

	// Валидируем структуру через go-playground/validator
	if err := validator.Get().Struct(settings); err != nil {
		return fmt.Errorf("ошибка валидации конфига: %w", err)
	}

	return nil
}
