// Copyright 2018 TED@Sogou, Inc. All rights reserved.

// @Author: wangzhongning@sogou-inc.com
// @Date: 2018-08-31 14:55:56

package dnspodapi

import "time"

// RespCommon defines common struct in http endpoint
type RespCommon struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
