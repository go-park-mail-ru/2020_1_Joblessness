// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package baseModels

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson80a4d695DecodeJoblessnessHahaModelsBase(in *jlexer.Lexer, out *VacancyOrganization) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "tag":
			out.Tag = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "phone":
			out.Phone = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "site":
			out.Site = string(in.String())
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
func easyjson80a4d695EncodeJoblessnessHahaModelsBase(out *jwriter.Writer, in VacancyOrganization) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	if in.Tag != "" {
		const prefix string = ",\"tag\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Tag))
	}
	if in.Email != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	if in.Phone != "" {
		const prefix string = ",\"phone\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Phone))
	}
	if in.Avatar != "" {
		const prefix string = ",\"avatar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Avatar))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Site != "" {
		const prefix string = ",\"site\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Site))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VacancyOrganization) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson80a4d695EncodeJoblessnessHahaModelsBase(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VacancyOrganization) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson80a4d695DecodeJoblessnessHahaModelsBase(l, v)
}
func easyjson80a4d695DecodeJoblessnessHahaModelsBase1(in *jlexer.Lexer, out *Vacancy) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "organization":
			(out.Organization).UnmarshalEasyJSON(in)
		case "name":
			out.Name = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "salaryFrom":
			out.SalaryFrom = int(in.Int())
		case "salaryTo":
			out.SalaryTo = int(in.Int())
		case "withTax":
			out.WithTax = bool(in.Bool())
		case "responsibilities":
			out.Responsibilities = string(in.String())
		case "conditions":
			out.Conditions = string(in.String())
		case "keywords":
			out.Keywords = string(in.String())
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
func easyjson80a4d695EncodeJoblessnessHahaModelsBase1(out *jwriter.Writer, in Vacancy) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	if true {
		const prefix string = ",\"organization\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Organization).MarshalEasyJSON(out)
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	if in.SalaryFrom != 0 {
		const prefix string = ",\"salaryFrom\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.SalaryFrom))
	}
	if in.SalaryTo != 0 {
		const prefix string = ",\"salaryTo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.SalaryTo))
	}
	if in.WithTax {
		const prefix string = ",\"withTax\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.WithTax))
	}
	if in.Responsibilities != "" {
		const prefix string = ",\"responsibilities\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Responsibilities))
	}
	if in.Conditions != "" {
		const prefix string = ",\"conditions\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Conditions))
	}
	if in.Keywords != "" {
		const prefix string = ",\"keywords\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Keywords))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vacancy) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson80a4d695EncodeJoblessnessHahaModelsBase1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vacancy) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson80a4d695DecodeJoblessnessHahaModelsBase1(l, v)
}
func easyjson80a4d695DecodeJoblessnessHahaModelsBase2(in *jlexer.Lexer, out *Vacancies) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Vacancies, 0, 8)
			} else {
				*out = Vacancies{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *Vacancy
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(Vacancy)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson80a4d695EncodeJoblessnessHahaModelsBase2(out *jwriter.Writer, in Vacancies) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vacancies) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson80a4d695EncodeJoblessnessHahaModelsBase2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vacancies) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson80a4d695DecodeJoblessnessHahaModelsBase2(l, v)
}