package crazy_nl_backend

type Constants struct {
	Permissions Permissions
}

type UserPermission struct {
	List   string
	Create string
	Edit   string
	Delete string
}

type Permissions struct {
	Sudo          string
	User          UserPermission
	Permissions   Permission
	Notifications NotificationPermissions
	Metrics       MetricPermissions
}

type Permission struct {
	List string
}

type NotificationPermissions struct {
	Create string
	View   string
	Delete string
}

type MetricPermissions struct {
	View string
}

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
		Edit:   "User.Edit",
		Delete: "User.Delete",
	}
}
