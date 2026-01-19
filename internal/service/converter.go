package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/konkovaanna23/shortener/internal/config/db"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	lengthURL = 6
)

const (
	ModeStoreDB      = "DB"
	ModeStoreFile    = "FILE"
	ModeStoreStorage = "STORAGE"
)

// Converter - основной сервис, который конвертирует URL.
type Converter struct {
	url              string
	storage          *model.Storage
	db               *sqlx.DB
	modeStore        string
	filePath         string
	fMx              sync.RWMutex
	userURLS         *model.UserURLS
	httpDeleteChan   chan *model.DescriptionURL
	manualDeleteChan chan *model.DescriptionURL
}

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	URLShort string `json:"result"`
}

// NewConverter создаёт новый экземпляр Converter — центрального компонента бизнес-логики
// сервиса сокращения URL.
//
// Converter отвечает за:
//   - Генерацию коротких ключей
//   - Хранение соответствий "короткий URL → исходный URL"
//   - Управление режимом хранения: в памяти, файле или БД
//   - Асинхронную очистку удалённых URL (через каналы)
//   - Сохранение состояния при остановке (если используется файл)
//
// Поведение зависит от переданных параметров:
//   - Если db != nil → данные хранятся в PostgreSQL (режим ModeStoreDB)
//   - Если db == nil и filePath задан → сохранение в JSON-файл (ModeStoreFile)
//   - Если db == nil и filePath пуст → данные в памяти (ModeStoreStorage)
//
// Запускает фоновый процесс обработки удаления URL через StartDeleteProcessor.
//
// Параметры:
//   - ctx: контекст для graceful shutdown и остановки фоновых процессов
//   - serverURL: базовый URL сервера (например, "http://localhost:8080"), используется для формирования полных коротких ссылок
//   - filePath: путь к файлу для сохранения данных (если используется файловое хранилище)
//   - db: подключение к базе данных *sqlx.DB (если используется БД)
//   - bufferSize: ёмкость каналов для асинхронного удаления URL
//   - batchSize: количество URL, обрабатываемых за одну итерацию фонового удаления
//   - timeFlushDel: интервал (в секундах) между проверками очереди на удаление
//
// Возвращает *Converter, готовый к использованию.
//
// Пример:
//
//	db, _ := sqlx.Connect("postgres", "...")
//	ctx := context.Background()
//
//	converter := NewConverter(
//	    ctx,
//	    "http://localhost:8080",
//	    "/data/urls.json",
//	    db,
//	    1000,     // bufferSize
//	    100,      // batchSize
//	    10,       // timeFlushDel (секунд)
//	)
func NewConverter(ctx context.Context, serverURL string, filePath string, db *sqlx.DB, bufferSize, batchSize, timeFlushDel int) *Converter {
	cvrt := &Converter{
		url:              serverURL,
		storage:          model.NewStorage(lengthURL),
		db:               db,
		filePath:         filePath,
		userURLS:         model.NewUserURLS(),
		httpDeleteChan:   make(chan *model.DescriptionURL, bufferSize),
		manualDeleteChan: make(chan *model.DescriptionURL, bufferSize),
	}

	if db != nil {
		cvrt.modeStore = ModeStoreDB
		logrus.Println("Установлен режим сохранения в БД")
	} else {
		if filePath != "" {
			cvrt.modeStore = ModeStoreFile
			logrus.Println("Установлен режим сохранения в файл")
		} else {
			cvrt.modeStore = ModeStoreStorage
			logrus.Println("Установлен режим сохранения в хранилище")
		}

	}
	cvrt.StartDeleteProcessor(ctx, batchSize, time.Duration(timeFlushDel)*time.Second)
	return cvrt
}

// AddURL добавляет новый URL в хранилище и возвращает короткую ссылку.
func (c *Converter) AddURL(url string, user string) (string, error) {
	if !c.isValidURL(url) {
		msg := fmt.Sprintf("URL [%s] не является валидным", url)
		return "", fmt.Errorf("%s", msg)
	}
	shortURL := c.storage.GenerateShortURL()
	var err error
	switch c.modeStore {
	case ModeStoreDB:
		url := []*model.DescriptionURL{{Short: shortURL, Original: url}}
		if err = c.TranStoreURLInDB(url, user); err != nil {
			logrus.Errorln("ошибка сохранения в базу:", err)
		}
		shortURL = url[0].Short
	case ModeStoreFile:
		if shortURL, err = c.StoreURLInFile(shortURL, url, user); err != nil {
			logrus.Errorln("ошибка сохранения в файл:", err)
		}
	case ModeStoreStorage:
		shortURL, err = c.storage.Add(url, shortURL)
		c.userURLS.AddURLForUser(user, shortURL)
	}
	return c.url + "/" + shortURL, err
}

// GetURL возвращает оригинальный URL по короткой ссылке.
func (c *Converter) GetURL(shortURL string) (string, error) {
	var URL string
	var err error
	switch c.modeStore {
	case ModeStoreDB:
		URL, err = c.GetOriginalURLFromDB(shortURL)
		if err != nil {
			logrus.Errorln(err)
			return "", err
		}
	case ModeStoreFile:
		URL, err = c.GetOriginalURLFromFile(shortURL)
		if err != nil {
			logrus.Errorln(err)
			return "", err
		}
	case ModeStoreStorage:
		URL, err = c.storage.Get(shortURL)
		if err != nil {
			return "", err
		}
	}
	return URL, nil
}

// isValidURL проверяет, является ли строка валидным URL.
func (c *Converter) isValidURL(input string) bool {
	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

// AddURLForRequest обрабатывает входящий запрос на добавление URL.
func (c *Converter) AddURLForRequest(url *URLRequest, user string) (*URLResponse, error) {
	if url == nil {
		return nil, errors.New("передана пустая структура")
	}
	result, err := c.AddURL(url.URL, user)
	if err != nil {
		if errors.Is(err, model.ErrorConflictURL) {
			return &URLResponse{URLShort: result}, err
		} else {
			return nil, err
		}
	}
	return &URLResponse{URLShort: result}, nil
}

// PingDB проверяет подключение к базе данных.
func (c *Converter) PingDB() error {
	err := db.Ping(c.db)
	if err != nil {
		return err
	}
	return nil
}

// AddURLForBatch добавляет несколько URL в хранилище и возвращает их.
func (c *Converter) AddURLForBatch(urls []*model.DescriptionURL, user string) ([]*model.DescriptionURL, error) {
	var errs []error
	for _, url := range urls {
		if !c.isValidURL(url.Original) {
			errs = append(errs, fmt.Errorf("URL [%s] не является валидным", url.Original))
		}
	}
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	for _, url := range urls {
		url.Short = c.storage.GenerateShortURL()
	}

	switch c.modeStore {
	case ModeStoreDB:
		err := c.TranStoreURLInDB(urls, user)
		if err != nil {
			logrus.Error("ошибка сохранения в базу:", err)
			return nil, err
		}
	case ModeStoreFile:
		err := c.StoreURLsInFile(urls, user)
		if err != nil {
			logrus.Error("ошибка сохранения в файл:", err)
			return nil, err
		}
	case ModeStoreStorage:
		for _, url := range urls {
			short, _ := c.storage.Add(url.Original, url.Short)
			url.Short = short
			c.userURLS.AddURLForUser(user, url.Short)

		}
	}
	for _, url := range urls {
		url.Original = ""
		url.Short = c.url + "/" + url.Short
	}
	return urls, nil
}

func (c *Converter) convertMapToDescriptionURL(m map[string]string) []*model.DescriptionURL {
	result := make([]*model.DescriptionURL, len(m))
	i := 0
	for key, value := range m {
		result[i] = &model.DescriptionURL{Short: c.url + "/" + key, Original: value}
		i++
	}
	return result
}

// GetURLsForUser возвращает список URL для указанного пользователя.
func (c *Converter) GetURLsForUser(user string) ([]*model.DescriptionURL, error) {
	var urls []*model.DescriptionURL
	var err error
	switch c.modeStore {
	case ModeStoreDB:
		urls, err = c.GetInfoUserURLFromDB(user)
		if err != nil {
			logrus.Error("ошибка чтения из базы:", err)
			return nil, err
		}

	case ModeStoreFile:
		urls, err = c.GetInfoUserURLFromFile(user)
		if err != nil {
			logrus.Error("ошибка чтения из файла:", err)
		}
	case ModeStoreStorage:
		urlList := c.userURLS.GetURLsForUser(user)
		resultMap := c.storage.GetURLMapForList(urlList)
		urls = c.convertMapToDescriptionURL(resultMap)
	}
	for _, url := range urls {
		url.Short = c.url + "/" + url.Short
	}
	return urls, nil
}

func (c *Converter) mergedDeleteChan(ctx context.Context, chans ...<-chan *model.DescriptionURL) <-chan *model.DescriptionURL {
	out := make(chan *model.DescriptionURL)

	var wg sync.WaitGroup
	wg.Add(len(chans))

	for _, ch := range chans {
		ch := ch
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					select {
					case out <- v:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// DeleteURLsForUser удаляет URL для указанного пользователя.
func (c *Converter) DeleteURLsForUser(urls []string, user string) {
	for _, url := range urls {

		select {
		case c.httpDeleteChan <- &model.DescriptionURL{Short: url, UserID: user}:
		default:
			msg := "Буфер переполнен"
			logrus.Info(msg)
		}
	}
}

func (c *Converter) runProcessDeleteURL(ctx context.Context, deletes <-chan *model.DescriptionURL, batchSize int, flushTimeout time.Duration) {
	logrus.Info("Старт процесса обновления удалённых URL")

	batch := make([]*model.DescriptionURL, 0, batchSize)
	timer := time.NewTimer(flushTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Infoln("Отмена контекста в процессе обновления удалённых URL")
			c.DeleteURLs(batch)
			return

		case del, ok := <-deletes:
			if !ok {
				logrus.Infoln("Канал для удаления данных закрыт")
				c.DeleteURLs(batch)
				return
			}
			batch = append(batch, del)
			if len(batch) == cap(batch) {
				c.DeleteURLs(batch)
				batch = batch[:0]
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(flushTimeout)
			}

		case <-timer.C:
			c.DeleteURLs(batch)
			batch = batch[:0]
			timer.Reset(flushTimeout)
		}
	}
}

// DeleteURLs удаляет URL из хранилища.
func (c *Converter) DeleteURLs(urls []*model.DescriptionURL) error {
	switch c.modeStore {
	case ModeStoreDB:
		err := c.TranDeleteURLsFromDB(urls)
		if err != nil {
			logrus.Error("ошибка удаления из базы:", err)
			return err
		}

	case ModeStoreFile:
		err := c.DeleteURLsFromFile(urls)
		if err != nil {
			logrus.Error("ошибка удаления из файла:", err)
		}
	case ModeStoreStorage:
		for _, url := range urls {
			if c.userURLS.ExistURLForUser(url.UserID, url.Short) {
				c.storage.Delete(url.Short)
			} else {
				logrus.Infof("Не существует URL=%s у пользователя %s", url.Short, url.UserID)
			}
		}
	}
	return nil
}

// StartDeleteProcessor запускает процесс обновления удалённых URL.
func (c *Converter) StartDeleteProcessor(ctx context.Context, batchSize int, flushTimeout time.Duration,
) {

	mergedChan := c.mergedDeleteChan(ctx,
		c.httpDeleteChan,
		c.manualDeleteChan,
	)

	go c.runProcessDeleteURL(ctx, mergedChan, batchSize, flushTimeout)
}
