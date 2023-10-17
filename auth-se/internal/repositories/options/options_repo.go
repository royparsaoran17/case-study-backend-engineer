package options

type GetAllUsersOptions struct {
	RoleID string
}

type GetAllUsersOption func(options *GetAllUsersOptions)

func GetAllUsersWithRoleFilter(roleID string) GetAllUsersOption {
	return func(options *GetAllUsersOptions) {
		options.RoleID = roleID
	}
}
