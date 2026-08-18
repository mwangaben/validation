# Go Validation Package for User Input

A Laravel-style validation package for Go.

## Installation

```bash
go get github.com/mwangaben/validation
```
## Usage
#### Basic Validation
```go
import "github.com/mwangaben/validation"

validation := validation.NewValidator()

data := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
}

rules := map[string][]string{
    "name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
    "email": {validation.Required(), validation.Email()},
    "age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
}

if validation.Validate(data, rules) {
    // Validation passed
} else {
    // Validation failed
    fmt.Println(validation.Error())
}
```


#### Database Validation
```go
// With GORM
db := gorm.Open(...)
validation := validation.NewValidatorWithDB(db)

data := map[string]interface{}{
    "email": "john@example.com",
}

rules := map[string][]string{
    "email": {validation.Unique("users", "email")},
}

validation.Validate(data, rules)
```
#### Custom Validator
```go
type UserValidator struct {
    *validation.Validator
}

func NewUserValidator(db interface{}) *UserValidator {
    return &UserValidator{
        Validation: validation.NewValidatorWithDB(db),
    }
}

func (uv *UserValidator) ValidateCreateUser(data map[string]interface{}) error {
    rules := map[string][]string{
        "name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
        "email": {validation.Required(), validation.Email(), validation.Unique("users", "email")},
        "age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
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
    
    if !validation.Validate(data, rules) {
        c.JSON(400, gin.H{"errors": validator.Errors()})
        return
    }
    
    // Create user...
}
```

#### 5. Decimal Validation

```go
data := map[string]interface{}{
    "price": 99.99,
    "rate":  3.456,
}

rules := map[string][]string{
    // Must have exactly 2 decimal places (9.99)...
    "price": {"required", "decimal:2"},
    
    // Must have between 2 and 4 decimal places...
    "rate": {"required", "decimal:2,4"},
}
```

#### 6. File Upload Validation

```go
file := &validation.UploadedFile{
    Name:        "avatar.jpg",
    Size:        1024 * 1024, // 1MB
    ContentType: "image/jpeg",
    Path:        "/tmp/uploaded_file.jpg", // Optional: for dimension checking
}

data := map[string]interface{}{
    "avatar": file,
    "photo":  file,
    "video":  file,
}

rules := map[string][]string{
    // Must be a successfully uploaded file...
    "avatar": {"file", "image", "dimensions:min_width=100,min_height=200"},
    
    // Must have a user-assigned extension...
    "photo": {"extensions:jpg,png,gif"},
    
    // Must have a MIME type corresponding to one of the listed extensions...
    "photo": {"mimes:jpg,bmp,png"},
    
    // Must match one of the given MIME types...
    "video": {"mimetypes:video/avi,video/mpeg,video/quicktime"},
    
    // Must be 512 kilobytes...
    "avatar": {"size:512"},
}
```

#### 7. Comparison Rules

```go
data := map[string]interface{}{
    "age":      25,
    "min_age":  18,
    "max_age":  65,
    "quantity": 10,
    "min_qty":  5,
}

rules := map[string][]string{
    // Must be greater than the given field or value...
    "age":      {"gt:min_age"},
    "age":      {"gt:21"},
    
    // Must be greater than or equal to the given field or value...
    "age":      {"gte:min_age"},
    
    // Must be less than the given field or value...
    "age":      {"lt:max_age"},
    
    // Must be less than or equal to the given field or value...
    "quantity": {"lte:max_qty"},
}
```

#### 8. Size Validation

```go
data := map[string]interface{}{
    "title": "Hello World", // 11 characters
    "seats": 10,
    "tags":  []string{"a", "b", "c", "d", "e"},
    "image": file, // 512 kilobytes
}

rules := map[string][]string{
    // Validate that a string is exactly 12 characters long...
    "title": {"size:12"},
    
    // Validate that a provided integer equals 10...
    "seats": {"integer", "size:10"},
    
    // Validate that an array has exactly 5 elements...
    "tags": {"array", "size:5"},
    
    // Validate that an uploaded file is exactly 512 kilobytes...
    "image": {"file", "size:512"},
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

##### Database Rules
- **unique:table,column** - Field must be unique in the database

- **unique:table,column,except_id** - Field must be unique except for the given ID

- **exists:table,column** - Field must exist in the database

- **exists_without_soft_delete:table,column** - Field must exist (including soft-deleted records)

##### List Rules

- **in:value1,value2,...** - Field must be in the given list

- **not_in:value1,value2,...** - Field must not be in the given list


#### Numeric Rules
- **decimal:min** - Field must have exactly min decimal places
   - Example: **decimal:2*** validates 9.99 but not 9.999

- **decimal:min,max** - Field must have between min and max decimal places
   - Example: decimal:2,4 validates 9.99, 9.999, and 9.9999

### Comparison Rules
- **gt:field** - Field must be greater than the given field
- **gt:value** - Field must be greater than the given value
- **gte:field** - Field must be greater than or equal to the given field
- **gte:value** - Field must be greater than or equal to the given value
- **lt:field** - Field must be less than the given field
- **lt:value** - Field must be less than the given value
- **lte:field** - Field must be less than or equal to the given field
- **lte:value** - Field must be less than or equal to the given value


**Note**: For comparison rules, the two fields must be of the same type. Strings, numerics, arrays, and files are evaluated using the same conventions as the size rule.

#### File Rules

- **file** - Field must be a successfully uploaded file
- **extensions:foo,bar,...** - File must have a user-assigned extension corresponding to one of the listed extensions
  - Example: extensions:jpg,png,gif
- **mimes:foo,bar,...** - File must have a MIME type corresponding to one of the listed extensions
  - Example: mimes:jpg,bmp,png
- **mimetypes:type/type,...** - File must match one of the given MIME types
   - Example: mimetypes:video/avi,video/mpeg,video/quicktime
   - Supports wildcards: mimetypes:image/*,video/*
- **dimensions:min_width=100,min_height=200** - Image must meet the given dimension constraints
   - Available parameters: min_width, min_height, max_width, max_height
   - Example: dimensions:min_width=100,min_height=200

#### Size Rule
- **size:value** - Field must have a size matching the given value
  - For string data: value corresponds to the number of characters
  - For numeric data: value corresponds to a given integer value (must have numeric or integer rule)
  - For arrays: size corresponds to the count of the array
  - For files: size corresponds to the file size in kilobytes
  - Examples:
    - size:12 - String must be exactly 12 characters
    - size:10 - Integer must equal 10
    - size:5 - Array must have exactly 5 elements
    - size:512 - File must be exactly 512 kilobytes
  
### File Handling
The package includes an UploadedFile struct for file validation:

```go
type UploadedFile struct {
    Name        string                 // Original filename
    Size        int64                  // File size in bytes
    ContentType string                 // MIME type of the file
    Path        string                 // Optional: path to temporary file
    Header      *multipart.FileHeader  // Optional: for multipart form data
}
```

### Helper Functions
The package provides helper functions to build validation rules:

```go
// Basic rules
validation.Required()
validation.Email()
validation.Min(10)
validation.Max(100)
validation.String()
validation.Int()
validation.Numeric()
validation.Boolean()
validation.Array()
validation.Between(1, 150)

// Database rules
validation.Unique("users", "email")
validation.UniqueExcept("users", "email", 1)
validation.Exists("users", "id")

// List rules
validation.In("admin", "user", "guest")
validation.NotIn("admin", "root")

// Numeric rules
validation.Decimal(2)        // decimal:2
validation.Decimal(2, 4)     // decimal:2,4

// Comparison rules
validation.Gt("field")
validation.Gte("field")
validation.Lt("field")
validation.Lte("field")

// File rules
validation.File()
validation.Extensions("jpg", "png")
validation.Mimes("jpg", "bmp")
validation.Mimetypes("image/*", "video/*")
validation.Dimensions(100, 200) // min_width=100,min_height=200

// Size rule
validation.Size(10)

// Other rules
validation.Confirmed()
validation.Date()
validation.URL()
validation.Alpha()
validation.AlphaNum()
validation.Phone()
validation.Password()
validation.UUID()
```

### Error Handling
#### Get All Errors

```go
errors := validator.Errors()
// Returns: map[string][]string{"field": ["error messages"]}
```

#### Get Specific Field Errors

```go
fieldErrors := validator.GetFieldErrors("email")
// Returns: []string{"The email must be a valid email address."}
``` 

### Check if Field Has Errors

```go
if validator.HasError("email") {
    // Handle email validation error
}
```

### Custom Error Messages
```go
validator := validation.NewValidator()
validator.SetMessage("required", "This field is required")
validator.SetMessage("email", "Please provide a valid email address")

// Or with field-specific messages
validator.WithMessages(map[string]string{
    "email.required": "Email address is required",
    "email.email":    "Please enter a valid email",
})
```

```text

This updated README now includes:
1. Documentation for all new rules (decimal, dimensions, gt, gte, lt, lte, file, extensions, mimes, mimetypes, size)
2. Usage examples for each new rule type
3. Helper functions for all new rules
4. Clear explanations of how each rule works
5. The UploadedFile struct documentation
6. Examples showing proper usage patterns
```





















