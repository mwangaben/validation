package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DBInterface defines the database interface needed for validation
type DBInterface interface {
	Table(name string) interface{}
	Where(query interface{}, args ...interface{}) interface{}
	Count(count *int64) interface{}
}

// validateRequired checks if value is not empty
func (v *Validator) validateRequired(value interface{}) bool {
	if value == nil {
		return false
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return strings.TrimSpace(val.String()) != ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0
	case reflect.Bool:
		return true
	case reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() > 0
	default:
		return true
	}
}

// validateEmail checks if value is a valid email
func (v *Validator) validateEmail(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	_, err := mail.ParseAddress(str)
	return err == nil
}

// validateMin checks if value meets minimum requirement
func (v *Validator) validateMin(value interface{}, minStr string) bool {
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(val.String())) >= min
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int()) >= min
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(val.Uint()) >= min
	case reflect.Float32, reflect.Float64:
		return int(val.Float()) >= min
	case reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() >= min
	default:
		return true
	}
}

// validateMax checks if value meets maximum requirement
func (v *Validator) validateMax(value interface{}, maxStr string) bool {
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(val.String())) <= max
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int()) <= max
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(val.Uint()) <= max
	case reflect.Float32, reflect.Float64:
		return int(val.Float()) <= max
	case reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() <= max
	default:
		return true
	}
}

// validateString checks if value is a string
func (v *Validator) validateString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// validateInt checks if value is an integer
func (v *Validator) validateInt(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return false
	default:
		return false
	}
}

// validateNumeric checks if value is numeric
func (v *Validator) validateNumeric(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// validateUnique checks if value is unique in database
func (v *Validator) validateUnique(value interface{}, table string, columns ...string) bool {
	if v.db == nil {
		return true
	}

	column := "id"
	exceptID := 0
	if len(columns) > 0 {
		column = columns[0]
	}
	if len(columns) > 1 {
		if id, err := strconv.Atoi(columns[1]); err == nil {
			exceptID = id
		}
	}

	var strValue string
	switch val := value.(type) {
	case string:
		strValue = val
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		strValue = fmt.Sprintf("%v", val)
	default:
		return true
	}

	// Use reflection to call database methods
	dbVal := reflect.ValueOf(v.db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}

	// Get the table
	tableMethod := dbVal.MethodByName("Table")
	if !tableMethod.IsValid() {
		return true
	}
	tableResult := tableMethod.Call([]reflect.Value{reflect.ValueOf(table)})
	if len(tableResult) == 0 {
		return true
	}

	// Get the where clause
	whereMethod := tableResult[0].MethodByName("Where")
	if !whereMethod.IsValid() {
		return true
	}

	query := fmt.Sprintf("%s = ?", column)
	whereResult := whereMethod.Call([]reflect.Value{
		reflect.ValueOf(query),
		reflect.ValueOf(strValue),
	})
	if len(whereResult) == 0 {
		return true
	}

	// If exceptID is provided, exclude that record
	if exceptID > 0 {
		whereMethod2 := whereResult[0].MethodByName("Where")
		if whereMethod2.IsValid() {
			whereResult[0] = whereMethod2.Call([]reflect.Value{
				reflect.ValueOf("id != ?"),
				reflect.ValueOf(exceptID),
			})[0]
		}
	}

	// Count
	countMethod := whereResult[0].MethodByName("Count")
	if !countMethod.IsValid() {
		return true
	}
	var count int64
	countResult := countMethod.Call([]reflect.Value{reflect.ValueOf(&count)})
	if len(countResult) == 0 {
		return true
	}

	return count == 0
}

// validateExists checks if value exists in database
func (v *Validator) validateExists(value interface{}, table string, columns ...string) bool {
	if v.db == nil {
		return true
	}

	column := "id"
	if len(columns) > 0 {
		column = columns[0]
	}

	var strValue string
	switch val := value.(type) {
	case string:
		strValue = val
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		strValue = fmt.Sprintf("%v", val)
	default:
		return true
	}

	// Use reflection to call database methods
	dbVal := reflect.ValueOf(v.db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}

	// Get the table
	tableMethod := dbVal.MethodByName("Table")
	if !tableMethod.IsValid() {
		return true
	}
	tableResult := tableMethod.Call([]reflect.Value{reflect.ValueOf(table)})
	if len(tableResult) == 0 {
		return true
	}

	// Get the where clause
	whereMethod := tableResult[0].MethodByName("Where")
	if !whereMethod.IsValid() {
		return true
	}

	query := fmt.Sprintf("%s = ?", column)
	whereResult := whereMethod.Call([]reflect.Value{
		reflect.ValueOf(query),
		reflect.ValueOf(strValue),
	})
	if len(whereResult) == 0 {
		return true
	}

	// Count
	countMethod := whereResult[0].MethodByName("Count")
	if !countMethod.IsValid() {
		return true
	}
	var count int64
	countResult := countMethod.Call([]reflect.Value{reflect.ValueOf(&count)})
	if len(countResult) == 0 {
		return true
	}

	return count > 0
}

// validateIn checks if value is in a list
func (v *Validator) validateIn(value interface{}, allowed []string) bool {
	if value == nil {
		return false
	}
	str := fmt.Sprintf("%v", value)
	for _, allowedVal := range allowed {
		if str == allowedVal {
			return true
		}
	}
	return false
}

// validateNotIn checks if value is not in a list
func (v *Validator) validateNotIn(value interface{}, disallowed []string) bool {
	if value == nil {
		return true
	}
	str := fmt.Sprintf("%v", value)
	for _, disallowedVal := range disallowed {
		if str == disallowedVal {
			return false
		}
	}
	return true
}

// validateConfirmed checks if value matches confirmation field
func (v *Validator) validateConfirmed(field string, value interface{}) bool {
	confirmationField := field + "_confirmation"
	confirmationValue, exists := v.data[confirmationField]
	if !exists {
		return false
	}
	return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", confirmationValue)
}

// validateDate checks if value is a valid date
func (v *Validator) validateDate(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
	}
	for _, format := range formats {
		if _, err := time.Parse(format, str); err == nil {
			return true
		}
	}
	return false
}

// validateURL checks if value is a valid URL
func (v *Validator) validateURL(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	_, err := url.ParseRequestURI(str)
	return err == nil
}

// validateAlpha checks if value contains only letters
func (v *Validator) validateAlpha(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z]+$", str)
	return matched
}

// validateAlphaNum checks if value contains only letters and numbers
func (v *Validator) validateAlphaNum(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", str)
	return matched
}

// validateBoolean checks if value is boolean
func (v *Validator) validateBoolean(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// validateArray checks if value is an array/slice
func (v *Validator) validateArray(value interface{}) bool {
	val := reflect.ValueOf(value)
	return val.Kind() == reflect.Array || val.Kind() == reflect.Slice
}

// validateBetween checks if value is between min and max
func (v *Validator) validateBetween(value interface{}, minStr, maxStr string) bool {
	min, err1 := strconv.Atoi(minStr)
	max, err2 := strconv.Atoi(maxStr)
	if err1 != nil || err2 != nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		length := len(strings.TrimSpace(val.String()))
		return length >= min && length <= max
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num := int(val.Int())
		return num >= min && num <= max
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num := int(val.Uint())
		return num >= min && num <= max
	case reflect.Float32, reflect.Float64:
		num := int(val.Float())
		return num >= min && num <= max
	case reflect.Array, reflect.Slice, reflect.Map:
		length := val.Len()
		return length >= min && length <= max
	default:
		return true
	}
}

// validatePhone checks if value is a valid phone number
func (v *Validator) validatePhone(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	// Basic phone validation (international format)
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, str)
	return matched
}

// validatePassword checks if password meets complexity requirements
func (v *Validator) validatePassword(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	if len(str) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(str)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(str)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(str)
	return hasUpper && hasLower && hasNumber
}

// validateUUID checks if value is a valid UUID
func (v *Validator) validateUUID(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, str)
	return matched
}
