package interviewPostgres

import (
	"database/sql"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	"time"
)

type InterviewRepository struct {
	db *sql.DB
}

func NewInterviewRepository(db *sql.DB) *InterviewRepository {
	return &InterviewRepository{db}
}

func (r *InterviewRepository) IsOrganizationVacancy(vacancyID, userID uint64) (err error) {
	var isAuthor bool
	checkAuthor := `SELECT v.organization_id = $2
				FROM vacancy v 
				WHERE v.id = $1`
	if err = r.db.QueryRow(checkAuthor, vacancyID, userID).Scan(&isAuthor); err != nil {
		return err
	}

	if !isAuthor {
		return interviewInterfaces.ErrOrganizationIsNotOwner
	}

	return err
}

func (r *InterviewRepository) ResponseSummary(sendSummary *baseModels.SendSummary) (err error) {
	if sendSummary.InterviewDate.IsZero() {
		sendSummary.InterviewDate = time.Now()
	}
	response := `UPDATE response 
				SET date = CURRENT_TIMESTAMP,
				    approved = $1,
				    rejected = $2,
				    interview_date = $3
				WHERE summary_id = $4
				AND vacancy_id = $5;`
	rows, err := r.db.Exec(response, sendSummary.Accepted, sendSummary.Denied, sendSummary.InterviewDate,
		sendSummary.SummaryID, sendSummary.VacancyID)
	if err != nil {
		return err
	}
	if rowsAf, _ := rows.RowsAffected(); rowsAf == 0 {
		return interviewInterfaces.ErrNoSummaryToRefresh
	}

	return nil
}

func (r *InterviewRepository) SaveMessage(message *chat.Message) (err error) {
	saveMessage := `INSERT INTO message (user_one_id, user_two_id, user_one, user_two, body)
				VALUES ($1, $2, $3, $4, $5);`
	_, err = r.db.Exec(saveMessage, message.UserOneId, message.UserTwoId, message.UserOne, message.UserTwo, message.Message)

	return err
}

func (r *InterviewRepository) getUserSendMessages(parameters *baseModels.ChatParameters) (result []*chat.Message, err error) {
	result = make([]*chat.Message, 0)
	//saveMessage := `SELECT user_one_id, user_two_id, user_one, user_two, body, created
    //				FROM message
	//				WHERE user_one_id = $1
	//				AND user_two_id = $2
	//				ORDER BY created desc
	//				LIMIT $3 OFFSET $4;`
	//rows, err := r.db.Query(saveMessage, parameters.From, parameters.To, 20, parameters.Page*20)
	saveMessage := `SELECT user_one_id, user_two_id, user_one, user_two, body, created
    				FROM message
					WHERE user_one_id = $1
					AND user_two_id = $2
					ORDER BY created desc;`
	rows, err := r.db.Query(saveMessage, parameters.From, parameters.To)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var message chat.Message

		err = rows.Scan(&message.UserOneId, &message.UserTwoId, &message.UserOne, &message.UserTwo, &message.Message,
			&message.Created)
		if err != nil {
			return nil, err
		}

		result = append(result, &message)
	}

	return result, nil
}

func (r *InterviewRepository) GetHistory(parameters *baseModels.ChatParameters) (result baseModels.Messages, err error) {
	from, err := r.getUserSendMessages(parameters)
	if err != nil {
		return baseModels.Messages{}, err
	}

	parameters.From, parameters.To = parameters.To, parameters.From
	to, err := r.getUserSendMessages(parameters)

	return baseModels.Messages{
		From: from,
		To:   to,
	}, err
}

func (r *InterviewRepository) GetResponseCredentials(summaryID, vacancyID uint64) (result *baseModels.SummaryCredentials, err error) {
	result = &baseModels.SummaryCredentials{}

	getPerson := `SELECT u.id, p.name
					FROM summary s
					JOIN users u on s.author = u.id
					JOIN person p on u.person_id = p.id
					WHERE s.id = $1`
	err = r.db.QueryRow(getPerson, summaryID).Scan(&result.UserID, &result.UserName)
	if err != nil {
		return result, err
	}

	getOrg := `SELECT u.id, o.name
					FROM vacancy v
					JOIN users u on v.organization_id = u.id
					JOIN organization o on u.organization_id = o.id
					WHERE v.id = $1`
	err = r.db.QueryRow(getOrg, vacancyID).Scan(&result.OrganizationID, &result.OrganizationName)

	return result, err
}

func (r *InterviewRepository) GetConversations(userID uint64) (result baseModels.Conversations, err error) {
	result = make(baseModels.Conversations, 0)
	getConversations := `SELECT u_to.id, u_to.avatar, 
        CASE WHEN u_to.person_id IS NOT NULL THEN per.name
            WHEN u_to.organization_id IS NOT NULL THEN org.name
        END 
        , u_to.tag, r.interview_date
					FROM response r
					JOIN summary s on r.summary_id = s.id
					JOIN vacancy v on r.vacancy_id = v.id
					JOIN users u_to on (s.author = u_to.id
					or v.organization_id = u_to.id)
					JOIN users u_from on (s.author = u_from.id
					or v.organization_id = u_from.id)
					LEFT JOIN organization org on org.id = u_to.organization_id
					LEFT JOIN person per on per.id = u_to.person_id
					WHERE u_from.id = $1
					AND u_to.id != $1
					AND ((rejected = false
					AND approved = false)
					    OR date >= (CURRENT_TIMESTAMP  - INTERVAL '1 DAY'));`
	rows, err := r.db.Query(getConversations, userID)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var title baseModels.ConversationTitle

		err = rows.Scan(&title.ChatterID, &title.Avatar, &title.ChatterName, &title.Tag, &title.InterviewDate)
		if err != nil {
			return nil, err
		}

		result = append(result, &title)
	}

	return result, nil
}
