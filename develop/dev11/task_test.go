package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// test the Event creation
func TestCreateEvent(t *testing.T) {
	calendarService := &CalendarServiceImpl{}
	handler := createEventHandler(calendarService)

	req, _ := http.NewRequest("POST", "/create_event", strings.NewReader("{\"id\":1,\"title\":\"Meeting\",\"description\":\"Discuss Project\",\"start_time\":\"2023-09-09T14:00:00Z\",\"end_time\":\"2023-09-09T15:00:00Z\"}"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// send a sample get request for events for day
func TestEventsForDay(t *testing.T) {
	calendarService := &CalendarServiceImpl{}
	handler := eventsForDayHandler(calendarService)

	req, _ := http.NewRequest("GET", "/events_for_day?date=2023-09-09", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// create and delete an event
func TestDeleteEvent(t *testing.T) {
	calendarService := &CalendarServiceImpl{}
	handler := createEventHandler(calendarService)

	req, _ := http.NewRequest("POST", "/create_event", strings.NewReader("{\"id\":1,\"title\":\"Meeting\",\"description\":\"Discuss Project\",\"start_time\":\"2023-09-09T14:00:00Z\",\"end_time\":\"2023-09-09T15:00:00Z\"}"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	deleteHandler := deleteEventHandler(calendarService)
	deleteReq, _ := http.NewRequest("POST", "/delete_event", strings.NewReader("{\"id\":1}"))
	deleteReq.Header.Set("Content-Type", "application/json")
	deleteResp := httptest.NewRecorder()

	deleteHandler.ServeHTTP(deleteResp, deleteReq)

	if status := deleteResp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
