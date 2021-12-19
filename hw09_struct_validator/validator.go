package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedType = errors.New("unsupported type. support only structures")
	ErrFormatRule      = errors.New("unsupported format rule")
	ErrUnknownRule     = errors.New("unknown rule")
	ErrValueIsInvalid  = errors.New("value does not satisfy condition")
	ErrNotNumber       = errors.New("not number")
	ErrNotValidRegexp  = errors.New("not valid regexp")

	supportedTypes = map[reflect.Kind]bool{
		reflect.String: true,
		reflect.Int:    true,
		reflect.Slice:  true,
		reflect.Array:  true,
	}

	ruleFunctions = map[string]func(value, ruleCondition string) (bool, error){
		"len": func(value, ruleCondition string) (bool, error) {
			lenCondition, err := strconv.Atoi(ruleCondition)
			if err != nil {
				return false, ErrNotNumber
			}
			return len([]rune(value)) == lenCondition, nil
		},
		"regexp": func(value, ruleCondition string) (bool, error) {
			reg, err := regexp.Compile(ruleCondition)
			if err != nil {
				return false, ErrNotValidRegexp
			}
			return reg.MatchString(value), nil
		},
		"in": func(value, ruleCondition string) (bool, error) { //nolint:unparam
			for _, str := range strings.Split(ruleCondition, ",") {
				if str == value {
					return true, nil
				}
			}
			return false, nil
		},
		"min": func(value, ruleCondition string) (bool, error) {
			minCondition, err := strconv.Atoi(ruleCondition)
			if err != nil {
				return false, ErrNotNumber
			}

			val, err := strconv.Atoi(value)
			if err != nil {
				return false, ErrNotNumber
			}

			return val >= minCondition, nil
		},
		"max": func(value, ruleCondition string) (bool, error) {
			maxCondition, err := strconv.Atoi(ruleCondition)
			if err != nil {
				return false, ErrNotNumber
			}

			val, err := strconv.Atoi(value)
			if err != nil {
				return false, ErrNotNumber
			}

			return maxCondition >= val, nil
		},
	}
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for _, vError := range v {
		sb.WriteString("Field: ")
		sb.WriteString(vError.Field)
		sb.WriteString(", Error: ")
		sb.WriteString(vError.Err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	var vErr ValidationErrors

	countField := rv.NumField()
	for i := 0; i < countField; i++ {
		fValue := rv.Field(i)

		kind := fValue.Kind()
		if _, ok := supportedTypes[kind]; !ok {
			continue
		}

		field := rv.Type().Field(i)
		tag := field.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}

		var values []interface{}

		if kind == reflect.Slice || kind == reflect.Array {
			if fValue.Len() == 0 {
				continue
			}
			first := fValue.Index(0).Interface()
			switch first.(type) {
			case int, string:
				for j := 0; j < fValue.Len(); j++ {
					values = append(values, fValue.Index(j).Interface())
				}
			default:
				continue
			}
		} else {
			values = append(values, fValue.Interface())
		}
		for _, value := range values {
			vErr = append(vErr, validateValue(value, field.Name, tag)...)
		}
	}

	if len(vErr) == 0 {
		return nil
	}

	return vErr
}

func validateValue(value interface{}, fieldName, someRuleStr string) ValidationErrors {
	var vErr ValidationErrors

	rulesStr := strings.Split(someRuleStr, "|")
	for _, ruleStr := range rulesStr {
		rule := strings.Split(ruleStr, ":")
		if len(rule) != 2 {
			vErr = append(vErr, ValidationError{Field: fieldName, Err: fmt.Errorf("%w, %s", ErrFormatRule, ruleStr)})
			continue
		}

		ruleName := rule[0]
		ruleCondition := rule[1]

		ruleFunction, ok := ruleFunctions[ruleName]
		if !ok {
			vErr = append(vErr, ValidationError{Field: fieldName, Err: fmt.Errorf("%w, %s", ErrUnknownRule, ruleName)})
			continue
		}

		ok, err := ruleFunction(fmt.Sprintf("%v", value), ruleCondition)
		if err != nil {
			vErr = append(vErr, ValidationError{Field: fieldName, Err: err})
			continue
		}

		if !ok {
			vErr = append(vErr, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("%w, %s, value:%s", ErrValueIsInvalid, ruleStr, fmt.Sprintf("%v", value)),
			})
		}
	}

	if len(vErr) == 0 {
		return nil
	}

	return vErr
}
