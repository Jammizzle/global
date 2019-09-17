package handler

import (
	"net/http"

	j "github.com/DocHQ/global/helpers/muxhelpers"

	"github.com/sirupsen/logrus"
)

func StatusOK(w http.ResponseWriter, r *http.Request) {
	if err := j.jsonResponse(w, http.StatusOK, j.jsonRes{"status": true}); err != nil {
		logrus.Error(err)
	}
}
