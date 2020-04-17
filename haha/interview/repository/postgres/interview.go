package interviewPostgres

import (
	"database/sql"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	"joblessness/haha/models"
	"joblessness/haha/utils/chat"
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

func (r *InterviewRepository) ResponseSummary(sendSummary *models.SendSummary) (err error) {
	response := `UPDATE response 
				SET date = CURRENT_TIMESTAMP,
				    approved = $1,
				    rejected = $2
				WHERE summary_id = $3
				AND vacancy_id = $4;`
	rows, err := r.db.Exec(response, sendSummary.Accepted, sendSummary.Denied, sendSummary.SummaryID, sendSummary.VacancyID)
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

func (r *InterviewRepository) getUserSendMessages(parameters *models.ChatParameters) (result []*chat.Message, err error) {
	result = make([]*chat.Message, 0)
	saveMessage := `SELECT user_one_id, user_two_id, user_one, user_two, body, created
    				FROM message
					WHERE user_one_id = $1
					AND user_two_id = $2
					ORDER BY created desc
					LIMIT $3 OFFSET $4;`
	rows, err := r.db.Query(saveMessage, parameters.From, parameters.To, 20, parameters.Page*20)
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

func (r *InterviewRepository) GetHistory(parameters *models.ChatParameters) (result models.Messages, err error) {
	from, err := r.getUserSendMessages(parameters)
	if err != nil {
		return models.Messages{}, err
	}

	parameters.From, parameters.To = parameters.To, parameters.From
	to, err := r.getUserSendMessages(parameters)

	return models.Messages{
		From: from,
		To:   to,
	}, err
}

func (r *InterviewRepository) GetResponseCredentials(summaryID, vacancyID uint64) (result *models.SummaryCredentials, err error) {
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
	err = r.db.QueryRow(getOrg, vacancyID).Scan(&result.OrganizationID, &result.OrganizationID)

	return result, err
}

func (r *InterviewRepository) GetConversations(userID uint64) (result models.Conversations, err error) {
	result = make(models.Conversations, 0)
	getConversations := `SELECT u.id, u.tag, r.interview_date
					FROM response r
					JOIN summary s on r.summary_id = s.id
					JOIN vacancy v on r.vacancy_id = v.id
					JOIN users u on (s.author = u.id
					and v.organization_id = u.id)
					WHERE u.id = $1
					AND rejected = false
					AND approved = false;`
	rows, err := r.db.Query(getConversations, userID)
	if err != nil {
		return result, err
	}
	rows.Close()

	for rows.Next() {
		var title models.ConversationTitle

		err = rows.Scan(&title.ChatterId, &title.ChatterName, &title.InterviewDate)
		if err != nil {
			return nil, err
		}

		result = append(result, &title)
	}

	return result, nil
}