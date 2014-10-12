package models

// Notification represents an application notification.
type Notification struct {
	ID        uint64 `db:"id" json:"id"`
	UserID    uint64 `db:"user_id" json:"userId"`
	Timestamp uint64 `db:"timestamp" json:"timestamp"`
	Read      bool   `db:"read" json:"read"`
	Text      string `db:"text" json:"text"`
	URI       string `db:"uri" json:"uri"`
}

// SQLReadFields returns the correct field order to scan SQL row results into the
// receiving Notification struct.
func (n *Notification) SQLReadFields() []interface{} {
	return []interface{}{
		&n.ID,
		&n.UserID,
		&n.Timestamp,
		&n.Read,
		&n.Text,
		&n.URI,
	}
}

// SQLWriteFields returns the correct field order for SQL write actions (such as
// insert or update), for the receiving Notification struct.
func (n *Notification) SQLWriteFields() []interface{} {
	return []interface{}{
		n.UserID,
		n.Timestamp,
		n.Read,
		n.Text,
		n.URI,

		// Last argument for WHERE clause
		n.ID,
	}
}
