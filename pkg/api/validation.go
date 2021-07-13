package api

import (
	"github.com/charlie1404/vqueue/pkg/storage"
	"github.com/go-playground/validator/v10"
)

type ApiValidator struct {
	*validator.Validate
}

func (v *ApiValidator) validateCreateQueueInput(ip *storage.CreateQueueInput) error {
	vv := validator.New()
	err := vv.Struct(ip)

	if err != nil {
		return err
		// TODO: make human friendly later

		// validationErrors := err.(validator.ValidationErrors)
		// log.Printf("%+v \n=========>\n", validationErrors)

		// for _, err := range err.(validator.ValidationErrors) {
		// 	fmt.Println("Namespace => ", err.Namespace())
		// 	fmt.Println("Field => ", err.Field())
		// 	fmt.Println("StructNamespace => ", err.StructNamespace())
		// 	fmt.Println("StructField => ", err.StructField())
		// 	fmt.Println("Tag => ", err.Tag())
		// 	fmt.Println("ActualTag => ", err.ActualTag())
		// 	fmt.Println("Kind => ", err.Kind())
		// 	fmt.Println("Type => ", err.Type())
		// 	fmt.Println("Value => ", err.Value())
		// 	fmt.Println("Param => ", err.Param())
		// }
	}

	return nil
}

func newValidator() ApiValidator {
	return ApiValidator{validator.New()}
}

func (v *ApiValidator) registerCustomValidations() {
	// look later
}
