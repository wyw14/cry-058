package domain

import "fmt"

type CodeError struct{ Code, Message, Field string }

func (e CodeError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Field)
	}
	return e.Code + ": " + e.Message
}
func Invalid(msg, field string) error { return CodeError{"invalid", msg, field} }
func NotFound(kind, id string) error {
	return CodeError{"not_found", kind + " " + id + " 不存在", ""}
}
func Conflict(msg string) error   { return CodeError{"conflict", msg, ""} }
func Forbidden(msg string) error  { return CodeError{"forbidden", msg, ""} }
func StateError(msg string) error { return CodeError{"state_conflict", msg, ""} }
