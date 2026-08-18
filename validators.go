package validation

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// File represents an uploaded file
type UploadedFile struct {
	Name        string
	Size        int64
	ContentType string
	Path        string // Optional: path to temporary file
	Header      *multipart.FileHeader
}

// isFile checks if value is a File type
func isFile(value interface{}) bool {
	_, ok := value.(*UploadedFile)
	return ok
}

// getFileValue returns the File object if value is a file
func getFileValue(value interface{}) (*UploadedFile, bool) {
	file, ok := value.(*UploadedFile)
	return file, ok
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
	case reflect.Ptr:
		return !val.IsNil()
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

// ... [existing validation methods remain unchanged] ...

// validateDecimal checks if value has the required number of decimal places
func (v *Validator) validateDecimal(value interface{}, minStr string, maxStr ...string) bool {
	minDecimals, err := strconv.Atoi(minStr)
	if err != nil {
		return true
	}

	maxDecimals := minDecimals
	if len(maxStr) > 0 {
		maxDecimals, err = strconv.Atoi(maxStr[0])
		if err != nil {
			return true
		}
	}

	var strValue string
	switch val := value.(type) {
	case float32:
		strValue = fmt.Sprintf("%f", val)
	case float64:
		strValue = fmt.Sprintf("%f", val)
	case string:
		strValue = val
	default:
		return false
	}

	// Check if it's a valid number
	if _, err := strconv.ParseFloat(strValue, 64); err != nil {
		return false
	}

	// Count decimal places
	parts := strings.Split(strValue, ".")
	if len(parts) != 2 {
		return minDecimals == 0
	}

	decimalPlaces := len(strings.TrimRight(parts[1], "0"))
	return decimalPlaces >= minDecimals && decimalPlaces <= maxDecimals
}

// validateDimensions checks image dimensions
// validateDimensions checks image dimensions
func (v *Validator) validateDimensions(value interface{}, params string) bool {
	file, ok := getFileValue(value)
	if !ok {
		return false
	}

	// Parse dimensions parameters
	paramsMap := make(map[string]int)
	for _, param := range strings.Split(params, ",") {
		parts := strings.Split(param, "=")
		if len(parts) == 2 {
			if val, err := strconv.Atoi(parts[1]); err == nil {
				paramsMap[parts[0]] = val
			}
		}
	}

	var reader io.Reader
	var closer io.Closer

	if file.Path != "" {
		imgFile, err := os.Open(file.Path)
		if err != nil {
			return false
		}
		reader = imgFile
		closer = imgFile
	} else if file.Header != nil {
		multipartFile, err := file.Header.Open()
		if err != nil {
			return false
		}
		reader = multipartFile
		closer = multipartFile
	} else {
		return false
	}

	if closer != nil {
		defer closer.Close()
	}

	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return false
	}

	if minWidth, ok := paramsMap["min_width"]; ok {
		if img.Width < minWidth {
			return false
		}
	}

	if minHeight, ok := paramsMap["min_height"]; ok {
		if img.Height < minHeight {
			return false
		}
	}

	if maxWidth, ok := paramsMap["max_width"]; ok {
		if img.Width > maxWidth {
			return false
		}
	}

	if maxHeight, ok := paramsMap["max_height"]; ok {
		if img.Height > maxHeight {
			return false
		}
	}

	return true
}

// getSizeValue returns the size based on the type conventions
func getSizeValue(value interface{}) (int, bool) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(val.String())), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(val.Uint()), true
	case reflect.Float32, reflect.Float64:
		return int(val.Float()), true
	case reflect.Array, reflect.Slice, reflect.Map:
		return val.Len(), true
	case reflect.Ptr:
		if file, ok := value.(*UploadedFile); ok {
			return int(file.Size / 1024), true // Size in kilobytes
		}
	}
	return 0, false
}

// validateGt checks if value is greater than the given field or value
func (v *Validator) validateGt(value interface{}, fieldOrValue string) bool {
	actualSize, ok := getSizeValue(value)
	if !ok {
		return false
	}

	// Check if it's a comparison with another field
	if otherValue, exists := v.data[fieldOrValue]; exists {
		otherSize, ok := getSizeValue(otherValue)
		if !ok {
			return false
		}
		return actualSize > otherSize
	}

	// Compare with a direct value
	compareValue, err := strconv.Atoi(fieldOrValue)
	if err != nil {
		return false
	}
	return actualSize > compareValue
}

// validateGte checks if value is greater than or equal to the given field or value
func (v *Validator) validateGte(value interface{}, fieldOrValue string) bool {
	actualSize, ok := getSizeValue(value)
	if !ok {
		return false
	}

	if otherValue, exists := v.data[fieldOrValue]; exists {
		otherSize, ok := getSizeValue(otherValue)
		if !ok {
			return false
		}
		return actualSize >= otherSize
	}

	compareValue, err := strconv.Atoi(fieldOrValue)
	if err != nil {
		return false
	}
	return actualSize >= compareValue
}

// validateLt checks if value is less than the given field or value
func (v *Validator) validateLt(value interface{}, fieldOrValue string) bool {
	actualSize, ok := getSizeValue(value)
	if !ok {
		return false
	}

	if otherValue, exists := v.data[fieldOrValue]; exists {
		otherSize, ok := getSizeValue(otherValue)
		if !ok {
			return false
		}
		return actualSize < otherSize
	}

	compareValue, err := strconv.Atoi(fieldOrValue)
	if err != nil {
		return false
	}
	return actualSize < compareValue
}

// validateLte checks if value is less than or equal to the given field or value
func (v *Validator) validateLte(value interface{}, fieldOrValue string) bool {
	actualSize, ok := getSizeValue(value)
	if !ok {
		return false
	}

	if otherValue, exists := v.data[fieldOrValue]; exists {
		otherSize, ok := getSizeValue(otherValue)
		if !ok {
			return false
		}
		return actualSize <= otherSize
	}

	compareValue, err := strconv.Atoi(fieldOrValue)
	if err != nil {
		return false
	}
	return actualSize <= compareValue
}

// validateFile checks if value is a successfully uploaded file
func (v *Validator) validateFile(value interface{}) bool {
	if _, ok := getFileValue(value); !ok {
		return false
	}
	return true
}

// validateExtensions checks if file has the required extension
func (v *Validator) validateExtensions(value interface{}, extensions ...string) bool {
	file, ok := getFileValue(value)
	if !ok {
		return false
	}

	ext := strings.ToLower(filepath.Ext(file.Name))
	if ext == "" {
		return false
	}
	ext = strings.TrimPrefix(ext, ".")

	for _, allowedExt := range extensions {
		if strings.ToLower(allowedExt) == ext {
			return true
		}
	}
	return false
}

// validateMimes checks if file has the required MIME type based on extension
func (v *Validator) validateMimes(value interface{}, mimes ...string) bool {
	file, ok := getFileValue(value)
	if !ok {
		return false
	}

	// Get extension from filename
	ext := strings.ToLower(filepath.Ext(file.Name))
	if ext == "" {
		return false
	}
	ext = strings.TrimPrefix(ext, ".")

	// Check if the file extension is in the allowed list
	extensionAllowed := false
	for _, allowedExt := range mimes {
		if strings.ToLower(allowedExt) == ext {
			extensionAllowed = true
			break
		}
	}

	if !extensionAllowed {
		return false
	}

	// If ContentType is provided, verify it matches
	if file.ContentType != "" {
		// Map extension to expected MIME type
		extensionMimeMap := map[string]string{
			"jpg":  "image/jpeg",
			"jpeg": "image/jpeg",
			"png":  "image/png",
			"gif":  "image/gif",
			"bmp":  "image/bmp",
			"pdf":  "application/pdf",
			"txt":  "text/plain",
			"csv":  "text/csv",
			"json": "application/json",
			"xml":  "application/xml",
			"zip":  "application/zip",
		}

		expectedMime, exists := extensionMimeMap[ext]
		if exists {
			// Check if the actual MIME type matches the expected one
			return strings.ToLower(file.ContentType) == expectedMime
		}
	}

	// If no ContentType provided, just check the extension
	return true
}

// validateMimetypes checks if file has the required MIME type
func (v *Validator) validateMimetypes(value interface{}, mimetypes ...string) bool {
	file, ok := getFileValue(value)
	if !ok {
		return false
	}

	actualMime := file.ContentType
	if actualMime == "" {
		return false
	}

	for _, allowedMime := range mimetypes {
		// Handle wildcard patterns like image/*
		if strings.HasSuffix(allowedMime, "/*") {
			prefix := strings.TrimSuffix(allowedMime, "*")
			if strings.HasPrefix(actualMime, prefix) {
				return true
			}
		} else if strings.ToLower(allowedMime) == strings.ToLower(actualMime) {
			return true
		}
	}
	return false
}

// validateSize checks if value has the required size
func (v *Validator) validateSize(value interface{}, sizeStr string) bool {
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return true
	}

	actualSize, ok := getSizeValue(value)
	if !ok {
		return false
	}

	return actualSize == size
}

// ... [rest of existing validation methods remain unchanged] ...

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
