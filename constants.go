package starter

// Constants structure for storing constants application wide
type Constants struct {
	Permissions Permissions
}

// UserPermission to store constants related to user permission
type UserPermission struct {
	List   string
	Create string
	Edit   string
	Delete string
}

// Permissions store all available types of permission this application supports
type Permissions struct {
	Sudo          string
	User          UserPermission
	Permissions   Permission
	Notifications NotificationPermissions
	Metrics       MetricPermissions
}

// Permission permission related permissions
type Permission struct {
	List string
}

// NotificationPermissions permissions related to notifications
type NotificationPermissions struct {
	Create string
	View   string
	Delete string
}

// MetricPermissions permissions related to metric
type MetricPermissions struct {
	View string
}

// Const returns a Constants struct with populated fields
func Const() Constants {
	return Constants{
		Permissions: Permissions{
			Sudo:          "sudo",
			User:          getUserPermissions(),
			Permissions:   getPermission(),
			Metrics:       getMetricPermissions(),
			Notifications: getNotificationPermissions(),
		},
	}
}

func getMetricPermissions() MetricPermissions {
	return MetricPermissions{
		View: "Metric.View",
	}
}

func getNotificationPermissions() NotificationPermissions {
	return NotificationPermissions{
		Create: "Notification.Create",
		Delete: "Notification.Delete",
		View:   "Notification.View",
	}
}

func getPermission() Permission {
	return Permission{
		List: "Permission.List",
	}
}

func getUserPermissions() UserPermission {
	return UserPermission{
		List:   "User.List",
		Create: "User.Create",
		Edit:   "User.Update",
		Delete: "User.Delete",
	}
}
