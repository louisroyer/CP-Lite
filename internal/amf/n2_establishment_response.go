// Copyright 2024 Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package amf

import (
	"net/http"

	"github.com/nextmn/json-api/jsonapi"
	"github.com/nextmn/json-api/jsonapi/n1n2"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (amf *Amf) N2EstablishmentResponse(c *gin.Context) {
	var ps n1n2.N2PduSessionRespMsg
	if err := c.BindJSON(&ps); err != nil {
		logrus.WithError(err).Error("could not deserialize")
		c.JSON(http.StatusBadRequest, jsonapi.MessageWithError{Message: "could not deserialize", Error: err})
		return
	}
	pduSession, err := amf.smf.CreateSessionDownlink(c, ps.UeInfo.Header.Ue, ps.UeInfo.Header.Dnn, ps.Gnb, ps.DownlinkTeid)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"ue-ip-addr": ps.UeInfo.Addr,
			"ue":         ps.UeInfo.Header.Ue,
			"gnb":        ps.UeInfo.Header.Gnb,
			"dnn":        ps.UeInfo.Header.Dnn,
		}).Error("could not create downlink path")
		c.JSON(http.StatusInternalServerError, jsonapi.MessageWithError{Message: "could not create downlink path", Error: err})
		return
	}
	logrus.WithFields(logrus.Fields{
		"ue":                ps.UeInfo.Header.Ue.String(),
		"gnb":               ps.UeInfo.Header.Gnb.String(),
		"ip-addr":           ps.UeInfo.Addr,
		"gtp-upf":           pduSession.UplinkFteid.Addr,
		"gtp-uplink-teid":   pduSession.UplinkFteid.Teid,
		"gtp-gnb":           pduSession.DownlinkFteid.Addr,
		"gtp-downlink-teid": pduSession.DownlinkFteid.Teid,
		"dnn":               ps.UeInfo.Header.Dnn,
	}).Info("New PDU Session Established")
}
