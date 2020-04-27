// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package baseModels

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	chat "joblessness/haha/utils/chat"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)
func easyjsonF9aeba63EncodeJoblessnessHahaModelsBase(out *jwriter.Writer, in Messages) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.From) != 0 {
		const prefix string = ",\"from\":"
		first = false
		out.RawString(prefix[1:])
		{
			out.RawByte('[')
			for v3, v4 := range in.From {
				if v3 > 0 {
					out.RawByte(',')
				}
				if v4 == nil {
					out.RawString("null")
				} else {
					easyjsonF9aeba63EncodeJoblessnessHahaUtilsChat(out, *v4)
				}
			}
			out.RawByte(']')
		}
	}
	if len(in.To) != 0 {
		const prefix string = ",\"to\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.To {
				if v5 > 0 {
					out.RawByte(',')
				}
				if v6 == nil {
					out.RawString("null")
				} else {
					easyjsonF9aeba63EncodeJoblessnessHahaUtilsChat(out, *v6)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Messages) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF9aeba63EncodeJoblessnessHahaModelsBase(w, v)
}
func easyjsonF9aeba63DecodeJoblessnessHahaUtilsChat(in *jlexer.Lexer, out *chat.Message) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "message":
			out.Message = string(in.String())
		case "userOneId":
			out.UserOneId = uint64(in.Uint64())
		case "userOne":
			out.UserOne = string(in.String())
		case "userTwoId":
			out.UserTwoId = uint64(in.Uint64())
		case "userTwo":
			out.UserTwo = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "vacancyId":
			out.VacancyID = uint64(in.Uint64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF9aeba63EncodeJoblessnessHahaUtilsChat(out *jwriter.Writer, in chat.Message) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Message != "" {
		const prefix string = ",\"message\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	if in.UserOneId != 0 {
		const prefix string = ",\"userOneId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.UserOneId))
	}
	if in.UserOne != "" {
		const prefix string = ",\"userOne\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.UserOne))
	}
	if in.UserTwoId != 0 {
		const prefix string = ",\"userTwoId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.UserTwoId))
	}
	if in.UserTwo != "" {
		const prefix string = ",\"userTwo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.UserTwo))
	}
	if true {
		const prefix string = ",\"created\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Created).MarshalJSON())
	}
	if in.VacancyID != 0 {
		const prefix string = ",\"vacancyId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.VacancyID))
	}
	out.RawByte('}')
}
func easyjsonF9aeba63EncodeJoblessnessHahaModelsBase1(out *jwriter.Writer, in Conversations) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			if v9 == nil {
				out.RawString("null")
			} else {
				(*v9).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Conversations) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF9aeba63EncodeJoblessnessHahaModelsBase1(w, v)
}
func easyjsonF9aeba63EncodeJoblessnessHahaModelsBase2(out *jwriter.Writer, in ConversationTitle) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ChatterID != 0 {
		const prefix string = ",\"chatter_id\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ChatterID))
	}
	if in.ChatterName != "" {
		const prefix string = ",\"chatter_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ChatterName))
	}
	if true {
		const prefix string = ",\"interview_date\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.InterviewDate).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ConversationTitle) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF9aeba63EncodeJoblessnessHahaModelsBase2(w, v)
}