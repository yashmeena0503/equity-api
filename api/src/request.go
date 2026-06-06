package src

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type PlatformType string
const (
	// NOTE: NEW_PLATFORM, wherever this note is mentioned, the new platform switch-cases are to be updated.
	PlatformDiscord PlatformType = "discord"
	PlatformMinecraft PlatformType = "minecraft"
)
func (p *PlatformType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch PlatformType(s) {
	// NOTE: NEW_PLATFORM
	case PlatformDiscord, PlatformMinecraft:
		*p = PlatformType(s); return nil
	default:
		return ErrUnknownPlatform(s)
	}
}

func WriteErrs(w http.ResponseWriter, status int, errs ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any { "errors": errs })
}

// This is like a template of objects.
type PlatformInfo struct {
	PlatformName PlatformType `json:"platform_name"`
	PlatformID string `json:"platform_id"`
}
// Call this to require all objects from the template struct in the request.
func (r PlatformInfo) validate() []string {
	var errs []string
	if r.PlatformName == "" { errs = append(errs, "platform_name is required") }
	if r.PlatformID == "" { errs = append(errs, "platform_id is required") }
	return errs
}

// HEAD: Request structs for Endpoints.

type MessageRequest struct {
	PlatformInfo
	MessageLength int `json:"message_length"`
}
func (r MessageRequest) validate() []string {
	errs := r.PlatformInfo.validate()
	if r.MessageLength <= 0 {errs = append(errs, "message_length must be greater than 0")}
	return errs
}

// End of Request structs.

func DecodeAndValidate[T interface{ validate() []string }](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrs(w, http.StatusBadRequest, "invalid JSON")
		return req, false
	}
	if errs := req.validate(); len(errs) > 0 {
		WriteErrs(w, http.StatusUnprocessableEntity, errs...)
		return req, false
	}
	return req, true
}

type ExpectedError struct {
	Status int
	Message string
}
func (e ExpectedError) Error() string { return e.Message }

func HandleRequest[
	T interface { validate() []string },
	R any,
](
	app *App,
	fn func(ctx context.Context, app *App, req T) (R, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := DecodeAndValidate[T](w, r)
		if !ok { return }
		result, err := fn(r.Context(), app, req)
		if err != nil {
			// Expected error as ExpectedError.
			if clientErr, ok := errors.AsType[ExpectedError](err); ok {
				WriteErrs(w, clientErr.Status, clientErr.Message)
				return
			}
			// Unexpected error being logged, and user being excused.
			slog.Error("HandleRequest error", "err", err)
			WriteErrs(w, 500, "Something went wrong!")
			return
		}
		// User being responded with 200.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
