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

//Request parameter validator
//All parameters should validate first
type FieldValidator struct {
	Validator
	Value string
}

//check the field exists, if not exit ,then check over with no errors.
func (this *FieldValidator) Optional() *FieldValidator {
	if !this.Exists || 0 == len(this.Value) {
		this.GoOn = false
	}
	return this
}

//check the field is empty, if not exit ,then check over with no errors.
func (this *FieldValidator) Empty() *FieldValidator {
	if this.GoOn {
		if 0 == len(this.Value) {
			this.GoOn = false
		}
	}
	return this
}

//check the field is not empty,
func (this *FieldValidator) NotEmpty(tip ...string) *FieldValidator {
	if this.GoOn && 0 == len(this.Value) {
		this.FireError(this.Key+" can not be empty.", tip)
	}
	return this
}

//check the field if match the reg expression
func (this *FieldValidator) Match(reg string, tip ...string) *FieldValidator {
	if this.GoOn && !regexp.MustCompile(reg).MatchString(this.Value) {
		this.FireError(this.Key+" is bad format.", tip)
	}
	return this
}

//check the field if it is a integer format
func (this *FieldValidator) IsInt(tip ...string) *FieldValidator {
	if this.GoOn && !IntReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not integer format.", tip)
	}
	return this
}

//check the field if it is a float format
func (this *FieldValidator) IsFloat(tip ...string) *FieldValidator {
	if this.GoOn && !FloatReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not float format.", tip)
	}
	return this
}

//check the field's length
//len(value) should in [min ,max] , min<=len(value)<=max
func (this *FieldValidator) Len(min, max int, tip ...string) *FieldValidator {
	var l = len(this.Value)
	if this.GoOn && (l < min || l > max) {
		this.FireError(fmt.Sprintf("%s'length should between %d and %d.", this.Key, min, max), tip)
	}
	return this
}

//check the field's byte length in utf8
func (this *FieldValidator) ByteLen(min, max int, tip ...string) *FieldValidator {
	bl := len([]byte(this.Value))
	if this.GoOn && (bl < min || bl > max) {
		this.FireError(fmt.Sprintf("the byte length of the %s should between %d and %d.", this.Key, min, max), tip)
	}
	return this
}

//check the field if it is in specified values
func (this *FieldValidator) In(values []string, tip ...string) *FieldValidator {
	exists := false
	for _, v := range values {
		if v == this.Value {
			exists = true
			break
		}
	}
	if this.GoOn && !exists {
		this.FireError(fmt.Sprintf("%s must in [%s].", this.Key, strings.Join(values, ",")), tip)
	}
	return this
}

//check the field if it equal the value
func (this *FieldValidator) Eq(value string, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Value == value) {
		this.FireError(fmt.Sprintf("%s is should equal %s.", this.Key, value), tip)
	}
	return this
}

//check the field if it not equal the value
func (this *FieldValidator) Neq(value string, tip ...string) *FieldValidator {
	if this.GoOn && this.Value == value {
		this.FireError(fmt.Sprintf("%s is should not equal %s.", this.Key, value), tip)
	}
	return this
}

//check the field if it great than the value
func (this *FieldValidator) Gt(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() > l) {
		this.FireError(fmt.Sprintf("%s must great than %d.", this.Key, l), tip)
	}
	return this
}

//check the field if it great than or equal the the value
func (this *FieldValidator) Ge(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() >= l) {
		this.FireError(fmt.Sprintf("%s must great than or equal %d.", this.Key, l), tip)
	}
	return this
}

//check the field if it less than the value
func (this *FieldValidator) Lt(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() < l) {
		this.FireError(fmt.Sprintf("%s must less than %d.", this.Key, l), tip)
	}
	return this
}

//check the field if it less than or equal the the value
func (this *FieldValidator) Le(l int, tip ...string) *FieldValidator {
	if this.GoOn && !(this.Int() <= l) {
		this.FireError(fmt.Sprintf("%s must less than or equal %d.", this.Key, l), tip)
	}
	return this
}

//check the field if it contains the string
func (this *FieldValidator) Contains(s string, tip ...string) *FieldValidator {
	if this.GoOn && !(strings.Contains(this.Value, s)) {
		this.FireError(fmt.Sprintf("%s shoud contains %s.", this.Key, s), tip)
	}
	return this
}

//check the field if it contains the string
func (this *FieldValidator) NotContains(s string, tip ...string) *FieldValidator {
	if this.GoOn && (strings.Contains(this.Value, s)) {
		this.FireError(fmt.Sprintf("%s shoud not contains %s.", this.Key, s), tip)
	}
	return this
}

//check the field if it is a email format
func (this *FieldValidator) IsEmail(tip ...string) *FieldValidator {
	if this.GoOn && !(EmailReg.MatchString(this.Value)) {
		this.FireError(this.Key+" is not a email format.", tip)
	}
	return this
}

//check the field if it is a url format
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

//check the field if it is a ip(ipv4 or ipv6) format
func (this *FieldValidator) IsIp(tip ...string) *FieldValidator {
	if this.GoOn && !(Ipv4Reg.MatchString(this.Value) || Ipv6Reg.MatchString(this.Value)) {
		this.FireError(this.Key+" is not a ip format.", tip)
	}
	return this
}

//check the field if it is a ipv4 format
func (this *FieldValidator) IsIpv4(tip ...string) *FieldValidator {
	if this.GoOn && !Ipv4Reg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a ipv4 format.", tip)
	}
	return this
}

//check the field if it is a ipv6 format
func (this *FieldValidator) IsIpv6(tip ...string) *FieldValidator {
	if this.GoOn && !Ipv6Reg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a ipv6 format.", tip)
	}
	return this
}

//check the field if it is all alpha.
func (this *FieldValidator) IsAlpha(tip ...string) *FieldValidator {
	if this.GoOn && !AlphaReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a alpha format.", tip)
	}
	return this
}

//check the field if it is a numeric format
func (this *FieldValidator) IsNumeric(tip ...string) *FieldValidator {
	if this.GoOn && !NumericReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a numeric format.", tip)
	}
	return this
}

//check the field if it is a alpha and numeric format
func (this *FieldValidator) IsAlphaNumeric(tip ...string) *FieldValidator {
	if this.GoOn && !AlphaNumericReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a alpha numeric format.", tip)
	}
	return this
}

//check the field if it is a base64 format
func (this *FieldValidator) IsBase64(tip ...string) *FieldValidator {
	if this.GoOn && !Base64Reg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a base64 format.", tip)
	}
	return this
}

//check the field if it is a hexadecimal format
func (this *FieldValidator) IsHexadecimal(tip ...string) *FieldValidator {
	if this.GoOn && !HexadecimalReg.MatchString(this.Value) {
		this.FireError(this.Key+" is not a haxadecimal format.", tip)
	}
	return this
}

//check the field if it is a lowercase string
func (this *FieldValidator) IsLowercase(tip ...string) *FieldValidator {
	if this.GoOn && this.Value != strings.ToLower(this.Value) {
		this.FireError(this.Key+" is not lowercase.", tip)
	}
	return this
}

//check the field if it is a uppercase string
func (this *FieldValidator) IsUppercase(tip ...string) *FieldValidator {
	if this.GoOn && this.Value != strings.ToUpper(this.Value) {
		this.FireError(this.Key+" is not uppercase.", tip)
	}
	return this
}

//check the field if it is a uuid format
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

//get field as int
func (this *FieldValidator) Int(tip ...string) int {
	r, e := strconv.Atoi(this.Value)
	if nil != e {
		this.FireError(this.Key+" is not integer.", tip)
	}
	return r
}

//get field as int64
func (this *FieldValidator) Int64(tip ...string) int64 {
	r, e := strconv.ParseInt(this.Value, 10, 64)
	if nil != e {
		this.FireError(this.Key+" is not integer.", tip)
	}
	return r
}

//get field as float64
func (this *FieldValidator) Float(tip ...string) float64 {
	r, e := strconv.ParseFloat(this.Value, 64)
	if nil != e {
		this.FireError(this.Key+" is not float.", tip)
	}

	return r
}

//get field as bool
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

//get field as Time object
func (this *FieldValidator) Date(format string, tip ...string) time.Time {
	r, e := time.Parse(format, this.Value)
	if nil != e {
		panic(e)
	}
	return r
}

//map the field as a md5 string
func (this *FieldValidator) Md5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(this.Value)))
}

//map the field as a sha1 string
func (this *FieldValidator) Sha1() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(this.Value)))
}
