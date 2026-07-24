# Go Validator

A Laravel-style validation package for Go.

## Installation

```bash
go get github.com/mwangaben/validation
```
## Usage
#### Basic Validation
```go
import "github.com/mwangaben/validation"

validator := validator.NewValidator()

data := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
}

rules := map[string][]string{
    "name":  {validator.Required(), validator.String(), validator.Min(2), validator.Max(100)},
    "email": {validator.Required(), validator.Email()},
    "age":   {validator.Required(), validator.Int(), validator.Between(1, 150)},
}

if validator.Validate(data, rules) {
    // Validation passed
} else {
    // Validation failed
    fmt.Println(validator.Error())
}
```


#### Database Validation
```go
// With GORM
db := gorm.Open(...)
validator := validator.NewValidatorWithDB(db)

data := map[string]interface{}{
    "email": "john@example.com",
}

rules := map[string][]string{
    "email": {validator.Unique("users", "email")},
}

validator.Validate(data, rules)
```
#### Custom Validator
```go
type UserValidator struct {
    *validator.Validator
}

func NewUserValidator(db interface{}) *UserValidator {
    return &UserValidator{
        Validator: validator.NewValidatorWithDB(db),
    }
}

func (uv *UserValidator) ValidateCreateUser(data map[string]interface{}) error {
    rules := map[string][]string{
        "name":  {validator.Required(), validator.String(), validator.Min(2), validator.Max(100)},
        "email": {validator.Required(), validator.Email(), validator.Unique("users", "email")},
        "age":   {validator.Required(), validator.Int(), validator.Between(1, 150)},
    }

    if !uv.Validate(data, rules) {
        return fmt.Errorf(uv.Error())
    }
    return nil
}
```

#### 1. User Registration
```go
data := map[string]interface{}{
    "name":     "John Doe",
    "email":    "john@example.com",
    "password": "SecurePass123",
}

rules := map[string][]string{
    "name":     {"required", "string", "min:2", "max:100"},
    "email":    {"required", "email", "unique:users,email"},
    "password": {"required", "string", "min:8", "password"},
}
```

#### 2. User Update
```go
data := map[string]interface{}{
    "email": "john@example.com",
    "id":    1,
}

rules := map[string][]string{
    "email": {"required", "email", "unique:users,email,id"},
    "id":    {"exists:users,id"},
}
```

#### 3. Foreign Key Validation
```go
data := map[string]interface{}{
    "user_id":    1,
    "product_id": "P001",
    "quantity":   5,
}

rules := map[string][]string{
    "user_id":    {"exists:users,id"},
    "product_id": {"exists:products,code"},
    "quantity":   {"int", "min:1", "max:100"},
}
```

#### 4. API Request Validation

```go
func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    data := map[string]interface{}{
        "name":  req.Name,
        "email": req.Email,
    }
    
    rules := map[string][]string{
        "name":  {"required", "string", "min:2", "max:100"},
        "email": {"required", "email", "unique:users,email"},
    }
    
    if !validator.Validate(data, rules) {
        c.JSON(400, gin.H{"errors": validator.Errors()})
        return
    }
    
    // Create user...
}
```


### Available Rules

- **required** - Field must be present and not empty

* **email** - Field must be a valid email address

- **min:n** - Field must be at least n characters/bytes

- **max:n** - Field must be at most n characters/bytes

- **string** - Field must be a string

- **int** - Field must be an integer

- **numeric** - Field must be numeric

- **unique:table,column** - Field must be unique in the database

- **exists:table,column** - Field must exist in the database

- **in:value1,value2,...** - Field must be in the given list

- **not_in:value1,value2,...** - Field must not be in the given list

- **confirmed** - Field must match {field}_confirmation

- **date** - Field must be a valid date

- **url** - Field must be a valid URL

- **alpha** - Field must contain only letters

- **alpha_num** - Field must contain only letters and numbers

- **boolean** - Field must be a boolean

- **array** - Field must be an array

- **between:min,max** - Field must be between min and max

- **phone** - Field must be a valid phone number

- **password** - Field must be a strong password

- **uuid** - Field must be a valid UUID

























