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

func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase(in *jlexer.Lexer, out *VacancyResponse) {
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
		case "user_id":
			out.UserID = uint64(in.Uint64())
		case "tag":
			out.Tag = string(in.String())
		case "vacancyId":
			out.VacancyID = uint64(in.Uint64())
		case "summaryId":
			out.SummaryID = uint64(in.Uint64())
		case "firstName":
			out.FirstName = string(in.String())
		case "lastName":
			out.LastName = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "accepted":
			out.Accepted = bool(in.Bool())
		case "denied":
			out.Denied = bool(in.Bool())
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
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase(out *jwriter.Writer, in VacancyResponse) {
	out.RawByte('{')
	first := true
	_ = first
	if in.UserID != 0 {
		const prefix string = ",\"user_id\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.UserID))
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
	if in.SummaryID != 0 {
		const prefix string = ",\"summaryId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.SummaryID))
	}
	if in.FirstName != "" {
		const prefix string = ",\"firstName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.FirstName))
	}
	if in.LastName != "" {
		const prefix string = ",\"lastName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.LastName))
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
	if in.Accepted {
		const prefix string = ",\"accepted\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Accepted))
	}
	if in.Denied {
		const prefix string = ",\"denied\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Denied))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VacancyResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VacancyResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase1(in *jlexer.Lexer, out *Summary) {
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
		case "author":
			(out.Author).UnmarshalEasyJSON(in)
		case "name":
			out.Name = string(in.String())
		case "salaryFrom":
			out.SalaryFrom = int(in.Int())
		case "salaryTo":
			out.SalaryTo = int(in.Int())
		case "keywords":
			out.Keywords = string(in.String())
		case "educations":
			if in.IsNull() {
				in.Skip()
				out.Educations = nil
			} else {
				in.Delim('[')
				if out.Educations == nil {
					if !in.IsDelim(']') {
						out.Educations = make([]Education, 0, 1)
					} else {
						out.Educations = []Education{}
					}
				} else {
					out.Educations = (out.Educations)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Education
					(v1).UnmarshalEasyJSON(in)
					out.Educations = append(out.Educations, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "experiences":
			if in.IsNull() {
				in.Skip()
				out.Experiences = nil
			} else {
				in.Delim('[')
				if out.Experiences == nil {
					if !in.IsDelim(']') {
						out.Experiences = make([]Experience, 0, 1)
					} else {
						out.Experiences = []Experience{}
					}
				} else {
					out.Experiences = (out.Experiences)[:0]
				}
				for !in.IsDelim(']') {
					var v2 Experience
					(v2).UnmarshalEasyJSON(in)
					out.Experiences = append(out.Experiences, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase1(out *jwriter.Writer, in Summary) {
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
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Author).MarshalEasyJSON(out)
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
	if len(in.Educations) != 0 {
		const prefix string = ",\"educations\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v3, v4 := range in.Educations {
				if v3 > 0 {
					out.RawByte(',')
				}
				(v4).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if len(in.Experiences) != 0 {
		const prefix string = ",\"experiences\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Experiences {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Summary) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Summary) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase1(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase2(in *jlexer.Lexer, out *Summaries) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Summaries, 0, 8)
			} else {
				*out = Summaries{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 *Summary
			if in.IsNull() {
				in.Skip()
				v7 = nil
			} else {
				if v7 == nil {
					v7 = new(Summary)
				}
				(*v7).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase2(out *jwriter.Writer, in Summaries) {
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
func (v Summaries) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Summaries) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase2(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase3(in *jlexer.Lexer, out *SendSummary) {
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
		case "vacancyId":
			out.VacancyID = uint64(in.Uint64())
		case "summaryId":
			out.SummaryID = uint64(in.Uint64())
		case "user_id":
			out.UserID = uint64(in.Uint64())
		case "organizationId":
			out.OrganizationID = uint64(in.Uint64())
		case "interview_date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.InterviewDate).UnmarshalJSON(data))
			}
		case "accepted":
			out.Accepted = bool(in.Bool())
		case "denied":
			out.Denied = bool(in.Bool())
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
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase3(out *jwriter.Writer, in SendSummary) {
	out.RawByte('{')
	first := true
	_ = first
	if in.VacancyID != 0 {
		const prefix string = ",\"vacancyId\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.VacancyID))
	}
	if in.SummaryID != 0 {
		const prefix string = ",\"summaryId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.SummaryID))
	}
	if in.UserID != 0 {
		const prefix string = ",\"user_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.UserID))
	}
	if in.OrganizationID != 0 {
		const prefix string = ",\"organizationId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.OrganizationID))
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
	if in.Accepted {
		const prefix string = ",\"accepted\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Accepted))
	}
	if in.Denied {
		const prefix string = ",\"denied\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Denied))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SendSummary) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SendSummary) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase3(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase4(in *jlexer.Lexer, out *OrgSummaries) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(OrgSummaries, 0, 8)
			} else {
				*out = OrgSummaries{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v10 *VacancyResponse
			if in.IsNull() {
				in.Skip()
				v10 = nil
			} else {
				if v10 == nil {
					v10 = new(VacancyResponse)
				}
				(*v10).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v10)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase4(out *jwriter.Writer, in OrgSummaries) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v11, v12 := range in {
			if v11 > 0 {
				out.RawByte(',')
			}
			if v12 == nil {
				out.RawString("null")
			} else {
				(*v12).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OrgSummaries) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase4(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OrgSummaries) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase4(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase5(in *jlexer.Lexer, out *Experience) {
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
		case "companyName":
			out.CompanyName = string(in.String())
		case "role":
			out.Role = string(in.String())
		case "responsibilities":
			out.Responsibilities = string(in.String())
		case "start":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Start).UnmarshalJSON(data))
			}
		case "stop":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Stop).UnmarshalJSON(data))
			}
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
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase5(out *jwriter.Writer, in Experience) {
	out.RawByte('{')
	first := true
	_ = first
	if in.CompanyName != "" {
		const prefix string = ",\"companyName\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.CompanyName))
	}
	if in.Role != "" {
		const prefix string = ",\"role\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Role))
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
	if true {
		const prefix string = ",\"start\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Start).MarshalJSON())
	}
	if true {
		const prefix string = ",\"stop\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Stop).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Experience) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase5(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Experience) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase5(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase6(in *jlexer.Lexer, out *Education) {
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
		case "institution":
			out.Institution = string(in.String())
		case "speciality":
			out.Speciality = string(in.String())
		case "graduated":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Graduated).UnmarshalJSON(data))
			}
		case "type":
			out.Type = string(in.String())
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
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase6(out *jwriter.Writer, in Education) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Institution != "" {
		const prefix string = ",\"institution\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Institution))
	}
	if in.Speciality != "" {
		const prefix string = ",\"speciality\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Speciality))
	}
	if true {
		const prefix string = ",\"graduated\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Graduated).MarshalJSON())
	}
	if in.Type != "" {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Education) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase6(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Education) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase6(l, v)
}
func easyjsonF381ebcaDecodeJoblessnessHahaModelsBase7(in *jlexer.Lexer, out *Author) {
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
		case "firstName":
			out.FirstName = string(in.String())
		case "lastName":
			out.LastName = string(in.String())
		case "gender":
			out.Gender = string(in.String())
		case "birthday":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Birthday).UnmarshalJSON(data))
			}
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
func easyjsonF381ebcaEncodeJoblessnessHahaModelsBase7(out *jwriter.Writer, in Author) {
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
	if in.FirstName != "" {
		const prefix string = ",\"firstName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.FirstName))
	}
	if in.LastName != "" {
		const prefix string = ",\"lastName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.LastName))
	}
	if in.Gender != "" {
		const prefix string = ",\"gender\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Gender))
	}
	if true {
		const prefix string = ",\"birthday\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Birthday).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Author) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF381ebcaEncodeJoblessnessHahaModelsBase7(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Author) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF381ebcaDecodeJoblessnessHahaModelsBase7(l, v)
}
