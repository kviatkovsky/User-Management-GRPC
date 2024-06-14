package utils

import (
	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users"
	"github.com/kviatkovsky/User-Management-gRPC/internal/grpc/user"
)

func PrepareUserToUpdate(usr *users.User, usrPayload *user.UpdateUserRequest) {
	if usrPayload.GetEmail() != "" && usr.Email != usrPayload.GetEmail() {
		usr.Email = usrPayload.GetEmail()
	}

	if usrPayload.GetPassword() != "" {
		usr.PassHash = []byte(usrPayload.GetPassword())
	}
}
