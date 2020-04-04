package xss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/microcosm-cc/bluemonday"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type XssHandler struct {
	defaultPolicy *bluemonday.Policy
	CustomPolicy *bluemonday.Policy
	FieldsToSkip []string
}

func NewXssHandler() *XssHandler {
	return &XssHandler{
		defaultPolicy: bluemonday.UGCPolicy(),
	}
}

func (s *XssHandler) SanitizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := r.Header.Get("Content-Type")
		method := r.Method

		golog.Debug("Body before XSS: ", r.Body)

		if !strings.Contains(content, "multipart/form-data") && (method == "POST" || method == "GET") {
			var jsonBod interface{}
			d := json.NewDecoder(r.Body)
			d.UseNumber()
			err := d.Decode(&jsonBod)
			if err == nil {
				xmj := jsonBod.(map[string]interface{})
				var sbuff bytes.Buffer
				buff := s.ConstructJson(xmj, sbuff)

				bodOut := buff.String()
				enc := json.NewEncoder(ioutil.Discard)
				err = enc.Encode(&bodOut)
				r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(bodOut)))
			}
		}

		golog.Debug("Body after XSS: ", r.Body)
		next.ServeHTTP(w, r)
	})
}

func (s *XssHandler) ConstructJson(xmj map[string]interface{}, buff bytes.Buffer) bytes.Buffer {
	buff.WriteByte('{')

	m := xmj
	for k, v := range m {
		buff.WriteByte('"')
		buff.WriteString(k)
		buff.WriteByte('"')
		buff.WriteByte(':')

		var b bytes.Buffer
		apndBuff := s.buildJsonApplyPolicy(v, b)
		buff.WriteString(apndBuff.String())
	}
	if len(m) > 0 {
		buff.Truncate(buff.Len() - 1) // remove last ','
	}
	buff.WriteByte('}')

	return buff
}


func (s *XssHandler) buildJsonApplyPolicy(interf interface{}, buff bytes.Buffer) bytes.Buffer {
	switch v := interf.(type) {
	case map[string]interface{}:
		var sbuff bytes.Buffer
		scnd := s.ConstructJson(v, sbuff)
		buff.WriteString(scnd.String())
		buff.WriteByte(',')
	case []interface{}:
		b := s.unravelSlice(v)
		buff.WriteString(b.String())
		buff.WriteByte(',')
	case json.Number:
		buff.WriteString(s.defaultPolicy.Sanitize(fmt.Sprintf("%v", v)))
		buff.WriteByte(',')
	case string:
		buff.WriteString(fmt.Sprintf("%q", s.defaultPolicy.Sanitize(v)))
		buff.WriteByte(',')
	case float64:
		buff.WriteString(s.defaultPolicy.Sanitize(strconv.FormatFloat(v, 'g', 0, 64)))
		buff.WriteByte(',')
	default:
		if v == nil {
			buff.WriteString(fmt.Sprintf("%s", "null"))
			buff.WriteByte(',')
		} else {
			buff.WriteString(s.defaultPolicy.Sanitize(fmt.Sprintf("%v", v)))
			buff.WriteByte(',')
		}
	}
	return buff
}

func (s *XssHandler) unravelSlice(slce []interface{}) bytes.Buffer {
	var buff bytes.Buffer
	buff.WriteByte('[')
	for _, n := range slce {
		switch nn := n.(type) {
		case map[string]interface{}:
			var sbuff bytes.Buffer
			scnd := s.ConstructJson(nn, sbuff)
			buff.WriteString(scnd.String())
			buff.WriteByte(',')
		case string:
			buff.WriteString(fmt.Sprintf("%q", s.defaultPolicy.Sanitize(nn)))
			buff.WriteByte(',')
		case json.Number:
			buff.WriteString(s.defaultPolicy.Sanitize(fmt.Sprintf("%v", nn)))
			buff.WriteByte(',')
		case float64:
			buff.WriteString(s.defaultPolicy.Sanitize(strconv.FormatFloat(nn, 'g', 0, 64)))
			buff.WriteByte(',')
		default:
			if nn == nil {
				buff.WriteString(fmt.Sprintf("%s", "null"))
				buff.WriteByte(',')
			} else {
				buff.WriteString(s.defaultPolicy.Sanitize(fmt.Sprintf("%v", nn)))
				buff.WriteByte(',')
			}
		}
	}
	if len(slce) > 0 {
		buff.Truncate(buff.Len() - 1) // remove last ','
	}
	buff.WriteByte(']')
	return buff
}