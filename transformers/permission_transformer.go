package transformers

import "github.com/fossapps/starter/db"

// ResponsePermission response friendly version of db.Permission
type ResponsePermission struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

// TransformPermission takes db data and returns response friendly data
func TransformPermission(permission db.Permission) ResponsePermission {
	return ResponsePermission{
		Key:         permission.Key,
		Description: permission.Description,
	}
}

// TransformPermissions takes slice of db data and returns response friendly data
func TransformPermissions(items []db.Permission) []ResponsePermission {
	var listItems []ResponsePermission
	for _, data := range items {
		listItems = append(listItems, TransformPermission(data))
	}
	return listItems
}
