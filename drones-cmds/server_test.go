package dronescmds

import (
	"os"
	"testing"
)

func TestParseRabbitMqUrlFromEnv(t *testing.T) {
	os.Setenv("RABBITMQ_HOST", "rabbitmq")
	os.Setenv("RABBITMQ_PORT", "55555")

	url := parseRabbitURL()
	expected := "amqp://guest:guest@rabbitmq:55555"
	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
	os.Unsetenv("RABBITMQ_HOST")
	os.Unsetenv("RABBITMQ_PORT")
}

func TestParseRabbitMqUrlFallback(t *testing.T) {
	url := parseRabbitURL()
	expected := "amqp://guest:guest@localhost:5672"
	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}
