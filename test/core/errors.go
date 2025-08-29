package core

import "fmt"

type InvalidStateError struct{ From, Expected string }

func (e *InvalidStateError) Error() string {
	return fmt.Sprintf("transition interdite: état=%s, attendu=%s", e.From, e.Expected)
}

type InvalidAttemptError struct{ Input, Reason string }

func (e *InvalidAttemptError) Error() string {
	return fmt.Sprintf("tentative invalide: %q (%s)", e.Input, e.Reason)
}

type CaptureError struct{ Word, Reason string }

func (e *CaptureError) Error() string {
	return fmt.Sprintf("impossible de capturer %q (%s)", e.Word, e.Reason)
}

type NegativePointsError struct{ Points int }

func (e *NegativePointsError) Error() string {
	return fmt.Sprintf("points négatifs interdits (%d)", e.Points)
}
