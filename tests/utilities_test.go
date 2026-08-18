package tests

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/mwangaben/validation"
	. "github.com/onsi/gomega"
)

// ============ NEW RULES TESTS ============

func TestDecimalValidation(t *testing.T) {
	runTest(t, func() {
		validator := validation.NewValidator()
		Expect(validator).ToNot(BeNil())

		t.Run("should validate number with exactly 2 decimal places", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"price": 9.99,
				}
				rules := map[string][]string{
					"price": {"decimal:2"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject number with wrong decimal places", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"price": 9.999,
				}
				rules := map[string][]string{
					"price": {"decimal:2"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should validate string with correct decimal places", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"price": "99.99",
				}
				rules := map[string][]string{
					"price": {"decimal:2"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate decimal with range", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"price": 10.123,
				}
				rules := map[string][]string{
					"price": {"decimal:2,4"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject decimal outside range", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"price": 10.1,
				}
				rules := map[string][]string{
					"price": {"decimal:2,4"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should reject non-numeric value", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"price": "abc",
				}
				rules := map[string][]string{
					"price": {"decimal:2"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})
	})
}

func TestDimensionsValidation(t *testing.T) {
	runTest(t, func() {
		validator := validation.NewValidator()
		Expect(validator).ToNot(BeNil())

		// Create a test image file
		testImagePath := createTestImage(t, 200, 100)
		defer os.Remove(testImagePath)

		t.Run("should validate image with correct dimensions", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name: "test.png",
					Path: testImagePath,
					Size: 1024,
				}
				data := map[string]interface{}{
					"avatar": file,
				}
				rules := map[string][]string{
					"avatar": {"dimensions:min_width=100,min_height=50"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject image with insufficient dimensions", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name: "test.png",
					Path: testImagePath,
					Size: 1024,
				}
				data := map[string]interface{}{
					"avatar": file,
				}
				rules := map[string][]string{
					"avatar": {"dimensions:min_width=300,min_height=200"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should reject non-file value for dimensions", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"avatar": "not_a_file",
				}
				rules := map[string][]string{
					"avatar": {"dimensions:min_width=100,min_height=100"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})
	})
}

func TestComparisonRules(t *testing.T) {
	runTest(t, func() {
		validator := validation.NewValidator()
		Expect(validator).ToNot(BeNil())

		t.Run("should validate gt with field comparison", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"min_age": 18,
					"age":     25,
				}
				rules := map[string][]string{
					"age": {"gt:min_age"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject when gt condition not met", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"min_age": 18,
					"age":     16,
				}
				rules := map[string][]string{
					"age": {"gt:min_age"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should validate gt with direct value", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"age": 25,
				}
				rules := map[string][]string{
					"age": {"gt:18"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate gte with equal values", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"min_age": 18,
					"age":     18,
				}
				rules := map[string][]string{
					"age": {"gte:min_age"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate lt with field comparison", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"max_age": 65,
					"age":     30,
				}
				rules := map[string][]string{
					"age": {"lt:max_age"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate lte with equal values", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"max_age": 65,
					"age":     65,
				}
				rules := map[string][]string{
					"age": {"lte:max_age"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate string length comparison", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"field1": "short",
					"field2": "much longer string",
				}
				rules := map[string][]string{
					"field2": {"gt:field1"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})
	})
}

func TestFileValidation(t *testing.T) {
	runTest(t, func() {
		validator := validation.NewValidator()
		Expect(validator).ToNot(BeNil())

		t.Run("should validate file type", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "test.jpg",
					Size:        1024,
					ContentType: "image/jpeg",
				}
				data := map[string]interface{}{
					"photo": file,
				}
				rules := map[string][]string{
					"photo": {"file"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject non-file for file rule", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"photo": "not_a_file",
				}
				rules := map[string][]string{
					"photo": {"file"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should validate file extension", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "test.jpg",
					Size:        1024,
					ContentType: "image/jpeg",
				}
				data := map[string]interface{}{
					"photo": file,
				}
				rules := map[string][]string{
					"photo": {"extensions:jpg,png"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject wrong extension", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "test.gif",
					Size:        1024,
					ContentType: "image/gif",
				}
				data := map[string]interface{}{
					"photo": file,
				}
				rules := map[string][]string{
					"photo": {"extensions:jpg,png"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should validate MIME type by extension", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "test.jpg",
					Size:        1024,
					ContentType: "image/jpeg",
				}
				data := map[string]interface{}{
					"photo": file,
				}
				rules := map[string][]string{
					"photo": {"mimes:jpg,bmp,png"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate MIME type directly", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "video.mp4",
					Size:        1024,
					ContentType: "video/mp4",
				}
				data := map[string]interface{}{
					"video": file,
				}
				rules := map[string][]string{
					"video": {"mimetypes:video/avi,video/mpeg,video/mp4"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate MIME type with wildcard", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "image.png",
					Size:        1024,
					ContentType: "image/png",
				}
				data := map[string]interface{}{
					"media": file,
				}
				rules := map[string][]string{
					"media": {"mimetypes:image/*,video/*"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})
	})
}

func TestSizeValidation(t *testing.T) {
	runTest(t, func() {
		validator := validation.NewValidator()
		Expect(validator).ToNot(BeNil())

		t.Run("should validate string size", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"title": "Hello World", // 11 characters
				}
				rules := map[string][]string{
					"title": {"size:11"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate numeric size", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"seats": 10,
				}
				rules := map[string][]string{
					"seats": {"integer", "size:10"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate array size", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"tags": []string{"a", "b", "c", "d", "e"},
				}
				rules := map[string][]string{
					"tags": {"array", "size:5"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should validate file size in kilobytes", func(t *testing.T) {
			runTest(t, func() {
				file := &validation.UploadedFile{
					Name:        "test.jpg",
					Size:        512 * 1024, // 512 KB
					ContentType: "image/jpeg",
				}
				data := map[string]interface{}{
					"image": file,
				}
				rules := map[string][]string{
					"image": {"file", "size:512"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject wrong size", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"title": "Hello",
				}
				rules := map[string][]string{
					"title": {"size:10"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})
	})
}

func TestRuleHelperFunctions(t *testing.T) {
	runTest(t, func() {
		t.Run("should create decimal rule", func(t *testing.T) {
			rule := validation.Decimal(2)
			Expect(rule).To(Equal("decimal:2"))
		})

		t.Run("should create decimal rule with range", func(t *testing.T) {
			rule := validation.Decimal(2, 4)
			Expect(rule).To(Equal("decimal:2,4"))
		})

		t.Run("should create dimensions rule", func(t *testing.T) {
			rule := validation.Dimensions(100, 200)
			Expect(rule).To(Equal("dimensions:min_width=100,min_height=200"))
		})

		t.Run("should create comparison rules", func(t *testing.T) {
			Expect(validation.Gt("field")).To(Equal("gt:field"))
			Expect(validation.Gte("field")).To(Equal("gte:field"))
			Expect(validation.Lt("field")).To(Equal("lt:field"))
			Expect(validation.Lte("field")).To(Equal("lte:field"))
		})

		t.Run("should create file rules", func(t *testing.T) {
			Expect(validation.File()).To(Equal("file"))
			Expect(validation.Extensions("jpg", "png")).To(Equal("extensions:jpg,png"))
			Expect(validation.Mimes("jpg", "bmp", "png")).To(Equal("mimes:jpg,bmp,png"))
			Expect(validation.Mimetypes("text/plain", "image/*")).To(Equal("mimetypes:text/plain,image/*"))
		})

		t.Run("should create size rule", func(t *testing.T) {
			Expect(validation.Size(10)).To(Equal("size:10"))
		})
	})
}

// Helper function to create test image
func createTestImage(t *testing.T, width, height int) string {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with some color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	file, err := os.CreateTemp("", "test_image_*.png")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}

	return file.Name()
}
