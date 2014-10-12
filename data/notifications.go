package data

import "github.com/mdlayher/deltaiota/data/models"

const (
	// sqlSelectNotificationsByUserID is the SQL statement used to select all Notifications
	// for a user, by the user's ID
	sqlSelectNotificationsByUserID = `
		SELECT * FROM notifications WHERE user_id = ?;
	`

	// sqlInsertNotification is the SQL statement used to insert a new Notification
	sqlInsertNotification = `
		INSERT INTO notifications (
			"user_id"
			, "timestamp"
			, "read"
			, "text"
			, "uri"
		) VALUES (?, ?, ?, ?, ?);
	`

	// sqlUpdateNotification is the SQL statement used to update an existing Notification
	sqlUpdateNotification = `
		UPDATE notifications SET
			"user_id" = ?
			, "timestamp" = ?
			, "read" = ?
			, "text" = ?
			, "uri" = ?
		WHERE id = ?;
	`

	// sqlDeleteNotification is the SQL statement used to delete an existing Notification
	sqlDeleteNotification = `
		DELETE FROM notifications WHERE id = ?;
	`

	// sqlDeleteNotificationsByUserID is the SQL statement used to delete all Notifications
	// for a user, by the user's ID
	sqlDeleteNotificationsByUserID = `
		DELETE FROM notifications WHERE user_id = ?;
	`
)

// SelectNotificationsByUserID returns a slice of Notifications by user ID from the database.
func (db *DB) SelectNotificationsByUserID(userID uint64) ([]*models.Notification, error) {
	return db.selectNotifications(sqlSelectNotificationsByUserID, userID)
}

// InsertNotification starts a transaction, inserts a new Notification, and attempts to commit
// the transaction.
func (db *DB) InsertNotification(n *models.Notification) error {
	return db.WithTx(func(tx *Tx) error {
		return tx.InsertNotification(n)
	})
}

// UpdateNotification starts a transaction, updates the input Notification by its ID, and attempts
// to commit the transaction.
func (db *DB) UpdateNotification(n *models.Notification) error {
	return db.WithTx(func(tx *Tx) error {
		return tx.UpdateNotification(n)
	})
}

// DeleteNotification starts a transaction, deletes the input Notification by its ID, and attempts
// to commit the transaction.
func (db *DB) DeleteNotification(n *models.Notification) error {
	return db.WithTx(func(tx *Tx) error {
		return tx.DeleteNotification(n)
	})
}

// DeleteNotificationsByUserID starts a transaction, deletes all Notifications with the matching user ID,
// and attempts to commit the transaction.
func (db *DB) DeleteNotificationsByUserID(userID uint64) error {
	return db.WithTx(func(tx *Tx) error {
		return tx.DeleteNotificationsByUserID(userID)
	})
}

// selectNotifications returns a slice of Notifications from the database, based upon an input
// SQL query and arguments
func (db *DB) selectNotifications(query string, args ...interface{}) ([]*models.Notification, error) {
	// Slice of notifications to return
	var notifications []*models.Notification

	// Invoke closure with prepared statement and wrapped rows,
	// passing any arguments from the caller
	err := db.withPreparedRows(query, func(rows *Rows) error {
		// Scan rows into a slice of Notifications
		var err error
		notifications, err = rows.ScanNotifications()

		// Return errors from scanning
		return err
	}, args...)

	// Return any matching notifications and error
	return notifications, err
}

// InsertNotification inserts a new Notification in the context of the current transaction.
func (tx *Tx) InsertNotification(n *models.Notification) error {
	// Execute SQL to insert Notification
	result, err := tx.Tx.Exec(sqlInsertNotification, n.SQLWriteFields()...)
	if err != nil {
		return err
	}

	// Retrieve generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Store generated ID
	n.ID = uint64(id)
	return nil
}

// UpdateNotification updates the input Notification by its ID, in the context of the
// current transaction.
func (tx *Tx) UpdateNotification(n *models.Notification) error {
	_, err := tx.Tx.Exec(sqlUpdateNotification, n.SQLWriteFields()...)
	return err
}

// DeleteNotification updates the input Notification by its ID, in the context of the
// current transaction.
func (tx *Tx) DeleteNotification(n *models.Notification) error {
	_, err := tx.Tx.Exec(sqlDeleteNotification, n.ID)
	return err
}

// DeleteNotificationsByUserID deletes all Notifications with the input user ID, in the
// context of the current transaction.
func (tx *Tx) DeleteNotificationsByUserID(userID uint64) error {
	_, err := tx.Tx.Exec(sqlDeleteNotificationsByUserID, userID)
	return err
}

// ScanNotifications returns a slice of Notifications from wrapped rows.
func (r *Rows) ScanNotifications() ([]*models.Notification, error) {
	// Iterate all returned rows
	var notifications []*models.Notification
	for r.Rows.Next() {
		// Scan new notification into struct, using specified fields
		n := new(models.Notification)
		if err := r.Rows.Scan(n.SQLReadFields()...); err != nil {
			return nil, err
		}

		// Discard any nil results
		if n == nil {
			continue
		}

		// Append notification to output slice
		notifications = append(notifications, n)
	}

	return notifications, nil
}
