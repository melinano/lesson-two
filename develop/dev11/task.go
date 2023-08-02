package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем.
В рамках задания необходимо работать строго со стандартной HTTP библиотекой.

В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов

Методы API:
	POST /create_event
	POST /update_event
	POST /delete_event
	GET /events_for_day
	GET /events_for_week
	GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."}
в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных
	   (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер
	   должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить
	   в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// Event structure to hold event data
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

// CalendarService - Interface for the service operations
type CalendarService interface {
	CreateEvent(event Event) error
	UpdateEvent(event Event) error
	DeleteEvent(id int) error
	EventsForDay(day time.Time) ([]Event, error)
	EventsForWeek(start time.Time) ([]Event, error)
	EventsForMonth(start time.Time) ([]Event, error)
}

// CalendarServiceImpl - concrete type implementing CalendarService
type CalendarServiceImpl struct {
	events []Event
	mu     sync.Mutex
}

func main() {
	calendarService := &CalendarServiceImpl{}

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", createEventHandler(calendarService))
	mux.HandleFunc("/update_event", updateEventHandler(calendarService))
	mux.HandleFunc("/delete_event", deleteEventHandler(calendarService))
	mux.HandleFunc("/events_for_day", eventsForDayHandler(calendarService))
	mux.HandleFunc("/events_for_week", eventsForWeekHandler(calendarService))
	mux.HandleFunc("/events_for_month", eventsForMonthHandler(calendarService))

	log.Fatal(http.ListenAndServe(":8080", loggingMiddleware(mux)))
}

func createEventHandler(service CalendarService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.CreateEvent(event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func updateEventHandler(service CalendarService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.UpdateEvent(event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteEventHandler(service CalendarService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.DeleteEvent(event.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func eventsForDayHandler(service CalendarService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		day, err := time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events, err := service.EventsForDay(day)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func eventsForWeekHandler(service CalendarService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		start, err := time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events, err := service.EventsForWeek(start)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func eventsForMonthHandler(service CalendarService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		start, err := time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events, err := service.EventsForMonth(start)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CreateEvent creates a new event in the calendar
func (c *CalendarServiceImpl) CreateEvent(event Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.events = append(c.events, event)
	return nil
}

// UpdateEvent updates an existing event in the calendar
func (c *CalendarServiceImpl) UpdateEvent(event Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, e := range c.events {
		if e.ID == event.ID {
			c.events[i] = event
			return nil
		}
	}

	return errors.New("event not found")
}

// DeleteEvent deletes an event from the calendar
func (c *CalendarServiceImpl) DeleteEvent(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, e := range c.events {
		if e.ID == id {
			// Delete the event without preserving order
			c.events[i] = c.events[len(c.events)-1]
			c.events = c.events[:len(c.events)-1]
			return nil
		}
	}

	return errors.New("event not found")
}

// EventsForDay returns all events for a specific day
func (c *CalendarServiceImpl) EventsForDay(day time.Time) ([]Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []Event
	for _, e := range c.events {
		if isSameDay(e.StartTime, day) {
			events = append(events, e)
		}
	}

	return events, nil
}

// EventsForWeek returns all events for a specific week
func (c *CalendarServiceImpl) EventsForWeek(start time.Time) ([]Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []Event
	for _, e := range c.events {
		if isSameWeek(e.StartTime, start) {
			events = append(events, e)
		}
	}

	return events, nil
}

// EventsForMonth returns all events for a specific month
func (c *CalendarServiceImpl) EventsForMonth(start time.Time) ([]Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []Event
	for _, e := range c.events {
		if isSameMonth(e.StartTime, start) {
			events = append(events, e)
		}
	}

	return events, nil
}

// helper function for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// isSameDay checks if two dates are in the same day
func isSameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

// isSameWeek checks if two dates are in the same week
func isSameWeek(a, b time.Time) bool {
	y1, w1 := a.ISOWeek()
	y2, w2 := b.ISOWeek()

	return y1 == y2 && w1 == w2
}

// isSameMonth checks if two dates are in the same month
func isSameMonth(a, b time.Time) bool {
	y1, m1, _ := a.Date()
	y2, m2, _ := b.Date()

	return y1 == y2 && m1 == m2
}
