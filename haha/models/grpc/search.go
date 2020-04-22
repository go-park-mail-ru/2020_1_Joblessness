package grpcModels

import (
	"github.com/golang/protobuf/ptypes"
	baseModels "joblessness/haha/models/base"
	searchRpc "joblessness/searchService/rpc"
)

func TransformPersonRPC(p *baseModels.Person) *searchRpc.Person {
	if p == nil {
		return nil
	}

	registered, _ := ptypes.TimestampProto(p.Registered)
	birthday, _ := ptypes.TimestampProto(p.Birthday)

	res := &searchRpc.Person{
		ID:         p.ID,
		Login:      p.Login,
		Password:   p.Password,
		Tag:        p.Tag,
		Email:      p.Email,
		Phone:      p.Phone,
		Registered: registered,
		Avatar:     p.Avatar,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		Gender:     p.Gender,
		Birthday:   birthday,
	}
	return res
}

func TransformPersonBase(p *searchRpc.Person) *baseModels.Person {
	if p == nil {
		return nil
	}

	registered, _ := ptypes.Timestamp(p.Registered)
	birthday, _ := ptypes.Timestamp(p.Birthday)

	res := &baseModels.Person{
		ID:         p.ID,
		Login:      p.Login,
		Password:   p.Password,
		Tag:        p.Tag,
		Email:      p.Email,
		Phone:      p.Phone,
		Registered: registered,
		Avatar:     p.Avatar,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		Gender:     p.Gender,
		Birthday:   birthday,
	}
	return res
}

func TransformOrganizationRPC(p *baseModels.Organization) *searchRpc.Organization {
	if p == nil {
		return nil
	}

	registered, _ := ptypes.TimestampProto(p.Registered)

	res := &searchRpc.Organization{
		ID:         p.ID,
		Login:      p.Login,
		Password:   p.Password,
		Tag:        p.Tag,
		Email:      p.Email,
		Phone:      p.Phone,
		Registered: registered,
		Avatar:     p.Avatar,
		Name:       p.Name,
		About:      p.About,
		Site:       p.Site,
	}
	return res
}

func TransformOrganizationBase(p *searchRpc.Organization) *baseModels.Organization {
	if p == nil {
		return nil
	}

	registered, _ := ptypes.Timestamp(p.Registered)

	res := &baseModels.Organization{
		ID:         p.ID,
		Login:      p.Login,
		Password:   p.Password,
		Tag:        p.Tag,
		Email:      p.Email,
		Phone:      p.Phone,
		Registered: registered,
		Avatar:     p.Avatar,
		Name:       p.Name,
		About:      p.About,
		Site:       p.Site,
	}
	return res
}

func TransformVacanciesRPC(v *baseModels.Vacancy) *searchRpc.Vacancy {
	if v == nil {
		return nil
	}

	res := &searchRpc.Vacancy{
		ID: v.ID,
		Organization: &searchRpc.VacancyOrganization{
			ID:     v.Organization.ID,
			Tag:    v.Organization.Tag,
			Email:  v.Organization.Email,
			Phone:  v.Organization.Phone,
			Avatar: v.Organization.Avatar,
			Name:   v.Organization.Name,
			Site:   v.Organization.Site,
		},
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       int64(v.SalaryFrom),
		SalaryTo:         int64(v.SalaryTo),
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities,
		Conditions:       v.Conditions,
		Keywords:         v.Keywords,
	}
	return res
}

func TransformVacancyBase(v *searchRpc.Vacancy) *baseModels.Vacancy {
	if v == nil {
		return nil
	}

	res := &baseModels.Vacancy{
		ID: v.ID,
		Organization: baseModels.VacancyOrganization{
			ID:     v.Organization.ID,
			Tag:    v.Organization.Tag,
			Email:  v.Organization.Email,
			Phone:  v.Organization.Phone,
			Avatar: v.Organization.Avatar,
			Name:   v.Organization.Name,
			Site:   v.Organization.Site,
		},
		Name:             v.Name,
		Description:      v.Description,
		SalaryFrom:       int(v.SalaryFrom),
		SalaryTo:         int(v.SalaryTo),
		WithTax:          v.WithTax,
		Responsibilities: v.Responsibilities,
		Conditions:       v.Conditions,
		Keywords:         v.Keywords,
	}
	return res
}

func TransformParamsRPC(p *baseModels.SearchParams) *searchRpc.SearchParams {
	if p == nil {
		return nil
	}

	res := &searchRpc.SearchParams{
		Request: p.Request,
		Since:   p.Since,
		Desc:    p.Desc,
	}
	return res
}
