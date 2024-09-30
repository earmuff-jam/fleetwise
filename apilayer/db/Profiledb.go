package db

import (
	"database/sql"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mohit2530/communityCare/model"
)

// FetchAllUserProfiles ...
func FetchAllUserProfiles(user string) ([]model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `
        SELECT 
            id,
            username,
            full_name,
            COALESCE(ENCODE(avatar_url::bytea, 'base64'), '') AS base64,
			email_address,
			phone_number,
			about_me,
			onlinestatus,
			appearance,
			grid_view,
			role,
			updated_at
        FROM community.profiles;
    `

	var profiles []model.Profile

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var draftProfile model.Profile
		var updated_at sql.NullTime
		var userName, fullName, avatarUrl, emailAddress, phoneNumber, aboutMe, role sql.NullString

		if err := rows.Scan(&draftProfile.ID, &userName, &fullName, &avatarUrl, &emailAddress, &phoneNumber, &aboutMe, &draftProfile.OnlineStatus, &draftProfile.Appearance, &draftProfile.GridView, &role, &updated_at); err != nil {
			return nil, err
		}
		draftProfile.Username = userName.String
		draftProfile.FullName = fullName.String
		draftProfile.AvatarUrl = avatarUrl.String
		draftProfile.EmailAddress = emailAddress.String
		draftProfile.PhoneNumber = phoneNumber.String
		draftProfile.AboutMe = aboutMe.String
		draftProfile.UpdatedAt = updated_at.Time

		profiles = append(profiles, draftProfile)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

// FetchUserProfile ...
func FetchUserProfile(user string, userID string) (*model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `
		SELECT 
			id,
			username,
			full_name,
			CASE 
				WHEN avatar_url IS NOT NULL THEN ENCODE(avatar_url::bytea,'base64')
				ELSE ''
			END AS base64,
			email_address,
			phone_number,
			about_me,
			onlinestatus,
			appearance,
			grid_view,
			role,
			updated_at
		FROM community.profiles
		WHERE id=$1;
	`

	var draftProfile model.Profile

	var profileID sql.NullString
	var userName sql.NullString
	var fullName sql.NullString
	var avatarUrl sql.NullString
	var emailAddress sql.NullString
	var phoneNumber sql.NullString
	var aboutMe sql.NullString
	var role sql.NullString
	var updated_at sql.NullTime

	rows, err := db.Query(sqlStr, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&profileID, &userName, &fullName, &avatarUrl, &emailAddress, &phoneNumber, &aboutMe, &draftProfile.OnlineStatus, &draftProfile.Appearance, &draftProfile.GridView, &role, &updated_at); err != nil {
			return nil, err
		}
		draftProfile.ID, _ = uuid.Parse(profileID.String)
		draftProfile.Username = userName.String
		draftProfile.FullName = fullName.String
		draftProfile.AvatarUrl = avatarUrl.String
		draftProfile.EmailAddress = emailAddress.String
		draftProfile.PhoneNumber = phoneNumber.String
		draftProfile.AboutMe = aboutMe.String
		draftProfile.UpdatedAt = updated_at.Time
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	draftProfile.Validate()
	return &draftProfile, nil
}

// UpdateUserProfile ...
func UpdateUserProfile(user string, userID string, draftProfile model.Profile) (*model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	sqlStr := `UPDATE community.profiles SET
		username=$2,
		full_name=$3,
		email_address=$4,
		phone_number=$5,
		about_me=$6,
		onlinestatus=$7,
		appearance=$8,
		grid_view=$9
		WHERE id=$1
		RETURNING id, username, full_name, avatar_url, email_address, phone_number, about_me, onlinestatus, appearance, grid_view;`

	var updatedProfile model.Profile
	var avatarUrl sql.NullString // Assuming avatar_url is a string column, not bytea

	row := tx.QueryRow(sqlStr,
		userID,
		draftProfile.Username,
		draftProfile.FullName,
		draftProfile.EmailAddress,
		draftProfile.PhoneNumber,
		draftProfile.AboutMe,
		draftProfile.OnlineStatus,
		draftProfile.Appearance,
		draftProfile.GridView,
	)

	err = row.Scan(
		&updatedProfile.ID,
		&updatedProfile.Username,
		&updatedProfile.FullName,
		&avatarUrl,
		&updatedProfile.EmailAddress,
		&updatedProfile.PhoneNumber,
		&updatedProfile.AboutMe,
		&updatedProfile.OnlineStatus,
		&updatedProfile.Appearance,
		&updatedProfile.GridView,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if avatarUrl.Valid {
		updatedProfile.AvatarUrl = avatarUrl.String
	} else {
		updatedProfile.AvatarUrl = ""
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	updatedProfile.Validate()
	return &updatedProfile, nil
}

// UpdateProfileAvatar ...
func UpdateProfileAvatar(user string, userID string, header *multipart.FileHeader, fileBytes []byte) (*model.Profile, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	var updatedProfile model.Profile

	sqlStr := `
		UPDATE community.profiles 
		SET avatar_url = $2
		WHERE id = $1
		RETURNING id, username, full_name, avatar_url, phone_number, goal, about_me, onlinestatus, role
	`

	var avatarUrl sql.NullString // Assuming avatar_url is a string column, not bytea
	var goal sql.NullString

	row := tx.QueryRow(sqlStr, userID, fileBytes)

	err = row.Scan(
		&updatedProfile.ID,
		&updatedProfile.Username,
		&updatedProfile.FullName,
		&avatarUrl,
		&updatedProfile.PhoneNumber,
		&goal,
		&updatedProfile.AboutMe,
		&updatedProfile.OnlineStatus,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if avatarUrl.Valid {
		updatedProfile.AvatarUrl = avatarUrl.String
	} else {
		updatedProfile.AvatarUrl = ""
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	updatedProfile.Validate()
	return &updatedProfile, nil
}

// FetchFavouriteItems ...
func FetchFavouriteItems(user string, userID string, limit int) ([]model.FavouriteItem, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := getFavouriteItems(db, userID, limit)
	if err != nil {
		log.Printf("unable to retrieve favourite items. error: %+v", err)
		return nil, err
	}

	return rows, nil
}

// getFavouriteItems ...
func getFavouriteItems(db *sql.DB, userID string, limit int) ([]model.FavouriteItem, error) {
	sqlStr := `SELECT
		fi.id,
		fi.category_id,
		c."name" AS category_name,
		s."name" AS category_status,
		fi.maintenance_plan_id,
		mp."name" AS maintenance_plan_name,
		ms."name" AS maintenance_status
			FROM community.favourite_items fi
		LEFT JOIN community.category c ON c.id = fi.category_id 
		LEFT JOIN community.statuses s ON s.id = c.status 
		LEFT JOIN community.maintenance_plan mp ON mp.id = fi.maintenance_plan_id
		LEFT JOIN community.statuses ms ON ms.id = mp.status WHERE fi.created_by = $1 
	FETCH FIRST $2 rows only;`

	rows, err := db.Query(sqlStr, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	favouriteItems := []model.FavouriteItem{}

	for rows.Next() {
		var draftFavItem model.FavouriteItem

		var categoryID sql.NullString
		var categoryName sql.NullString
		var categoryStatus sql.NullString
		var maintenancePlanID sql.NullString
		var maintenancePlanName sql.NullString
		var maintenancePlanStatus sql.NullString
		if err := rows.Scan(&draftFavItem.ID, &categoryID, &categoryName, &categoryStatus, &maintenancePlanID, &maintenancePlanName, &maintenancePlanStatus); err != nil {
			return nil, err
		}
		draftFavItem.CategoryID = categoryID.String
		draftFavItem.CategoryName = categoryName.String
		draftFavItem.CategoryStatus = categoryStatus.String
		draftFavItem.MaintenancePlanID = maintenancePlanID.String
		draftFavItem.MaintenancePlanName = maintenancePlanName.String
		draftFavItem.MaintenancePlanStatus = maintenancePlanStatus.String

		favouriteItems = append(favouriteItems, draftFavItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return favouriteItems, nil
}

// SaveFavouriteItem ...
func SaveFavouriteItem(user string, userID string, draftFavItem model.FavouriteItem) ([]model.FavouriteItem, error) {
	db, err := SetupDB(user)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	parsedCreatorID, err := uuid.Parse(draftFavItem.CreatedBy)
	if err != nil {
		return nil, err
	}

	var draftFavItemID string
	var sharableGroups = make([]uuid.UUID, 0)
	sharableGroups = append(sharableGroups, parsedCreatorID)

	categoryID := sql.NullString{}
	if draftFavItem.CategoryID != "" {
		categoryID = sql.NullString{String: draftFavItem.CategoryID, Valid: true}
	}

	maintenancePlanID := sql.NullString{}
	if draftFavItem.MaintenancePlanID != "" {
		maintenancePlanID = sql.NullString{String: draftFavItem.MaintenancePlanID, Valid: true}
	}

	sqlStr := `INSERT INTO community.favourite_items (category_id, maintenance_plan_id, created_by, updated_by, sharable_groups)
			VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		sqlStr,
		categoryID,
		maintenancePlanID,
		parsedCreatorID,
		parsedCreatorID,
		pq.Array(sharableGroups),
	).Scan(&draftFavItemID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rows, err := getFavouriteItems(db, userID, 1000)
	if err != nil {
		log.Printf("unable to retrieve favourite items. error: %+v", err)
		return nil, err
	}

	return rows, nil
}

// RemoveFavItem ...
func RemoveFavItem(user string, userID string, itemID string) (string, error) {
	db, err := SetupDB(user)
	if err != nil {
		return "", err
	}
	defer db.Close()

	parsedCreatorID, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}

	sqlStr := `DELETE FROM community.favourite_items fi WHERE fi.id = $1 AND $2::UUID = ANY(fi.sharable_groups);`
	_, err = db.Exec(sqlStr, itemID, parsedCreatorID)
	if err != nil {
		log.Printf("unable to delete note ID. error: %+v", itemID)
		return "", err
	}

	return itemID, nil
}

// UpdateProfileImage ...
func UpdateProfileImage(user string, userID string, imageURL string) (bool, error) {

	db, err := SetupDB(user)
	if err != nil {
		return false, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("unable to start trasanction with selected db pool. error: %+v", err)
		return false, err
	}
	sqlStr := `UPDATE community.profiles p
		SET associated_image_url = $1,
			updated_at = $4,
			updated_by = $2
			WHERE $2::UUID = ANY(p.sharable_groups) 
		RETURNING p.id;`

	var updatedProfileID string
	err = tx.QueryRow(sqlStr, imageURL, userID, time.Now()).Scan(&updatedProfileID)
	if err != nil {
		log.Printf("unable to update profile id. error: %+v", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("unable to commit. error: %+v", err)
		return false, err
	}

	return true, nil
}
