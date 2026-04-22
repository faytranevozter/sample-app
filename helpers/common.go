package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func Dump(values ...any) {
	for _, v := range values {
		fmt.Println(ToJSON(v, "\t"))
	}
}

func ToJSON(i any, indent string) string {
	s, _ := json.MarshalIndent(i, "", indent)
	return string(s)
}

func IsValidEmail(email string) bool {
	// Define the stricter email regex pattern (based on RFC 5322)
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

	// Compile the regex
	return regexp.MustCompile(emailRegex).MatchString(email)
}
