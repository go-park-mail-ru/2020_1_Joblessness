package seed

import (
	"joblessness/haha/auth"
	"joblessness/haha/models"
)

type Seeder struct {
	userUseCase auth.AuthUseCase
}

func NewSeeder(userUseCase *auth.AuthUseCase) *Seeder {
	return &Seeder{
		userUseCase: *userUseCase,
	}
}

func (s *Seeder) CreatePersons() (err error) {
	has, err := s.userUseCase.HasPersons()
	if has || err != nil {
		return err
	}

	err = s.userUseCase.RegisterPerson(&models.Person{
		Login:      "vvputin",
		Password:   "password",
		Tag:        "vvputin",
		Email:      "vvputin@mail.ru",
		Phone:      "88005553535",
		FirstName:  "Vova",
		LastName:   "Putin",
		Gender:     "male",
		Birthday:	"2000-10-10",
	})
	if err != nil {
		return err
	}

	err = s.userUseCase.RegisterPerson(&models.Person{
		Login:      "johncena",
		Password:   "password",
		Tag:        "johncena",
		Email:      "johncena@mail.ru",
		Phone:      "88005553535",
		FirstName:  "John",
		LastName:   "Cena",
		Gender:     "male",
		Birthday:	"2000-10-10",
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Seeder) CreateOrganizations() (err error) {
	has, err := s.userUseCase.HasOrganizations()
	if has || err != nil {
		return err
	}

	err = s.userUseCase.RegisterOrganization(&models.Organization{
		Login:      "haharu",
		Password:   "password",
		Tag:        "haharu",
		Email:      "haharu@mail.ru",
		Phone:      "88005553535",
		Name:       "Haha.ru",
		Site:       "haha.ru",
	})
	if err != nil {
		return err
	}

	err = s.userUseCase.RegisterOrganization(&models.Organization{
		Login:      "mailru",
		Password:   "password",
		Tag:        "mailru",
		Email:      "mailru@mail.ru",
		Phone:      "88005553535",
		Name:       "Mail.ru",
		Site:       "mail.ru",
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Seeder) Seed() (err error) {
	err = s.CreatePersons()
	if err != nil {
		return err
	}

	err = s.CreateOrganizations()
	if err != nil {
		return err
	}

	return nil
}
