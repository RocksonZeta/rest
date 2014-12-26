package rest

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var EmailReg = regexp.MustCompile("(?i)^((([a-z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+(\\.([a-z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(\\\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-z]|\\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\\d|-|\\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\\d|-|\\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))$")

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
	this.Req.ParamErrors = append(this.Req.ParamErrors, Error)
}

type FieldValidator struct {
	Validator
	Value string
}

func (this *FieldValidator) Optional() *FieldValidator {
	if !this.Exists {
		this.GoOn = false
	}
	return this
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
		this.FireError(this.Key+"'length should between %d and %d.", tip, min, max)
	}
	return this
}
func (this *FieldValidator) ByteLen(min, max int, tip ...string) *FieldValidator {
	bl := len([]byte(this.Value))
	if this.GoOn && (bl < min || bl > max) {
		this.FireError("the byte length of the "+this.Key+" should between %d and %d.", tip, min, max)
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

func (this *FieldValidator) Eq(value string, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Value == value) {
		this.FireError(this.Key+" is should equal %s.", tip, value)
	}
	return this
}
func (this *FieldValidator) Neq(value string, tip ...string) *FieldValidator {
	if this.GoOn && this.Value == value {
		this.FireError(this.Key+" is should not equal %s.", tip, value)
	}
	return this
}

func (this *FieldValidator) Gt(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() > l) {
		this.FireError(fmt.Sprintf("%s must great than %d.", this.Key, l), tip)
	}
	return this
}
func (this *FieldValidator) Ge(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() >= l) {
		this.FireError(fmt.Sprintf("%s must great than or equal %d.", this.Key, l), tip)
	}
	return this
}
func (this *FieldValidator) Lt(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() < l) {
		this.FireError(fmt.Sprintf("%s must less than %d.", this.Key, l), tip)
	}
	return this
}
func (this *FieldValidator) Le(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() <= l) {
		this.FireError(fmt.Sprintf("%s must less than or equal %d.", this.Key, l), tip)
	}
	return this
}
func (this *FieldValidator) Contains(s string, tip ...string) *FieldValidator {
	if this.GoOn && !(strings.Contains(this.Value, s)) {
		this.FireError(fmt.Sprintf("%s shoud contains %s.", this.Key, s), tip)
	}
	return this
}
func (this *FieldValidator) NotContains(s string, tip ...string) *FieldValidator {
	if this.GoOn && (strings.Contains(this.Value, s)) {
		this.FireError(fmt.Sprintf("%s shoud not contains %s.", this.Key, s), tip)
	}
	return this
}
func (this *FieldValidator) IsEmail(tip ...string) *FieldValidator {
	if this.GoOn && !(EmailReg.MatchString(this.Value)) {
		this.FireError(this.Key+" is not a email format.", tip)
	}
	return this
}
func (this *FieldValidator) IsUrl(tip ...string) *FieldValidator {
	b := true
	protocols := [3]string{"http", "https", "ftp"}
	if 2083 <= len(this.Value) {
		b = false
	} else {
		u, e := url.Parse(this.Value)
		if nil != e || 0 == len(u.Scheme) || 0 == len(u.Host) {
			b = false
		}
		for _, p := range protocols {
			if p == u.Scheme {
				b = true
				break
			}
			b = false
		}
	}
	if this.GoOn && !b {
		this.FireError(this.Key+" is not a url format.", tip)
	}
	return this
}
func (this *FieldValidator) IsIp(tip ...string) *FieldValidator {
	if this.GoOn && !(Ipv4Reg.MatchString(this.Value) || Ipv6Reg.MatchString(this.Value)) {
		this.FireError(this.Key+" is not a ip format.", tip)
	}
	return this
}
func (this *FieldValidator) IsIpv4(tip ...string) *FieldValidator {
	if this.GoOn && !Ipv4Reg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a ipv4 format.", tip)
	}
	return this
}
func (this *FieldValidator) IsIpv6(tip ...string) *FieldValidator {
	if this.GoOn && !Ipv6Reg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a ipv6 format.", tip)
	}
	return this
}
func (this *FieldValidator) IsAlpha(tip ...string) *FieldValidator {
	if this.GoOn && !AlphaReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a alpha format.", tip)
	}
	return this
}
func (this *FieldValidator) IsNumeric(tip ...string) *FieldValidator {
	if this.GoOn && !NumericReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a numeric format.", tip)
	}
	return this
}
func (this *FieldValidator) IsAlphaNumeric(tip ...string) *FieldValidator {
	if this.GoOn && !AlphaNumericReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a alpha numeric format.", tip)
	}
	return this
}
func (this *FieldValidator) IsBase64(tip ...string) *FieldValidator {
	if this.GoOn && !Base64Reg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a base64 format.", tip)
	}
	return this
}
func (this *FieldValidator) IsHexadecimal(tip ...string) *FieldValidator {
	if this.GoOn && !HexadecimalReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a haxadecimal format.", tip)
	}
	return this
}
func (this *FieldValidator) IsLowercase(tip ...string) *FieldValidator {
	if this.GoOn && this.Value != strings.ToLower(this.Value) {
		this.FireError(this.Key+" is not lowercase.", tip)
	}
	return this
}
func (this *FieldValidator) IsUppercase(tip ...string) *FieldValidator {
	if this.GoOn && this.Value != strings.ToUpper(this.Value) {
		this.FireError(this.Key+" is not uppercase.", tip)
	}
	return this
}
func (this *FieldValidator) IsUUID(tip ...string) *FieldValidator {
	if this.GoOn && !UUIDReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a UUID format.", tip)
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
func (this *FieldValidator) Md5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(this.Value)))
}
func (this *FieldValidator) Sha1() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(this.Value)))
}
