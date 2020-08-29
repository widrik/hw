/*
* Code generated automatically
* THIS FILE SHOULD NOT BE EDITED.
 */
package models

import (
	"fmt"
	"regexp"
)

type ValidationError struct {
	Field string
	Err   string
}

func (user User) Validate() ([]ValidationError, error) {
	var errs []ValidationError

	if len(user.ID) != 36 {
		errs = append(errs, ValidationError{
			Field: "ID",

			Err: "field ID must be length is 36",
		})
	}

	if user.Age < 18 {
		errs = append(errs, ValidationError{
			Field: "Age",

			Err: "field Age must be min is 18",
		})
	}

	if user.Age > 50 {
		errs = append(errs, ValidationError{
			Field: "Age",

			Err: "field Age must be max is 50",
		})
	}

	{
		matched, err := regexp.MatchString("^\\w+@\\w+\\.\\w+$", user.Email)
		if err != nil {
			return errs, err
		}
		if matched == false {
			errs = append(errs, ValidationError{
				Field: "Email",

				Err: "field Email must be in regexp ^\\w+@\\w+\\.\\w+$",
			})
		}
	}

	if !(user.Role == "admin" || user.Role == "stuff") {
		errs = append(errs, ValidationError{
			Field: "Role",

			Err: "field Role must be in range [admin stuff]",
		})
	}

	for i := range user.Phones {

		if len(user.Phones[i]) != 11 {
			errs = append(errs, ValidationError{
				Field: "Phones",

				Err: fmt.Sprintf("field Phones must be length is 11 in index %d", i),
			})
		}

	}

	return errs, nil
}

func (app App) Validate() ([]ValidationError, error) {
	var errs []ValidationError

	if len(app.Version) != 5 {
		errs = append(errs, ValidationError{
			Field: "Version",

			Err: "field Version must be length is 5",
		})
	}

	return errs, nil
}

func (response Response) Validate() ([]ValidationError, error) {
	var errs []ValidationError

	if !(response.Code == 200 || response.Code == 404 || response.Code == 500) {
		errs = append(errs, ValidationError{
			Field: "Code",

			Err: "field Code must be in range [200 404 500]",
		})
	}

	return errs, nil
}
