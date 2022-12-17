package auth

import "errors"

func (s *Service) HasPermission(token Token, permission string) bool {
	userPermissions, err := s.getUserPermissions(token.UserId)
	if err != nil {
		return false
	}
	groupPermissions, err := s.getGroupPermissions(token.Scope)
	if err != nil {
		return false
	}
	permissions := append(userPermissions, groupPermissions...)
	return s.isInArray(permissions, permission)
}

func (s *Service) getUserPermissions(userId int) ([]Permission, error) {
	var permissions []Permission
	err := s.Server.DB.Select(&permissions, "SELECT * FROM user_permissions WHERE user_id = ?", string(rune(userId)))
	return permissions, err
}

func (s *Service) CreatePermission(permission Permission) (Permission, error) {
	var perm Permission
	result, err := s.Server.DB.NamedExec("INSERT INTO permissions (name, description) VALUES (:name, :description)", permission)
	if err != nil {
		return perm, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return perm, err
	}
	perm, err = s.GetPermissionById(int(id))
	return perm, err
}

func (s *Service) UpdatePermission(permission Permission) (Permission, error) {
	var perm Permission
	_, err := s.Server.DB.NamedExec("UPDATE permissions SET name = :name, description = :description WHERE id = :id", permission)
	if err != nil {
		return perm, err
	}
	perm, err = s.GetPermissionById(permission.ID)
	return perm, err
}

func (s *Service) DeletePermission(id int) error {
	_, err := s.Server.DB.Exec("DELETE FROM permissions WHERE id = ?", string(rune(id)))
	return err
}

func (s *Service) ChangeUserScope(userId int, scope Group) error {
	if s.userHasScope(userId, scope) {
		return errors.New("user already has the scope " + scope.Name)
	}
	_, err := s.Server.DB.Exec("UPDATE users SET scope = ? WHERE id = ?", string(rune(scope.ID)), string(rune(userId)))
	return err
}

func (s *Service) userHasScope(userId int, scope Group) bool {
	user, err := s.User.GetUserById(string(rune(userId)))
	if err != nil {
		return false
	}
	return user.Scope == scope.ID
}

func (s *Service) AddPermissionToGroup(groupId int, permission Permission) error {
	if s.groupHasPermission(groupId, permission.Name) {
		return errors.New("group already has permission")
	}
	_, err := s.Server.DB.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?)", string(rune(groupId)), string(rune(permission.ID)))
	return err
}

func (s *Service) RemovePermissionFromGroup(groupId int, permission Permission) error {
	if !s.groupHasPermission(groupId, permission.Name) {
		return errors.New("group does not have permission")
	}
	_, err := s.Server.DB.Exec("DELETE FROM group_permissions WHERE group_id = ? AND permission_id = ?", string(rune(groupId)), string(rune(permission.ID)))
	return err
}

func (s *Service) RemovePermissionFromUser(userId int, permission Permission) error {
	if !s.userHasPermission(userId, permission.Name) {
		return errors.New("user does not have permission")
	}
	_, err := s.Server.DB.Exec("DELETE FROM user_permissions WHERE user_id = ? AND permission_id = ?", string(rune(userId)), string(rune(permission.ID)))
	return err
}

func (s *Service) AddPermissionToUser(userId int, permission Permission) error {
	if s.userHasPermission(userId, permission.Name) {
		return errors.New("user already has permission")
	}
	_, err := s.Server.DB.Exec("INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?)", string(rune(userId)), string(rune(permission.ID)))
	return err
}

func (s *Service) groupHasPermission(groupId int, permission string) bool {
	permissions, err := s.getGroupPermissions(groupId)
	if err != nil {
		return false
	}
	return s.isInArray(permissions, permission)
}

func (s *Service) userHasPermission(userId int, permission string) bool {
	userPermissions, err := s.getUserPermissions(userId)
	if err != nil {
		return false
	}
	return s.isInArray(userPermissions, permission)
}

func (s *Service) GetPermissionById(id int) (Permission, error) {
	var permission Permission
	err := s.Server.DB.Get(&permission, "SELECT * FROM permissions WHERE id = ?", string(rune(id)))
	return permission, err
}

func (s *Service) getGroupPermissions(groupId int) ([]Permission, error) {
	var permissions []Permission
	err := s.Server.DB.Select(&permissions, "SELECT * FROM group_permissions WHERE group_id = ?", string(rune(groupId)))
	return permissions, err
}

func (s *Service) getGroupByUser(user Token) (Group, error) {
	var group Group
	err := s.Server.DB.Get(&group, "SELECT * FROM groups WHERE id = ?", string(rune(user.Scope)))
	return group, err
}

func (s *Service) getPermissions() ([]Permission, error) {
	var permissions []Permission
	err := s.Server.DB.Select(&permissions, "SELECT * FROM permissions")
	return permissions, err
}

func (s *Service) isInArray(array []Permission, permission string) bool {
	for _, item := range array {
		if item.Name == permission {
			return true
		}
	}
	return false
}
