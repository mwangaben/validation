package validation

import (
	"fmt"
	"strings"
)

// Common validation rules
const (
	RuleRequired                = "required"
	RuleEmail                   = "email"
	RuleMin                     = "min"
	RuleMax                     = "max"
	RuleString                  = "string"
	RuleInt                     = "int"
	RuleNumeric                 = "numeric"
	RuleUnique                  = "unique"
	RuleExists                  = "exists"
	RuleIn                      = "in"
	RuleNotIn                   = "not_in"
	RuleConfirmed               = "confirmed"
	RuleDate                    = "date"
	RuleURL                     = "url"
	RuleAlpha                   = "alpha"
	RuleAlphaNum                = "alpha_num"
	RuleBoolean                 = "boolean"
	RuleArray                   = "array"
	RuleBetween                 = "between"
	RulePhone                   = "phone"
	RulePassword                = "password"
	RuleUUID                    = "uuid"
	RuleMaxLength               = "max_length"
	RuleMinLength               = "min_length"
	ExistsWithoutSoftDeleteRule = "exists_without_soft_delete"
	RuleDecimal                 = "decimal"
	RuleDimensions              = "dimensions"
	RuleGt                      = "gt"
	RuleGte                     = "gte"
	RuleLt                      = "lt"
	RuleLte                     = "lte"
	RuleFile                    = "file"
	RuleExtensions              = "extensions"
	RuleMimes                   = "mimes"
	RuleMimetypes               = "mimetypes"
	RuleSize                    = "size"
)

// Helper functions to build rules

// Required returns required rule
func Required() string {
	return RuleRequired
}

// Email returns email rule
func Email() string {
	return RuleEmail
}

// Min returns min rule with value
func Min(value int) string {
	return fmt.Sprintf("%s:%d", RuleMin, value)
}

// Max returns max rule with value
func Max(value int) string {
	return fmt.Sprintf("%s:%d", RuleMax, value)
}

// String returns string rule
func String() string {
	return RuleString
}

// Int returns int rule
func Int() string {
	return RuleInt
}

// Numeric returns numeric rule
func Numeric() string {
	return RuleNumeric
}

// Unique returns unique rule for table and column
func Unique(table string, column ...string) string {
	col := "id"
	if len(column) > 0 {
		col = column[0]
	}
	return fmt.Sprintf("%s:%s,%s", RuleUnique, table, col)
}

// UniqueExcept returns unique rule with exception for a specific ID
func UniqueExcept(table, column string, exceptID int) string {
	return fmt.Sprintf("%s:%s,%s,%d", RuleUnique, table, column, exceptID)
}

// Exists returns exists rule for table and column
func Exists(table string, column ...string) string {
	col := "id"
	if len(column) > 0 {
		col = column[0]
	}
	return fmt.Sprintf("%s:%s,%s", RuleExists, table, col)
}

// In returns in rule with values
func In(values ...string) string {
	return fmt.Sprintf("%s:%s", RuleIn, strings.Join(values, ","))
}

// NotIn returns not_in rule with values
func NotIn(values ...string) string {
	return fmt.Sprintf("%s:%s", RuleNotIn, strings.Join(values, ","))
}

// Confirmed returns confirmed rule
func Confirmed() string {
	return RuleConfirmed
}

// Date returns date rule
func Date() string {
	return RuleDate
}

// URL returns URL rule
func URL() string {
	return RuleURL
}

// Alpha returns alpha rule
func Alpha() string {
	return RuleAlpha
}

// AlphaNum returns alpha_num rule
func AlphaNum() string {
	return RuleAlphaNum
}

// Boolean returns boolean rule
func Boolean() string {
	return RuleBoolean
}

// Array returns array rule
func Array() string {
	return RuleArray
}

// Between returns between rule with min and max
func Between(min, max int) string {
	return fmt.Sprintf("%s:%d,%d", RuleBetween, min, max)
}

// Phone returns phone rule
func Phone() string {
	return RulePhone
}

// Password returns password rule
func Password() string {
	return RulePassword
}

// UUID returns uuid rule
func UUID() string {
	return RuleUUID
}

// MaxLength returns max length rule
func MaxLength(length int) string {
	return fmt.Sprintf("%s:%d", RuleMaxLength, length)
}

// MinLength returns min length rule
func MinLength(length int) string {
	return fmt.Sprintf("%s:%d", RuleMinLength, length)
}

// Decimal returns decimal rule with min and optional max decimal places
func Decimal(min int, max ...int) string {
	if len(max) > 0 {
		return fmt.Sprintf("%s:%d,%d", RuleDecimal, min, max[0])
	}
	return fmt.Sprintf("%s:%d", RuleDecimal, min)
}

// Dimensions returns dimensions rule
func Dimensions(minWidth, minHeight int) string {
	return fmt.Sprintf("%s:min_width=%d,min_height=%d", RuleDimensions, minWidth, minHeight)
}

// Gt returns greater than rule
func Gt(field string) string {
	return fmt.Sprintf("%s:%s", RuleGt, field)
}

// Gte returns greater than or equal rule
func Gte(field string) string {
	return fmt.Sprintf("%s:%s", RuleGte, field)
}

// Lt returns less than rule
func Lt(field string) string {
	return fmt.Sprintf("%s:%s", RuleLt, field)
}

// Lte returns less than or equal rule
func Lte(field string) string {
	return fmt.Sprintf("%s:%s", RuleLte, field)
}

// File returns file rule
func File() string {
	return RuleFile
}

// Extensions returns extensions rule
func Extensions(extensions ...string) string {
	return fmt.Sprintf("%s:%s", RuleExtensions, strings.Join(extensions, ","))
}

// Mimes returns mimes rule
func Mimes(mimes ...string) string {
	return fmt.Sprintf("%s:%s", RuleMimes, strings.Join(mimes, ","))
}

// Mimetypes returns mimetypes rule
func Mimetypes(mimetypes ...string) string {
	return fmt.Sprintf("%s:%s", RuleMimetypes, strings.Join(mimetypes, ","))
}

// Size returns size rule
func Size(size int) string {
	return fmt.Sprintf("%s:%d", RuleSize, size)
}
