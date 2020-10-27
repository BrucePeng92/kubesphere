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
package tenant

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"kubesphere.io/kubesphere/pkg/apis/tenant/v1alpha1"
	"kubesphere.io/kubesphere/pkg/constants"
	"kubesphere.io/kubesphere/pkg/informers"
	"kubesphere.io/kubesphere/pkg/models"
	iam "kubesphere.io/kubesphere/pkg/models/iam"
	ws "kubesphere.io/kubesphere/pkg/models/workspaces"
	"kubesphere.io/kubesphere/pkg/server/params"
	"kubesphere.io/kubesphere/pkg/simple/client"
	baomi "kubesphere.io/kubesphere/pkg/utils/baomi"
)

var (
	workspaces = workspaceSearcher{}
	namespaces = namespaceSearcher{}
)

func CreateNamespace(workspaceName string, namespace *v1.Namespace, username string) (*v1.Namespace, error) {
	if namespace.Labels == nil {
		namespace.Labels = make(map[string]string, 0)
	}
	if username != "" {
		namespace.Annotations[constants.CreatorAnnotationKey] = username
	}

	user, err := iam.GetUserInfo(username)
	if err != nil {
		return nil, err
	}
	userBaomi := user.Baomi

	namespace.Annotations["scheduler.alpha.kubernetes.io/node-selector"] = "baomi=" + userBaomi
	namespace.Labels[constants.WorkspaceLabelKey] = workspaceName

	return client.ClientSets().K8s().Kubernetes().CoreV1().Namespaces().Create(namespace)
}

func DescribeWorkspace(username, workspaceName string) (*v1alpha1.Workspace, error) {
	workspace, err := informers.KsSharedInformerFactory().Tenant().V1alpha1().Workspaces().Lister().Get(workspaceName)
	if err != nil {
		return nil, err
	}
	annotations := workspace.GetAnnotations()
	workspaceBaomi := annotations["baomi"]
	user, err := iam.GetUserInfo(username)
	if err != nil {
		return nil, err
	}
	userBaomi := user.Baomi
	pass, err := baomi.IsContain(userBaomi, workspaceBaomi)
	if pass || err != nil {
		err := k8serr.NewNotFound(v1alpha1.Resource("workspace"), workspaceName)
		return nil, err
	}

	workspace = appendAnnotations(username, workspace)

	return workspace, nil
}

func ListWorkspaces(username string, conditions *params.Conditions, orderBy string, reverse bool, limit, offset int) (*models.PageableResponse, error) {

	workspaces, err := workspaces.search(username, conditions, orderBy, reverse)

	if err != nil {
		return nil, err
	}

	// limit offset
	result := make([]interface{}, 0)
	for i, workspace := range workspaces {
		if len(result) < limit && i >= offset {
			workspace := appendAnnotations(username, workspace)
			result = append(result, workspace)
		}
	}

	return &models.PageableResponse{Items: result, TotalCount: len(workspaces)}, nil
}

func appendAnnotations(username string, workspace *v1alpha1.Workspace) *v1alpha1.Workspace {
	workspace = workspace.DeepCopy()
	if workspace.Annotations == nil {
		workspace.Annotations = make(map[string]string)
	}
	ns, err := ListNamespaces(username, &params.Conditions{Match: map[string]string{constants.WorkspaceLabelKey: workspace.Name}}, "", false, 1, 0)
	if err == nil {
		workspace.Annotations["kubesphere.io/namespace-count"] = strconv.Itoa(ns.TotalCount)
	}
	devops, err := ListDevopsProjects(workspace.Name, username, &params.Conditions{}, "", false, 1, 0)
	if err == nil {
		workspace.Annotations["kubesphere.io/devops-count"] = strconv.Itoa(devops.TotalCount)
	}
	userCount, err := ws.WorkspaceUserCount(workspace.Name)
	if err == nil {
		workspace.Annotations["kubesphere.io/member-count"] = strconv.Itoa(userCount)
	}
	return workspace
}

func ListNamespaces(username string, conditions *params.Conditions, orderBy string, reverse bool, limit, offset int) (*models.PageableResponse, error) {

	namespaces, err := namespaces.search(username, conditions, orderBy, reverse)

	if err != nil {
		return nil, err
	}

	// limit offset
	result := make([]interface{}, 0)
	for i, v := range namespaces {
		if len(result) < limit && i >= offset {
			result = append(result, v)
		}
	}

	return &models.PageableResponse{Items: result, TotalCount: len(namespaces)}, nil
}
