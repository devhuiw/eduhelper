package utils

import (
	"context"
	"encoding/json"
	"strconv"
)

func GetUserIDFromContext(ctx context.Context) *int64 {
	val := ctx.Value("user_id")
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case int64:
		return &v
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return &i
		}
	}
	return nil
}

func PtrToStr(s string) *string {
	return &s
}

func PtrToJSON(v interface{}) *string {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	str := string(b)
	return &str
}
