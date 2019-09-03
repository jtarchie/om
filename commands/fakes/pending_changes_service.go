// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/om/api"
)

type PendingChangesService struct {
	ListStagedPendingChangesStub        func() (api.PendingChangesOutput, error)
	listStagedPendingChangesMutex       sync.RWMutex
	listStagedPendingChangesArgsForCall []struct {
	}
	listStagedPendingChangesReturns struct {
		result1 api.PendingChangesOutput
		result2 error
	}
	listStagedPendingChangesReturnsOnCall map[int]struct {
		result1 api.PendingChangesOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PendingChangesService) ListStagedPendingChanges() (api.PendingChangesOutput, error) {
	fake.listStagedPendingChangesMutex.Lock()
	ret, specificReturn := fake.listStagedPendingChangesReturnsOnCall[len(fake.listStagedPendingChangesArgsForCall)]
	fake.listStagedPendingChangesArgsForCall = append(fake.listStagedPendingChangesArgsForCall, struct {
	}{})
	fake.recordInvocation("ListStagedPendingChanges", []interface{}{})
	fake.listStagedPendingChangesMutex.Unlock()
	if fake.ListStagedPendingChangesStub != nil {
		return fake.ListStagedPendingChangesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listStagedPendingChangesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PendingChangesService) ListStagedPendingChangesCallCount() int {
	fake.listStagedPendingChangesMutex.RLock()
	defer fake.listStagedPendingChangesMutex.RUnlock()
	return len(fake.listStagedPendingChangesArgsForCall)
}

func (fake *PendingChangesService) ListStagedPendingChangesCalls(stub func() (api.PendingChangesOutput, error)) {
	fake.listStagedPendingChangesMutex.Lock()
	defer fake.listStagedPendingChangesMutex.Unlock()
	fake.ListStagedPendingChangesStub = stub
}

func (fake *PendingChangesService) ListStagedPendingChangesReturns(result1 api.PendingChangesOutput, result2 error) {
	fake.listStagedPendingChangesMutex.Lock()
	defer fake.listStagedPendingChangesMutex.Unlock()
	fake.ListStagedPendingChangesStub = nil
	fake.listStagedPendingChangesReturns = struct {
		result1 api.PendingChangesOutput
		result2 error
	}{result1, result2}
}

func (fake *PendingChangesService) ListStagedPendingChangesReturnsOnCall(i int, result1 api.PendingChangesOutput, result2 error) {
	fake.listStagedPendingChangesMutex.Lock()
	defer fake.listStagedPendingChangesMutex.Unlock()
	fake.ListStagedPendingChangesStub = nil
	if fake.listStagedPendingChangesReturnsOnCall == nil {
		fake.listStagedPendingChangesReturnsOnCall = make(map[int]struct {
			result1 api.PendingChangesOutput
			result2 error
		})
	}
	fake.listStagedPendingChangesReturnsOnCall[i] = struct {
		result1 api.PendingChangesOutput
		result2 error
	}{result1, result2}
}

func (fake *PendingChangesService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listStagedPendingChangesMutex.RLock()
	defer fake.listStagedPendingChangesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PendingChangesService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
