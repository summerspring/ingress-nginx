/*
Copyright 2015 The Kubernetes Authors.

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

package controller

import (
	"reflect"
	"testing"

	"k8s.io/ingress-nginx/pkg/ingress"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/auth"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/authreq"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/cors"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/ipwhitelist"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/proxy"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/ratelimit"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/redirect"
	"k8s.io/ingress-nginx/pkg/ingress/annotations/rewrite"
)

type fakeError struct{}

func (fe *fakeError) Error() string {
	return "fakeError"
}

func TestMergeLocationAnnotations(t *testing.T) {
	// initial parameters
	loc := ingress.Location{}
	annotations := map[string]interface{}{
		"Path":               "/checkpath",
		"IsDefBackend":       true,
		"Backend":            "foo_backend",
		"BasicDigestAuth":    auth.BasicDigest{},
		DeniedKeyName:        &fakeError{},
		"CorsConfig":         cors.CorsConfig{},
		"ExternalAuth":       authreq.External{},
		"RateLimit":          ratelimit.RateLimit{},
		"Redirect":           redirect.Redirect{},
		"Rewrite":            rewrite.Redirect{},
		"Whitelist":          ipwhitelist.SourceRange{},
		"Proxy":              proxy.Configuration{},
		"UsePortInRedirects": true,
	}

	// create test table
	type fooMergeLocationAnnotationsStruct struct {
		fName string
		er    interface{}
	}
	fooTests := []fooMergeLocationAnnotationsStruct{}
	for name, value := range annotations {
		fva := fooMergeLocationAnnotationsStruct{name, value}
		fooTests = append(fooTests, fva)
	}

	// execute test
	mergeLocationAnnotations(&loc, annotations)

	// check result
	for _, foo := range fooTests {
		fv := reflect.ValueOf(loc).FieldByName(foo.fName).Interface()
		if !reflect.DeepEqual(fv, foo.er) {
			t.Errorf("Returned %v but expected %v for the field %s", fv, foo.er, foo.fName)
		}
	}
	if _, ok := annotations[DeniedKeyName]; ok {
		t.Errorf("%s should be removed after mergeLocationAnnotations", DeniedKeyName)
	}
}

func TestIntInSlice(t *testing.T) {
	fooTests := []struct {
		i    int
		list []int
		er   bool
	}{
		{1, []int{1, 2}, true},
		{3, []int{1, 2}, false},
		{1, nil, false},
		{0, nil, false},
	}

	for _, fooTest := range fooTests {
		r := intInSlice(fooTest.i, fooTest.list)
		if r != fooTest.er {
			t.Errorf("returned %t but expected %t for s=%v & list=%v", r, fooTest.er, fooTest.i, fooTest.list)
		}
	}
}

func TestSysctlFSFileMax(t *testing.T) {
	i := sysctlFSFileMax()
	if i < 1 {
		t.Errorf("returned %v but expected > 0", i)
	}
}

func TestSysctlSomaxconn(t *testing.T) {
	i := sysctlSomaxconn()
	if i < 511 {
		t.Errorf("returned %v but expected >= 511", i)
	}
}
