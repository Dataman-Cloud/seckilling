/*
Copyright 2014 Rohith All rights reserved.

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

package marathon

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllTasks(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	tasks, err := endpoint.Client.AllTasks(nil)
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, len(tasks.Tasks), 2)
}

func TestAllStagingTasks(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	tasks, err := endpoint.Client.AllTasks(url.Values{"status": []string{"staging"}})
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, len(tasks.Tasks), 0)
}

func TestTaskEndpoints(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	endpoints, err := endpoint.Client.TaskEndpoints(fakeAppNameBroken, 8080, true)
	assert.NoError(t, err)
	assert.NotNil(t, endpoints)
	assert.Equal(t, len(endpoints), 1, t)
	assert.Equal(t, endpoints[0], "10.141.141.10:31045", t)

	endpoints, err = endpoint.Client.TaskEndpoints(fakeAppNameBroken, 8080, false)
	assert.NoError(t, err)
	assert.NotNil(t, endpoints)
	assert.Equal(t, len(endpoints), 2, t)
	assert.Equal(t, endpoints[0], "10.141.141.10:31045", t)
	assert.Equal(t, endpoints[1], "10.141.141.10:31234", t)

	endpoints, err = endpoint.Client.TaskEndpoints(fakeAppNameBroken, 80, true)
	assert.Error(t, err)
}

func TestKillApplicationTasks(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	tasks, err := endpoint.Client.KillApplicationTasks(fakeAppName, "", false)
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
}
