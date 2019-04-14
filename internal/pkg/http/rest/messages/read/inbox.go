// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package read

import (
	"net/http"

	"github.com/mailchain/mailchain/internal/pkg/http/rest/errs"
	"github.com/mailchain/mailchain/internal/pkg/http/rest/params/path"
	"github.com/mailchain/mailchain/internal/pkg/mail"
	"github.com/pkg/errors"
)

func doRead(inboxFunc func(messageID mail.ID) error, w http.ResponseWriter, r *http.Request) {
	messageID, err := path.MessageID(r)
	if err != nil {
		errs.JSONHandler(w, http.StatusNotAcceptable, errors.WithMessage(err, "invalid `message_id`"))
		return
	}

	if err := inboxFunc(messageID); err != nil {
		errs.JSONHandler(w, http.StatusUnprocessableEntity, errors.WithStack(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
