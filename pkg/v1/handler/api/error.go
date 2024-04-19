package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// TODO show multiple errors
func errorInvalidArguments(err error) gin.H {
	errors := strings.Split(err.Error(), "Key:")
	var parsedErrors []string
	for _, e := range errors {
		// check error is not empty
		if e != "" {
			var parsedTag *string
			var parsedField *string
			var parsedError string
			// remove ' from string regexp
			re := regexp.MustCompile(`(\w+)'`)
			// trim error
			trimmedError := strings.TrimLeft(e, " ")
			//* get fields
			fields := strings.Split(trimmedError, " ")

			//* get tag
			tag := fields[len(fields)-2]
			// remove ' from tag
			tagMatch := re.FindStringSubmatch(tag)
			if len(tagMatch) > 1 {
				parsedTag = &tagMatch[1]
			}

			//* get field
			field := fields[4]
			fieldMatch := re.FindStringSubmatch(field)
			if len(fieldMatch) > 1 {
				firstLetter := strings.ToLower(fieldMatch[1][:1])
				otherLetters := fieldMatch[1][1:]
				convertedField := strings.Join([]string{firstLetter, otherLetters}, "")
				parsedField = &convertedField
			}
			if parsedField != nil && parsedTag != nil {
				switch *parsedTag {
				case "required":
					{
						parsedError = fmt.Sprintf("Field %v is required", *parsedField)
					}
				case "email":
					{
						parsedError = fmt.Sprintf("Field %v should be email", *parsedField)
					}
				case "min":
					{
						parsedError = fmt.Sprintf("Field %v should be more", *parsedField)
					}
				}
			}
			parsedErrors = append(parsedErrors, parsedError)
		}
	}
	if parsedErrors[0] != "" {
		return gin.H{"error": parsedErrors[0]}
	}
	return gin.H{"error": err.Error()}
}

func errorCode(code *codes.Code) int {
	// var httpCode int
	httpCode := http.StatusInternalServerError
	if code != nil {
		switch *code {
		case codes.AlreadyExists:
			{
				httpCode = http.StatusBadRequest
			}
		case codes.Unimplemented:
			{
				httpCode = http.StatusInternalServerError
			}
		case codes.Internal:
			{
				httpCode = http.StatusInternalServerError
			}
		case codes.NotFound:
			{
				httpCode = http.StatusNotFound
			}
		case codes.Unauthenticated:
			{
				httpCode = http.StatusUnauthorized
			}
		case codes.PermissionDenied:
			{
				httpCode = http.StatusUnauthorized
			}
		}
	}
	return httpCode
}
