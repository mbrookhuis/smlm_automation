// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"go.uber.org/zap"

	sumamodels "smlm_automation/pkg/models/susemanager"
	"smlm_automation/pkg/util/rest"
)

// all public interface here for the model

// ISuseManager - description
type ISuseManager interface {
	GetSystemGroupName(negName string) string
	SetK3sDetails(auth AuthParams, systemgroupName string, k3sconfiginput map[string]interface{}) error
	ChangeChannels(auth AuthParams, systemID int, targetedVersion string) error
	GetHost(negName string, sessionKey string) (*AuthParams, error)
	InstallPackages(auth AuthParams, systemID int, pkgs []string, timeout int) error
	GetAuth(sessionkey string) (*AuthParams, error)
}

// ISuseManagerAPI - description
type ISuseManagerAPI interface {
	SuseManagerCall(body []byte, method string, hostname string, path string, sessionKey string) (output *rest.HTTPHelperStruct, er error)
}

type (
	// SumanConfig - SUSE Manager credentials
	SumanConfig struct {
		Host     string
		Login    string
		Password string
		Insecure bool
	}
	// Proxy - general needed information
	Proxy struct {
		cfg               *SumanConfig
		contentTypeHeader map[string]string
		suse              ISuseManagerAPI
		logger            *zap.Logger
		retrycount        int
	}
)

// NewProxy - ope new connection
//
// param: s
// param: suse
// param: logger
// param: retrycount
// return:
func NewProxy(s *SumanConfig, suse ISuseManagerAPI, logger *zap.Logger, retrycount int) IProxy {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	return &Proxy{
		cfg:               s,
		contentTypeHeader: header,
		suse:              suse,
		logger:            logger,
		retrycount:        retrycount,
	}
}

// IProxy - all public interface here for the model
type IProxy interface {
	// activationkey
	ActivationKeyAddChildChannels(auth AuthParams, keyName string, childChannels []string) (int, error)
	ActivationKeyAddServerGroups(auth AuthParams, keyName string, groups []int) (int, error)
	ActivationKeyCreate(auth AuthParams, keyName string, baseChannel string, entitlement []string) (string, error)
	ActivationKeyDelete(auth AuthParams, keyName string) (int, error)
	ActivationKeyGetDetails(auth AuthParams, keyName string) (sumamodels.ActivationkeyGetDetails, error)
	ActivationKeyListActivationKeys(auth AuthParams) ([]sumamodels.ActivationkeyGetDetails, error)
	ActivationKeyRemovePackages(auth AuthParams, keyName string, pckgs []sumamodels.ActivationkeyPackages) (int, error)

	// authentication
	GetSessionKey(body []byte, host string) (string, error)
	SumanLogin() (string, error)
	SumanLogout(auth AuthParams) error
	CheckResponseProgress(auth AuthParams, response *rest.HTTPHelperStruct, timeOut int, systemID int, funcName string) error

	// configchannel
	ConfigChannelListGlobals(auth AuthParams) ([]sumamodels.ConfigChannelListGlobals, error)

	// contentmanagement
	ContentManagementListProjects(auth AuthParams) ([]sumamodels.ContentManagementListProjects, error)
	ContentManagementCreate(auth AuthParams, projectLabel string, name string, description string) (sumamodels.ContentManagementListProjects, error)
	ContentManagementAttachSource(auth AuthParams, projectLabel string, sourceType string, sourceLabel string) (sumamodels.ContentManagementSource, error)
	ContentManagementListFilters(auth AuthParams) ([]sumamodels.ContentManagementFilter, error)
	ContentManagementCreateFilter(auth AuthParams, name string, rule string, entityType string, criteria sumamodels.FilterCriteria) (sumamodels.ContentManagementFilter, error)
	ContentManagementAttachFilter(auth AuthParams, projectLabel string, filterID int) (sumamodels.ContentManagementFilter, error)
	ContentManagementCreateEnvironment(auth AuthParams, projectLabel string, predecessorLabel string, envlabel string, name string, description string) (sumamodels.ContentManagementEnvironment, error)
	ContentManagementBuildProject(auth AuthParams, projectLabel string) (int, error)

	// system
	SystemGetID(auth AuthParams, systemName string) ([]sumamodels.System, error)
	CheckProgress(auth AuthParams, actionID int, timeout int, action string, systemID int) (int, error)
	ListCompleteSystem(auth AuthParams, actionID int) ([]interface{}, error)
	ListInprogressSystem(auth AuthParams, actionID int) ([]interface{}, error)
	ListLatestInstallablePackages(auth AuthParams, systemID int) ([]sumamodels.InstallablePackage, error)
	SchedulePackageRefresh(auth AuthParams, systemID int) error
	ScheduleScriptRun(auth AuthParams, systemID int, timeout int, script string) error
	SystemGetScriptResult(auth AuthParams, actionID int, resultCompleted int) (string, error)
	SystemListActiveSystems(auth AuthParams) ([]sumamodels.ActiveSystem, error)
	SystemListInstalledPackages(auth AuthParams, systemID int) ([]sumamodels.InstalledPackage, error)
	SystemScheduleApplyHighstate(auth AuthParams, systemID int, timeout int) error
	SystemScheduleApplyStates(auth AuthParams, systemID int, stateNames []string, timeout int) error
	SystemScheduleChangeChannels(auth AuthParams, systemID int, basechannel string, childChannel []sumamodels.ChannelSoftwareListChildren) error
	SystemGetSubscribedBaseChannel(auth AuthParams, systemID int) (sumamodels.SubscribedBaseChannel, error)
	SystemScheduleReboot(auth AuthParams, systemID int, timeout int) error

	// sync
	GetSlaves(sessionKey string) ([]sumamodels.Slaves, error)
	SyncSlaveGetSlaveByName(auth AuthParams, slaveFQDN string) (sumamodels.Slaves, error)
	SyncSlaveDelete(auth AuthParams, slaveID int) (int, error)
	SyncSlaveCreate(auth AuthParams, slaveFQDN string, isEnabled bool, allowAllOrgs bool) (sumamodels.Slaves, error)
	SyncMasterGetMasterByLabel(auth AuthParams, slaveFQDN string) (sumamodels.SlavesIssMaster, error)
	SyncMasterDelete(auth AuthParams, masterID int) (int, error)
	SyncMasterCreate(auth AuthParams, masterFQDN string) (sumamodels.SlavesIssMaster, error)
	SyncMasterMakeDefault(auth AuthParams, masterID int) (int, error)
	SyncMasterSetCaCert(auth AuthParams, masterID int, caCert string) (int, error)

	// Channels
	ChannelListSoftwareChannels(auth AuthParams) ([]sumamodels.ChannelListSoftwareChannels, error)
	ChannelSoftwareListChildren(auth AuthParams, label string) ([]sumamodels.ChannelSoftwareListChildren, error)
	ChannelSoftwareCreateRepo(auth AuthParams, label string, typeRepo string, url string) (sumamodels.ChannelSoftwareCreateRepo, error)
	ChannelSoftwareCreate(auth AuthParams, label string, name string, summary string, archLabel string, parentLabel string) (int, error)
	ChannelSoftwareAssociateRepo(auth AuthParams, channelLabel string, repoLabel string) (sumamodels.ChannelSoftwareListChildren, error)
	ChannelSoftwareSyncRepo(auth AuthParams, channelLabel string) (int, error)
	ChannelSoftwareIsExisting(auth AuthParams, label string) (bool, error)
	//	SystemGetSubscribedBaseChannel(auth AuthParams, systemID int) (*sumamodels.SubscribedChannel, error)

	// Add func for formula
	GetFormulasByServerID(auth AuthParams, systemID int) ([]string, error)
	GetFormulasByGroupID(auth AuthParams, groupID int) ([]string, error)
	GetGroupFormulaData(auth AuthParams, groupID int, formulaName string) (interface{}, error)
	GetSystemFormulaData(auth AuthParams, sid int, formulaname string) (interface{}, error)
	SetGroupFormulaData(auth AuthParams, groupID int, formulaName string, formulaData interface{}) (int, error)
	SetSystemFormulaData(auth AuthParams, systemID int, formulaName string, formulaData interface{}) (int, error)
	FormulaSetFormulasOfGroup(auth AuthParams, systemID int, formulaNames []string) (int, error)
	FormulaSetFormulasOfSystem(auth AuthParams, systemID int, formulaNames []string) (int, error)

	// SystemGroup
	SystemGroupCreate(auth AuthParams, groupName string, description string) (*sumamodels.SystemGroupGetDetails, error)
	SystemGroupGetDetails(auth AuthParams, groupName string) (*sumamodels.SystemGroupGetDetails, error)
	SystemGroupListSystemsMinimal(auth AuthParams, groupName string) ([]sumamodels.SystemGroupListSystemsMinimal, error)
	SystemGroupListActiveSystemsInGroup(auth AuthParams, groupName string) ([]int, error)

	// KickstartTree
	KickstartTreeGetDetails(auth AuthParams, distributionName string) (sumamodels.KickstartTreeGetDetails, error)
	KickstartTreeCreate(auth AuthParams, treeLabel string, basePath string, channelLabel string, installType string) (int, error)
	KickstartTreeCreateKernelOptions(auth AuthParams, treeLabel string, basePath string, channelLabel string, installType string, kernelOptions string, postKernelOptions string) (int, error)
	KickstartImportRawFile(auth AuthParams, profileLabel string, virtType string, channelLabel string, dataXML string) (int, error)
	KickstartListKickstarts(auth AuthParams) ([]sumamodels.KickstartListProfiles, error)
	KickstartDeleteProfile(auth AuthParams, profileName string) (int, error)
	KickstartProfileSetVariables(auth AuthParams, profileLabel string, profileVariables interface{}) (int, error)
}
