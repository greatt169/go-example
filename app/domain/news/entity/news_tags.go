package entity

import uuid "github.com/satori/go.uuid"

// Свзяь many-to-many для новостей и акций
type NewsTags struct {
	NewsId uuid.UUID
	TagId  uuid.UUID
}
