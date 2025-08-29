package core

import (
	"testing"
)

func TestNegativePointsError(t *testing.T) {
	err := &NegativePointsError{Points: -50}

	if err.Error() == "" {
		t.Error("NegativePointsError.Error() devrait retourner un message")
	}

	if err.Points != -50 {
		t.Errorf("Points = %d, attendu -50", err.Points)
	}
}

func TestCaptureError(t *testing.T) {
	err := &CaptureError{Word: "test", Reason: "mot invalide"}

	if err.Error() == "" {
		t.Error("CaptureError.Error() devrait retourner un message")
	}

	if err.Word != "test" {
		t.Errorf("Word = %q, attendu %q", err.Word, "test")
	}

	if err.Reason != "mot invalide" {
		t.Errorf("Reason = %q, attendu %q", err.Reason, "mot invalide")
	}
}

func TestInvalidAttemptError(t *testing.T) {
	err := &InvalidAttemptError{Input: "invalid", Reason: "format incorrect"}

	if err.Error() == "" {
		t.Error("InvalidAttemptError.Error() devrait retourner un message")
	}

	if err.Input != "invalid" {
		t.Errorf("Input = %q, attendu %q", err.Input, "invalid")
	}

	if err.Reason != "format incorrect" {
		t.Errorf("Reason = %q, attendu %q", err.Reason, "format incorrect")
	}
}

func TestInvalidStateError(t *testing.T) {
	err := &InvalidStateError{From: "IDLE", Expected: "ENCOUNTERED"}

	if err.Error() == "" {
		t.Error("InvalidStateError.Error() devrait retourner un message")
	}

	if err.From != "IDLE" {
		t.Errorf("From = %q, attendu %q", err.From, "IDLE")
	}

	if err.Expected != "ENCOUNTERED" {
		t.Errorf("Expected = %q, attendu %q", err.Expected, "ENCOUNTERED")
	}
}
