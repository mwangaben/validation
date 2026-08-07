package validation

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

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
func (v *Validator) validateUnique(value interface{}, data map[string]interface{}, table string, columns ...string) bool {
	if v.db == nil {
		return true
	}

	column := "id"
	exceptID := 0

	if len(columns) > 0 {
		column = columns[0]
	}

	if len(columns) > 1 {
		exceptColumn := columns[1]
		if val, exists := data[exceptColumn]; exists {
			switch v := val.(type) {
			case int:
				exceptID = v
			case int64:
				exceptID = int(v)
			case float64:
				exceptID = int(v)
			case string:
				if id, err := strconv.Atoi(v); err == nil {
					exceptID = id
				}
			}
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

	// Try to assert v.db as *gorm.DB
	db, ok := v.db.(*gorm.DB)
	if !ok {
		return true
	}

	var count int64
	query := db.Table(table).Where(column+" = ?", strValue)

	if exceptID > 0 {
		query = query.Where("id != ?", exceptID)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return true
	}

	return count == 0
}

// validateExists checks if value exists in database (with optional soft delete support)
func (v *Validator) validateExists(value interface{}, table string, columns ...string) bool {
	if v.db == nil {
		return false
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
		return false
	}

	db, ok := v.db.(*gorm.DB)
	if !ok {
		return false
	}

	var count int64
	query := db.Table(table).Where(column+" = ?", strValue)

	// Check if table has deleted_at column (soft delete support)
	if hasDeletedAtColumn(db, table) {
		query = query.Where("deleted_at IS NULL")
	}

	result := query.Count(&count)
	if result.Error != nil {
		return false
	}

	return count > 0
}

// hasDeletedAtColumn checks if a table has a deleted_at column
func hasDeletedAtColumn(db *gorm.DB, table string) bool {
	var count int64
	err := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND column_name = 'deleted_at'
	`, table).Count(&count).Error

	if err != nil {
		return false
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

// ExistsWithoutSoftDelete returns a rule that checks if a value exists in a table (including soft-deleted)
func ExistsWithoutSoftDelete(table, column string) string {
	return "exists_without_soft_delete:" + table + "," + column
}

// validateExistsWithoutSoftDelete checks if value exists in table (including soft-deleted)
func (v *Validator) validateExistsWithoutSoftDelete(value interface{}, table string, columns ...string) bool {
	if v.db == nil {
		return false
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
		return false
	}

	db, ok := v.db.(*gorm.DB)
	if !ok {
		return false
	}

	var count int64
	query := db.Table(table).Where(column+" = ?", strValue)
	result := query.Count(&count)
	if result.Error != nil {
		return false
	}

	return count > 0
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
