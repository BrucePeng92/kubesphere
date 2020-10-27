/*

 Copyright 2019 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

package quotas

import (
	"net/http"

	"github.com/emicklei/go-restful"

	"kubesphere.io/kubesphere/pkg/server/errors"

	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"kubesphere.io/kubesphere/pkg/constants"
	"kubesphere.io/kubesphere/pkg/models/quotas"
)

func GetNamespaceQuotas(req *restful.Request, resp *restful.Response) {
	namespace := req.PathParameter("namespace")
	username := req.HeaderParameter(constants.UserNameHeader)
	quota, err := quotas.GetNamespaceQuotas(namespace, username)

	if err != nil {
		if k8serr.IsNotFound(err) {
			resp.WriteHeaderAndEntity(http.StatusNotFound, errors.Wrap(err))
		} else {
			resp.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		}
		return
	}

	resp.WriteAsJson(quota)
}

func GetClusterQuotas(req *restful.Request, resp *restful.Response) {
	quota, err := quotas.GetClusterQuotas()

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	resp.WriteAsJson(quota)
}
