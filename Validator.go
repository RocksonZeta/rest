package rest

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//var EmailReg = regexp.MustCompile("(?i)^((([a-z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+(\\.([a-z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+)*)|((\x22)((((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(([\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(\\([\x01-\x09\x0b\x0c\x0d-\x7f]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))))*(((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(\x22)))@((([a-z]|\\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\\d|-|\\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\\d|-|\\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))$")

var CreditCardReg = regexp.MustCompile("^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$")

var Ipv4Reg = regexp.MustCompile("^(\\d?\\d?\\d)\\.(\\d?\\d?\\d)\\.(\\d?\\d?\\d)\\.(\\d?\\d?\\d)$")
var Ipv6Reg = regexp.MustCompile("^::|^::1|^([a-fA-F0-9]{1,4}::?){1,7}([a-fA-F0-9]{1,4})$")

var UUIDReg = regexp.MustCompile("^[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}$")

var AlphaReg = regexp.MustCompile("^[a-zA-Z]+$")
var AlphaNumericReg = regexp.MustCompile("^[a-zA-Z0-9]+$")
var NumericReg = regexp.MustCompile("^-?[0-9]+$")
var IntReg = regexp.MustCompile("^(?:-?(?:0|[1-9][0-9]*))$")
var FloatReg = regexp.MustCompile("^(?:-?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$")
var HexadecimalReg = regexp.MustCompile("^[0-9a-fA-F]+$")

var AsciiReg = regexp.MustCompile("^[\x00-\x7F]+$")
var MultibyteReg = regexp.MustCompile("[^\x00-\x7F]")
var FullWidthReg = regexp.MustCompile("[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]")
var HalfWidthReg = regexp.MustCompile("[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]")

var Base64Reg = regexp.MustCompile("^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$")

type Validator struct {
	Key    string
	GoOn   bool
	Exists bool
	Error  string
	Req    *Request
}

func (this *Validator) HasError() bool {
	return 0 == len(this.Error)
}

func (this *Validator) FireError(defaultTip string, tip []string, args ...interface{}) {
	this.GoOn = false
	var Error string
	if 0 == len(tip) {
		if 0 == len(args) {
			Error = defaultTip
		} else {
			Error = fmt.Sprintf(defaultTip, args...)
		}
	} else {
		if 0 == len(args) {
			Error = tip[0]
		} else {
			Error = fmt.Sprintf(tip[0], args...)
		}
	}
	this.Req.Context.ParamErrors = append(this.Req.Context.ParamErrors, Error)
}

func (this *Validator) Optional() *Validator {
	if !this.Exists {
		this.GoOn = false
	}
	return this
}

type FieldValidator struct {
	Validator
	Value string
}

func (this *FieldValidator) Empty() *FieldValidator {
	if this.GoOn {
		if 0 == len(this.Value) {
			this.GoOn = false
		}
	}
	return this
}
func (this *FieldValidator) NotEmpty(tip ...string) *FieldValidator {
	if this.GoOn && 0 == len(this.Value) {
		this.FireError(this.Key+" can not be empty.", tip)
	}
	return this
}

func (this *FieldValidator) Match(reg string, tip ...string) *FieldValidator {
	if this.GoOn && !regexp.MustCompile(reg).MatchString(this.Value) {
		this.FireError(this.Key+" is bad format.", tip)
	}
	return this
}
func (this *FieldValidator) IsInt(tip ...string) *FieldValidator {
	if this.GoOn && !IntReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not integer format.", tip)
	}
	return this
}

func (this *FieldValidator) IsFloat(tip ...string) *FieldValidator {
	if this.GoOn && !FloatReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not float format.", tip)
	}
	return this
}

//len(value) should in [min ,max] , min<=len(value)<=max
func (this *FieldValidator) Len(min, max int, tip ...string) *FieldValidator {
	var l = len(this.Value)
	if this.GoOn && (l < min || l > max) {
		this.FireError(this.Key+" should between %d and %d.", tip, min, max)
	}
	return this
}
func (this *FieldValidator) In(values []string, tip ...string) *FieldValidator {
	exists := false
	for _, v := range values {
		if v == this.Value {
			exists = true
			break
		}
	}
	if this.GoOn && !exists {
		this.FireError(this.Key+" must in [%s].", tip, strings.Join(values, ","))
	}
	return this
}

func (this *FieldValidator) eq(value string, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Value == value) {
		this.FireError(this.Key+" is should equal %s.", tip, value)
	}
	return this
}
func (this *FieldValidator) neq(value string, tip ...string) *FieldValidator {
	if this.GoOn && this.Value == value {
		this.FireError(this.Key+" is should not equal %s.", tip, value)
	}
	return this
}

func (this *FieldValidator) gt(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() > l) {
		this.FireError(fmt.Sprintf("%s must great than %d. ", this.Key, l), tip)
	}
	return this
}
func (this *FieldValidator) ge(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() >= l) {
		this.FireError(fmt.Sprintf("%s must great than or equal %d. ", this.Key, l), tip)
	}
	return this
}

//Sanitizers

//to string
func (this *FieldValidator) String() string {
	return this.Value
}

func (this *FieldValidator) Int(tip ...string) int {
	r, e := strconv.Atoi(this.Value)
	if nil != e {
		this.FireError(this.Key+" is not integer.", tip)
	}
	return r
}
func (this *FieldValidator) Int64(tip ...string) int64 {
	r, e := strconv.ParseInt(this.Value, 10, 64)
	if nil != e {
		this.FireError(this.Key+" is not integer.", tip)
	}
	return r
}

func (this *FieldValidator) Float(tip ...string) float64 {
	r, e := strconv.ParseFloat(this.Value, 64)
	if nil != e {
		this.FireError(this.Key+" is not float.", tip)
	}

	return r
}

func (this *FieldValidator) Bool() bool {
	if 0 == len(this.Value) {
		return false
	}
	low := strings.ToLower(this.Value)
	if "false" == low {
		return false
	}
	if "0" == low {
		return false
	}
	return true
}
func (this *FieldValidator) Date(format string, tip ...string) time.Time {
	r, e := time.Parse(format, this.Value)
	if nil != e {
		panic(e)
	}
	return r
}
