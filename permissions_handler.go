package golang_starter

import (
	"net/http"

	"gopkg.in/matryer/respond.v1"
	"golang_starter/transformers"
)

func (s *Server) ListPermissions() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.Db.Clone()
		defer db.Close()
		permissionList, err := db.Permissions().List()
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "Server Error")
		}
		respond.With(w, r, http.StatusOK, transformers.TransformPermissions(permissionList))
	})
}
