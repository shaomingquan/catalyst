package gene

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func fetchParam(ctx *gin.Context, name string) interface{} {
	// 1, from  urlencoded
	valFromUrlencoded := ctx.PostForm(name)
	if valFromUrlencoded != "" {
		return valFromUrlencoded
	}

	// 2, from query
	valFromQuery := ctx.Query(name)
	if valFromQuery != "" {
		return valFromQuery
	}

	// 3, from json body, when you post binary, dont use validator, make sure you are posting simple json
	_jsonBody, ok := ctx.Get("webcore:jsonBodyStore")
	var jsonBody []byte
	if !ok {
		jsonBody, _ = ctx.GetRawData()
		ctx.Set("webcore:jsonBodyStore", jsonBody)
	} else {
		jsonBody = _jsonBody.([]byte)
	}

	valsFromJson := map[string]interface{}{}
	json.Unmarshal(jsonBody, &valsFromJson)

	println(string(jsonBody))
	valFromJson := valsFromJson[name]

	stringFromJson, ok := valFromJson.(string)
	println(stringFromJson)
	if ok {
		return stringFromJson
	}
	float64FromJson, ok := valFromJson.(float64)
	if ok {
		return float64FromJson
	}
	boolFromJson, ok := valFromJson.(bool)
	if ok {
		return boolFromJson
	}
	return ""
}

// run time param parsers

func ParamTostring(ctx *gin.Context, name string) string {
	_val := fetchParam(ctx, name)
	val, ok := _val.(string)
	if ok {
		return val
	}
	return ""
}

func ParamTofloat64(ctx *gin.Context, name string) float64 {
	_val := fetchParam(ctx, name)
	val, ok := _val.(float64)
	zero := float64(0)
	if ok {
		return val
	}
	valstring, ok := _val.(string)
	if ok {
		float64val, err := strconv.ParseFloat(valstring, 0)
		if err != nil {
			return zero
		}
		return float64val
	}
	return zero
}

func ParamToint(ctx *gin.Context, name string) int {
	_val := fetchParam(ctx, name)
	val, ok := _val.(int)
	if ok {
		return val
	}
	return 0
}

func ParamTobool(ctx *gin.Context, name string) bool {
	_val := fetchParam(ctx, name)
	val, ok := _val.(bool)
	if ok {
		return val
	}
	return false
}
