package validators

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const validatorTag = "validate"

type Validator interface {
	Validate(v interface{}) error
}

type ValidatorFunc func(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error

type ValidatorLibrary map[string]ValidatorFunc

func (a ValidatorLibrary) Validate(value interface{}) error {
	t := reflect.TypeOf(value)
	// 去指针化
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(value)
	// 去指针化
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tags, exist := t.Field(i).Tag.Lookup(validatorTag)
		// 如果没有验证的字段,则直接跳过
		if !exist {
			continue
		}
		// 根据分号分隔
		// required(m=姓名不能为空);max_length(m=姓名长度不能大于10,value=10)
		// required(m=姓名不能为空)  max_length(m=姓名长度不能大于10,value=10)
		validatorFuncNames := strings.Split(tags, ";")

		for _, validatorFuncName := range validatorFuncNames {
			// 正则获取取校验函数名和参数
			matchList := validateParamRegexp.FindAllStringSubmatch(validatorFuncName, -1)
			if len(matchList) > 0 {
				// 如果匹配到了,则取第一个
				result := matchList[0]
				key := result[1]
				// 获取校验函数
				validator, exist := a[key]
				// 没有获取到
				// 再往下则没有意义
				// 返回错误信息
				if !exist {
					return fmt.Errorf("%s validator does not exist", result[1])
				}
				// 获取参数
				// m=姓名长度不能大于10,value=10
				paramsText := result[2]
				// 根据,分隔
				paramsList := strings.Split(paramsText, ",")
				// 判断分隔后的长度是否大于2, 可能出现情况
				// m=姓名长度不能大于10,value=10,,,,,,
				// 这种情况的value就等于10,,,,,,
				// 我们的参数只有两个, 直接取前面两个就OK了
				if len(paramsList) > 2 {
					paramsList[1] = strings.TrimPrefix(paramsText, paramsList[0]+",")
					paramsList = paramsList[:2]
				}
				var validateParam Param

				for i := 0; i < len(paramsList); i++ {
					param := paramsList[i]
					// 首位去空格化
					for strings.HasPrefix(param, " ") {
						param = strings.TrimPrefix(param, " ")
					}
					// 将参数进行key-value分隔
					item := strings.Split(param, "=")
					if len(item) != 2 {
						return fmt.Errorf("%s syntax error", param)
					}

					// 分别进行赋值
					k, v := item[0], item[1]
					switch k {
					case "m", "message":
						validateParam.Message = v
					case "v", "value":
						validateParam.Value = v
					default:
						return fmt.Errorf("unexpect param got %s", k)
					}
				}

				// to validate
				if err := validator(t, v, i, validateParam); err != nil {
					return NewValidationError(err, t.Field(i).Name, key)
				}
			}
		}
	}
	return nil
}

var validatorLibrary = ValidatorLibrary{}

// RegisterValidator add more validator in
func RegisterValidator(key string, v ValidatorFunc) error {
	if _, exist := validatorLibrary[key]; exist {
		return fmt.Errorf("key %s has already exist")
	}
	validatorLibrary[key] = v
	return nil
}

func isValid(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() != ""
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int16, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return !value.IsZero()
	case reflect.Bool:
		return value.Bool()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return value.Len() > 0
	case reflect.Struct:
		return !reflect.DeepEqual(value.Interface(), reflect.New(value.Type()).Elem().Interface())
	case reflect.Ptr, reflect.Interface:
		return !value.IsNil()
	}
	return true
}

// Required except not zero value
func Required(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if isValid(valueOf.Field(index)) {
		return nil
	}
	return errors.New(param.Message)
}

func MaxLength(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() == reflect.String {
		except, err := strconv.Atoi(param.Value)
		if err != nil {
			return err
		}
		if utf8.RuneCountInString(field.String()) >= except {
			return errors.New(param.Message)
		}
	}
	return nil
}

func MinLength(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() == reflect.String {
		except, err := strconv.Atoi(param.Value)
		if err != nil {
			return err
		}
		if utf8.RuneCountInString(field.String()) <= except {
			return errors.New(param.Message)
		}
	}
	return nil
}

func isGt(field reflect.Value, v string) (bool, error) {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Int() > value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Uint() > value, nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
		return field.Float() > value, nil
	}
	return false, unsupportedError
}

func GreaterThan(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isGt(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func isLt(field reflect.Value, v string) (bool, error) {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Int() < value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Uint() < value, nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
		return field.Float() < value, nil
	}
	return false, unsupportedError
}

func LowerThan(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isLt(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func isGe(field reflect.Value, v string) (bool, error) {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Int() >= value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Uint() >= value, nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
		return field.Float() >= value, nil
	}
	return false, nil
}

func GreaterEqual(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isGe(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func isLe(field reflect.Value, v string) (bool, error) {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Int() <= value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Uint() <= value, nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
		return field.Float() <= value, nil
	}
	return false, unsupportedError
}

func LowerEqual(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isLe(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func isEqual(field reflect.Value, v string) (bool, error) {
	switch field.Kind() {
	case reflect.String:
		return field.String() == v, nil
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Int() == value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Uint() == value, nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
		return field.Float() == value, nil
	case reflect.Bool:
		value, err := strconv.ParseBool(v)
		if err != nil {
			return false, err
		}
		return field.Bool() == value, nil
	}
	return false, unsupportedError
}

func Equal(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isEqual(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func EqualField(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	f := valueOf.FieldByName(param.Value)
	if !f.IsValid() {
		return fmt.Errorf("no field named %s", param.Value)
	}
	if !reflect.DeepEqual(valueOf.Field(index).Interface(), f.Interface()) {
		return errors.New(param.Message)
	}
	return nil
}

func Method(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Value == "" {
		return valueParamRequiredError
	}
	method := valueOf.MethodByName(param.Value)
	if !method.IsValid() {
		return fmt.Errorf("no method named %s", param.Value)
	}
	if m, ok := method.Interface().(func(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error); ok {
		return m(typeOf, valueOf, index, param)
	}
	return fmt.Errorf("method %s must be `func(typeOf reflect.Type, valueOf reflect.Value, index int, param validate.ValidateParam) error` type ", param.Value)
}

func Regexp(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)
	reg := regexp.MustCompile(param.Value)
	if field.Kind() != reflect.String {
		return fmt.Errorf("Regexp only support string type")
	}
	if !reg.MatchString(field.String()) {
		return errors.New(param.Message)
	}
	return nil
}

func Email(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() != reflect.String {
		return errors.New("Email only support string type")
	}
	if !emailRegex.MatchString(field.String()) {
		return errors.New(param.Message)
	}
	return nil
}

func UUID(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() != reflect.String {
		return errors.New("UUID only support string type")
	}
	if !uuidRegexp.MatchString(field.String()) {
		return errors.New(param.Message)
	}
	return nil
}

func Phone(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() != reflect.String {
		return errors.New("Phone only support string type")
	}
	if !phoneRegexp.MatchString(field.String()) {
		return errors.New(param.Message)
	}
	return nil
}

func isBetween(field reflect.Value, min, max string) (bool, error) {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		mi, err := strconv.ParseInt(min, 10, 64)
		if err != nil {
			return false, err
		}
		ma, err := strconv.ParseInt(max, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Int() >= mi && field.Int() <= ma, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mi, err := strconv.ParseUint(min, 10, 64)
		if err != nil {
			return false, err
		}
		ma, err := strconv.ParseUint(max, 10, 64)
		if err != nil {
			return false, err
		}
		return field.Uint() >= mi && field.Uint() <= ma, nil
	case reflect.Float32, reflect.Float64:
		mi, err := strconv.ParseFloat(min, 64)
		if err != nil {
			return false, err
		}
		ma, err := strconv.ParseFloat(max, 64)
		if err != nil {
			return false, err
		}
		return field.Float() >= mi && field.Float() <= ma, nil
	default:
		return false, nil
	}
}

func Between(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	value := strings.Split(param.Value, ",")
	if len(value) != 2 {
		return fmt.Errorf("%s syntax error", param.Value)
	}
	min := strings.ReplaceAll(value[0], " ", "")
	max := strings.ReplaceAll(value[1], " ", "")
	ok, err := isBetween(valueOf.Field(index), min, max)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func NotBetween(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	value := strings.Split(param.Value, ",")
	if len(value) != 2 {
		return fmt.Errorf("%s syntax error", param.Value)
	}
	min := strings.ReplaceAll(value[0], " ", "")
	max := strings.ReplaceAll(value[1], " ", "")
	ok, err := isBetween(valueOf.Field(index), min, max)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(param.Message)
	}
	return nil
}

func NotEqual(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isEqual(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(param.Message)
	}
	return nil
}

func NotEqualField(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	f := valueOf.FieldByName(param.Value)
	if !f.IsValid() {
		return fmt.Errorf("no field named %s", param.Value)
	}
	if reflect.DeepEqual(valueOf.Field(index).Interface(), f.Interface()) {
		return errors.New(param.Message)
	}
	return nil
}

func Url(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() != reflect.String {
		return errors.New("Url only support string type")
	}
	if !urlRegexp.MatchString(field.String()) {
		return errors.New(param.Message)
	}
	return nil
}

func isContains(field reflect.Value, value string) (bool, error) {
	switch field.Kind() {
	case reflect.String:
		values := strings.Split(value, ",")
		v := field.String()
		for _, str := range values {
			if str == v {
				return true, nil
			}
		}
		return false, nil
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		values := strings.Split(value, ",")
		v := field.Int()
		for _, item := range values {
			i, err := strconv.ParseInt(strings.ReplaceAll(item, " ", ""), 10, 64)
			if err != nil {
				return false, err
			}
			if v == i {
				return true, nil
			}
		}
		return false, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		values := strings.Split(value, ",")
		v := field.Uint()
		for _, item := range values {
			i, err := strconv.ParseUint(strings.ReplaceAll(item, " ", ""), 10, 64)
			if err != nil {
				return false, err
			}

			if v == i {
				return true, nil
			}
		}
		return false, nil
	case reflect.Float32, reflect.Float64:
		values := strings.Split(value, ",")
		v := field.Float()
		for _, item := range values {
			i, err := strconv.ParseFloat(strings.ReplaceAll(item, " ", ""), 64)
			if err != nil {
				return false, err
			}
			if v == i {
				return true, nil
			}
		}
		return false, nil
	}
	return false, unsupportedError
}

func Contains(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isContains(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(param.Message)
	}
	return nil
}

func NotContains(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	ok, err := isContains(valueOf.Field(index), param.Value)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(param.Message)
	}
	return nil
}

func DateFormat(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)
	if field.Kind() != reflect.String {
		return errors.New("DateFormat only support string type")
	}
	if _, err := time.Parse(param.Value, field.String()); err != nil {
		return errors.New(param.Message)
	}
	return nil
}

func GreatThanField(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)
	anotherField := valueOf.FieldByName(param.Value)
	if !anotherField.IsValid() {
		return errors.New("no field name " + param.Value)
	}
	if field.Kind() != anotherField.Kind() {
		return fmt.Errorf("%s and %s are different kind", typeOf.Field(index).Name, param.Value)
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if field.Int() > anotherField.Int() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() > anotherField.Uint() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Float32, reflect.Float64:
		if field.Float() > anotherField.Float() {
			return nil
		}
		return errors.New(param.Message)
	}
	return errors.New("unsupported type")
}

func GreatEqualField(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)
	anotherField := valueOf.FieldByName(param.Value)
	if !anotherField.IsValid() {
		return errors.New("no field name " + param.Value)
	}
	if field.Kind() != anotherField.Kind() {
		return fmt.Errorf("%s and %s are different kind", typeOf.Field(index).Name, param.Value)
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if field.Int() >= anotherField.Int() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() >= anotherField.Uint() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Float32, reflect.Float64:
		if field.Float() >= anotherField.Float() {
			return nil
		}
		return errors.New(param.Message)
	}
	return errors.New("unsupported type")
}

func LowerThanField(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)
	anotherField := valueOf.FieldByName(param.Value)
	if !anotherField.IsValid() {
		return errors.New("no field name " + param.Value)
	}
	if field.Kind() != anotherField.Kind() {
		return fmt.Errorf("%s and %s are different kind", typeOf.Field(index).Name, param.Value)
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if field.Int() < anotherField.Int() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() < anotherField.Uint() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Float32, reflect.Float64:
		if field.Float() < anotherField.Float() {
			return nil
		}
		return errors.New(param.Message)
	}
	return errors.New("unsupported type")
}

func LowerEqualField(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)
	anotherField := valueOf.FieldByName(param.Value)
	if !anotherField.IsValid() {
		return errors.New("no field name " + param.Value)
	}
	if field.Kind() != anotherField.Kind() {
		return fmt.Errorf("%s and %s are different kind", typeOf.Field(index).Name, param.Value)
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if field.Int() <= anotherField.Int() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() <= anotherField.Uint() {
			return nil
		}
		return errors.New(param.Message)
	case reflect.Float32, reflect.Float64:
		if field.Float() <= anotherField.Float() {
			return nil
		}
		return errors.New(param.Message)
	}
	return errors.New("unsupported type")
}

func Round(typeOf reflect.Type, valueOf reflect.Value, index int, param Param) error {
	if param.Message == "" {
		return messageParamRequiredError
	}
	if param.Value == "" {
		return valueParamRequiredError
	}
	field := valueOf.Field(index)

	v, err := strconv.Atoi(param.Value)
	if err != nil {
		return err
	}
	var text string
	switch field.Kind() {
	case reflect.Float32:
		text = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		text = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		return errors.New("round support float type only")
	}
	s := strings.Split(text, ".")
	if len(s) == 2 && len(s[1]) > v {
		return errors.New(param.Message)
	}
	return nil
}

func init() {
	RegisterValidator("required", Required)
	RegisterValidator("max_length", MaxLength)
	RegisterValidator("min_length", MinLength)
	RegisterValidator("gt", GreaterThan)
	RegisterValidator("lt", LowerThan)
	RegisterValidator("ge", GreaterEqual)
	RegisterValidator("le", LowerEqual)
	RegisterValidator("equal", Equal)
	RegisterValidator("eq", Equal)
	RegisterValidator("equal_field", EqualField)
	RegisterValidator("ef", EqualField)
	RegisterValidator("method", Method)
	RegisterValidator("regexp", Regexp)
	RegisterValidator("email", Email)
	RegisterValidator("uuid", UUID)
	RegisterValidator("phone", Phone)
	RegisterValidator("between", Between)
	RegisterValidator("not_between", NotBetween)
	RegisterValidator("not_equal", NotEqual)
	RegisterValidator("ne", NotEqual)
	RegisterValidator("not_equal_field", NotEqualField)
	RegisterValidator("nef", NotEqualField)
	RegisterValidator("url", Url)
	RegisterValidator("contains", Contains)
	RegisterValidator("not_contains", NotContains)
	RegisterValidator("date_format", DateFormat)
	RegisterValidator("great_than_field", GreatThanField)
	RegisterValidator("gtf", GreatThanField)
	RegisterValidator("lower_than_field", LowerThanField)
	RegisterValidator("ltf", LowerThanField)
	RegisterValidator("great_equal_field", GreatEqualField)
	RegisterValidator("gef", GreatEqualField)
	RegisterValidator("lower_equal_field", LowerEqualField)
	RegisterValidator("lef", LowerEqualField)
	RegisterValidator("round", Round)
}
